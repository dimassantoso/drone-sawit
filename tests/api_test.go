package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const ApiUrl = "http://localhost:8080"

func TestApi(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip API tests")
	}

	testcases := getTestCases()
	ctx := context.Background()
	client := &http.Client{}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			for idx := range tc.Steps {
				step := &tc.Steps[idx]
				request, err := step.Request(t, ctx, &tc)
				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("Accept", "application/json")
				require.NoError(t, err)

				// Send request
				response, err := client.Do(request)

				require.NoError(t, err)
				defer response.Body.Close()

				// Check response
				ReadJsonResult(t, response, step)
				step.Expect(t, ctx, &tc, response, step.Result)
			}
		})
	}
}

func getTestCases() []TestCase {
	return []TestCase{
		//----- Test for API
		{
			Name: "Test Error 1",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest("POST", ApiUrl+"/estate", nil)
					},
					Expect: ExpectBadRequest(),
				},
			},
		},
		{
			Name: "Test Error 2: Invalid Format",
			Steps: []TestCaseStep{
				{
					Request: SendRequestNewEstate(-1, -5),
					Expect:  ExpectBadRequest(),
				},
			},
		},
		{
			Name: "Test Error: Create Tree Out of Bound",
			Steps: []TestCaseStep{
				{
					Request: SendRequestNewEstate(10, 20),
					Expect:  ExpectNewEstateOk(),
				},
				{
					Request: SendRequestNewTree(5, 0, 0),
					Expect:  ExpectBadRequest(),
				},
			},
		},
		CreateNormalTestCase("Normal 1", []any{
			[]any{CreateEstate, 10, 20},
			[]any{CreateTree, 10, 5, 5},
			[]any{CreateTree, 20, 6, 5},
		}),
		CreateNormalTestCase("Normal 2", []any{
			[]any{CreateEstate, 5, 1},
			[]any{CreateTree, 10, 2, 1},
			[]any{CreateTree, 20, 3, 1},
			[]any{CreateTree, 10, 4, 1},
			[]any{GetStats, 3, 10, 20, 10},
			[]any{GetDronePlan, 0, 42},
		}),
	}
}

type TestCase struct {
	Name  string
	Steps []TestCaseStep
}

type RequestFunc func(*testing.T, context.Context, *TestCase) (*http.Request, error)
type ExpectFunc func(*testing.T, context.Context, *TestCase, *http.Response, map[string]any)

type TestCaseStep struct {
	Request RequestFunc
	Expect  ExpectFunc
	Result  map[string]any
}

func ResponseContains(t *testing.T, resp *http.Response, text string) {
	body, err := io.ReadAll(resp.Body)
	bodyStr := string(body)
	require.NoError(t, err)
	require.Contains(t, bodyStr, text)
}

func ReadJsonResult(t *testing.T, resp *http.Response, step *TestCaseStep) {
	var result map[string]any
	err := json.NewDecoder(resp.Body).Decode(&result)
	step.Result = result
	require.NoError(t, err)
}

func RequireIsUUID(t *testing.T, value string) {
	_, err := uuid.Parse(value)
	require.NoError(t, err)
}

const (
	CreateEstate = iota
	CreateTree
	GetStats
	GetDronePlan
)

func CreateNormalTestCase(name string, a []any) TestCase {
	tc := TestCase{}
	tc.Name = name

	for _, step := range a {
		switch step.([]any)[0].(int) {
		case CreateEstate:
			tc.Steps = append(tc.Steps, TestCaseStep{
				Request: SendRequestNewEstate(step.([]any)[1].(int), step.([]any)[2].(int)),
				Expect:  ExpectNewEstateOk(),
			})
		case CreateTree:
			tc.Steps = append(tc.Steps, TestCaseStep{
				Request: SendRequestNewTree(step.([]any)[1].(int), step.([]any)[2].(int), step.([]any)[3].(int)),
				Expect:  ExpectNewTreeOk(),
			})
		case GetStats:
			tc.Steps = append(tc.Steps, TestCaseStep{
				Request: SendRequestGetStats(),
				Expect:  ExpectGetStatsOk(step.([]any)[1].(int), step.([]any)[2].(int), step.([]any)[3].(int), step.([]any)[4].(int)),
			})
		case GetDronePlan:
			tc.Steps = append(tc.Steps, TestCaseStep{
				Request: SendRequestGetDronePlan(step.([]any)[1].(int)),
				Expect:  ExpectGetDronePlanOk(step.([]any)[2].(int)),
			})
		}

	}
	return tc
}

func SendRequestNewEstate(length, width int) RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
		req := map[string]int{
			"length": length,
			"width":  width,
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)
		return http.NewRequest("POST", ApiUrl+"/estate", bytes.NewReader(body))
	}
}

func ExpectNewEstateOk() ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		RequireReturnIsUUID(t, resp, data)
	}
}

func SendRequestNewTree(height, x, y int) RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
		req := map[string]int{
			"height": height,
			"x":      x,
			"y":      y,
		}
		id := tc.Steps[0].Result["id"].(string)
		body, err := json.Marshal(req)
		require.NoError(t, err)
		return http.NewRequest("POST", ApiUrl+"/estate/"+id+"/tree", bytes.NewReader(body))
	}
}

func ExpectNewTreeOk() ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		RequireReturnIsUUID(t, resp, data)
	}
}

func SendRequestGetStats() RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
		id := tc.Steps[0].Result["id"].(string)
		return http.NewRequest("GET", ApiUrl+"/estate/"+id+"/stats", nil)
	}
}

func ExpectGetStatsOk(count, min, max, median int) ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		RequireStats(t, resp, data, count, min, max, median)
	}
}

func SendRequestGetDronePlan(distance int) RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
		id := tc.Steps[0].Result["id"].(string)
		var url string

		if distance == 0 {
			url = fmt.Sprintf("%s/estate/%s/drone-plan", ApiUrl, id)
		} else {
			url = fmt.Sprintf("%s/estate/%s/drone-plan?distance=%d", ApiUrl, id, distance)
		}
		return http.NewRequest("GET", url, nil)
	}
}

func ExpectGetDronePlanOk(distance int) ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		RequireDistance(t, resp, data, distance)
	}
}

func RequireReturnIsUUID(t *testing.T, resp *http.Response, data map[string]any) {
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	RequireIsUUID(t, data["id"].(string))
}

func RequireStats(t *testing.T, resp *http.Response, data map[string]any, count, min, max, median int) {
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, count, int(data["count"].(float64)))
	require.Equal(t, min, int(data["min"].(float64)))
	require.Equal(t, max, int(data["max"].(float64)))
	require.Equal(t, median, int(data["median"].(float64)))
}

func RequireDistance(t *testing.T, resp *http.Response, data map[string]any, distance int) {
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, distance, int(data["distance"].(float64)))
}

func ExpectBadRequest() ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}
