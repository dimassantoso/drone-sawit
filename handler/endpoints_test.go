package handler

import (
	"errors"
	"github.com/dimassantoso/drone-sawit/generated"
	mockrepo "github.com/dimassantoso/drone-sawit/mocks/repository"
	"github.com/dimassantoso/drone-sawit/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostEstate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().CreateEstate(gomock.Any(), gomock.Any()).Return(nil)

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(`{"width": 3, "length": 4}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstate(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), `"id"`)
	})

	t.Run("Failed: invalid input", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(`{"width": -1}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstate(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "width and length must be between 1 and 50000")
	})

	t.Run("Failed: in repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().CreateEstate(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(`{"width": 5, "length": 10}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstate(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `failed`)
	})

	t.Run("Failed: invalid body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(`{"width": 3, "length":`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstate(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestServer_PostEstateIdTree(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{
			BaseModel: repository.BaseModel{
				ID: uuid.NewString(),
			},
			Width:  10,
			Length: 10,
		}, nil)
		mockRepo.EXPECT().FindEstateTree(gomock.Any(), gomock.Any()).Return(repository.EstateTree{}, errors.New("not found"))
		mockRepo.EXPECT().CreateEstateTree(gomock.Any(), gomock.Any()).Return(nil)

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate/:id/tree", strings.NewReader(`{"x": 1, "y": 1, "height": 30}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstateIdTree(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), `"id"`)
	})

	t.Run("Failed: invalid estate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{}, errors.New("not found"))

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate/:id/tree", strings.NewReader(`{"x": 1, "y": 1, "height": 30}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstateIdTree(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), `not found`)
	})

	t.Run("Failed: invalid estate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate/:id/tree", strings.NewReader(`{"x": 1, "y": 1, "height": 50}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstateIdTree(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `x and y must be greatest equal 1, height must be between 1 and 30`)
	})

	t.Run("Failed: Out of Bound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{
			BaseModel: repository.BaseModel{
				ID: uuid.NewString(),
			},
			Width:  10,
			Length: 10,
		}, nil)

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate/:id/tree", strings.NewReader(`{"x": 200, "y": 1, "height": 30}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstateIdTree(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `"coordinate out of bound"`)
	})

	t.Run("Failed: plot have tree", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{
			BaseModel: repository.BaseModel{
				ID: uuid.NewString(),
			},
			Width:  10,
			Length: 10,
		}, nil)
		mockRepo.EXPECT().FindEstateTree(gomock.Any(), gomock.Any()).Return(repository.EstateTree{X: 1, Y: 1}, nil)

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate/:id/tree", strings.NewReader(`{"x": 1, "y": 1, "height": 30}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstateIdTree(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), `plot`)
	})

	t.Run("Failed: invalid body payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate/:id/tree", strings.NewReader(`{"x": 1, "`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstateIdTree(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Failed : error create", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{
			BaseModel: repository.BaseModel{
				ID: uuid.NewString(),
			},
			Width:  10,
			Length: 10,
		}, nil)
		mockRepo.EXPECT().FindEstateTree(gomock.Any(), gomock.Any()).Return(repository.EstateTree{}, errors.New("not found"))
		mockRepo.EXPECT().CreateEstateTree(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/estate/:id/tree", strings.NewReader(`{"x": 1, "y": 1, "height": 30}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.PostEstateIdTree(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestServer_GetEstateIdStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{}, nil)
		mockRepo.EXPECT().CountEstateTree(gomock.Any(), gomock.Any()).Return(10)
		mockRepo.EXPECT().GetEstateTreeStats(gomock.Any(), gomock.Any()).Return(repository.EstateTreeStats{
			Min:    10,
			Max:    30,
			Median: 15,
		}, nil)

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetEstateIdStats(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"count"`)
		assert.Contains(t, rec.Body.String(), `"max"`)
		assert.Contains(t, rec.Body.String(), `"min"`)
		assert.Contains(t, rec.Body.String(), `"median"`)
	})

	t.Run("Failed: estate not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{}, errors.New("not found"))

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetEstateIdStats(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Failed : fetch stats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{}, nil)
		mockRepo.EXPECT().CountEstateTree(gomock.Any(), gomock.Any()).Return(10)
		mockRepo.EXPECT().GetEstateTreeStats(gomock.Any(), gomock.Any()).Return(repository.EstateTreeStats{}, errors.New("unexpected error"))

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetEstateIdStats(c, estateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestServer_GetEstateIdDronePlan(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{Width: 5, Length: 1}, nil)
		mockRepo.EXPECT().FindAllMapEstateTree(gomock.Any(), gomock.Any()).Return(map[repository.CoordinatePoint]repository.EstateTree{
			repository.CoordinatePoint{X: 2, Y: 1}: {X: 2, Y: 1, Height: 5},
			repository.CoordinatePoint{X: 3, Y: 1}: {X: 3, Y: 1, Height: 3},
			repository.CoordinatePoint{X: 4, Y: 1}: {X: 4, Y: 1, Height: 4},
		}, nil)

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/:id/drone-plan", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetEstateIdDronePlan(c, estateID, generated.GetEstateIdDronePlanParams{MaxDistance: nil})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `54`)
	})

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{Width: 6, Length: 3}, nil)
		mockRepo.EXPECT().FindAllMapEstateTree(gomock.Any(), gomock.Any()).Return(map[repository.CoordinatePoint]repository.EstateTree{
			repository.CoordinatePoint{X: 3, Y: 1}: {X: 3, Y: 1, Height: 10},
			repository.CoordinatePoint{X: 3, Y: 2}: {X: 3, Y: 2, Height: 30},
			repository.CoordinatePoint{X: 4, Y: 2}: {X: 4, Y: 2, Height: 14},
			repository.CoordinatePoint{X: 6, Y: 2}: {X: 6, Y: 2, Height: 24},
			repository.CoordinatePoint{X: 5, Y: 3}: {X: 5, Y: 3, Height: 6},
		}, nil)

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/:id/drone-plan", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetEstateIdDronePlan(c, estateID, generated.GetEstateIdDronePlanParams{MaxDistance: nil})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `312`)
	})

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		maxDistance := []int{100, 80, 1000, 15}
		for _, v := range maxDistance {
			estateID := uuid.NewString()

			mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
			mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{Width: 6, Length: 3}, nil)
			mockRepo.EXPECT().FindAllMapEstateTree(gomock.Any(), gomock.Any()).Return(map[repository.CoordinatePoint]repository.EstateTree{
				repository.CoordinatePoint{X: 3, Y: 1}: {X: 3, Y: 1, Height: 10},
				repository.CoordinatePoint{X: 3, Y: 2}: {X: 3, Y: 2, Height: 30},
				repository.CoordinatePoint{X: 4, Y: 2}: {X: 4, Y: 2, Height: 14},
				repository.CoordinatePoint{X: 6, Y: 2}: {X: 6, Y: 2, Height: 24},
				repository.CoordinatePoint{X: 5, Y: 3}: {X: 5, Y: 3, Height: 6},
			}, nil)

			handler := NewServer(NewServerOptions{Repository: mockRepo})

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/estate/:id/drone-plan", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GetEstateIdDronePlan(c, estateID, generated.GetEstateIdDronePlanParams{MaxDistance: &v})
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), `rest`)
		}
	})

	t.Run("Failed : not found estate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{}, errors.New("not found"))

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/:id/drone-plan", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetEstateIdDronePlan(c, estateID, generated.GetEstateIdDronePlanParams{MaxDistance: nil})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Failed : failed fetch tree data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		estateID := uuid.NewString()

		mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
		mockRepo.EXPECT().FindEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{Width: 5, Length: 1}, nil)
		mockRepo.EXPECT().FindAllMapEstateTree(gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid query"))

		handler := NewServer(NewServerOptions{Repository: mockRepo})

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/estate/:id/drone-plan", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetEstateIdDronePlan(c, estateID, generated.GetEstateIdDronePlanParams{MaxDistance: nil})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
