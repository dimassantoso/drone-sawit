package handler

import (
	"fmt"
	"github.com/dimassantoso/drone-sawit/generated"
	"github.com/dimassantoso/drone-sawit/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
)

func (s *Server) PostEstate(c echo.Context) error {
	ctx := c.Request().Context()

	var (
		req         generated.EstateRequest
		errResponse generated.ErrorResponse
	)

	if err := c.Bind(&req); err != nil {
		errResponse.Message = "invalid request body"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	if req.Width <= 0 || req.Length <= 0 || req.Width > 50000 || req.Length > 50000 {
		errResponse.Message = "width and length must be between 1 and 50000"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	estate := repository.Estate{
		BaseModel: repository.BaseModel{
			ID: uuid.NewString(),
		},
		Width:  req.Width,
		Length: req.Length,
	}

	if err := s.Repository.CreateEstate(ctx, &estate); err != nil {
		errResponse.Message = "failed to create estate"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	return c.JSON(http.StatusCreated, generated.EstateResponse{
		Id: estate.ID,
	})
}

func (s *Server) PostEstateIdTree(c echo.Context, estateID string) error {
	ctx := c.Request().Context()

	var (
		req         generated.EstateTreeRequest
		errResponse generated.ErrorResponse
	)

	if err := c.Bind(&req); err != nil {
		errResponse.Message = err.Error()
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	if req.X < 1 || req.Y < 1 || req.Height < 1 || req.Height > 30 {
		errResponse.Message = "x and y must be greatest equal 1, height must be between 1 and 30"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	estate, err := s.Repository.FindEstate(ctx, &repository.FilterEstate{ID: estateID})
	if err != nil {
		errResponse.Message = fmt.Sprintf("estate %s not found", estateID)
		return c.JSON(http.StatusNotFound, errResponse)
	}

	if estate.Length < req.X || estate.Width < req.Y {
		errResponse.Message = "coordinate out of bound"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	_, err = s.Repository.FindEstateTree(ctx, &repository.FilterEstateTree{
		EstateID: estateID,
		X:        req.X,
		Y:        req.Y,
	})
	if err == nil {
		errResponse.Message = "plot already has tree"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	data := repository.EstateTree{
		BaseModel: repository.BaseModel{
			ID: uuid.NewString(),
		},
		EstateID: estateID,
		X:        req.X,
		Y:        req.Y,
		Height:   req.Height,
	}
	if err = s.Repository.CreateEstateTree(ctx, &data); err != nil {
		errResponse.Message = "failed to create estate tree"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	return c.JSON(http.StatusCreated, generated.EstateTreeResponse{
		Id: data.ID,
	})
}

func (s *Server) GetEstateIdStats(c echo.Context, estateID string) error {
	ctx := c.Request().Context()
	var errResponse generated.ErrorResponse
	_, err := s.Repository.FindEstate(ctx, &repository.FilterEstate{ID: estateID})
	if err != nil {
		errResponse.Message = fmt.Sprintf("estate %s not found", estateID)
		return c.JSON(http.StatusNotFound, errResponse)
	}

	countEstateTree := s.Repository.CountEstateTree(ctx, &repository.FilterEstateTree{EstateID: estateID})
	var stats repository.EstateTreeStats
	if countEstateTree > 0 {
		stats, err = s.Repository.GetEstateTreeStats(ctx, &repository.FilterEstateTree{EstateID: estateID})
		if err != nil {
			errResponse.Message = err.Error()
			return c.JSON(http.StatusBadRequest, errResponse)
		}
	}

	return c.JSON(http.StatusOK, generated.EstateStatsResponse{
		Count:  countEstateTree,
		Max:    stats.Max,
		Min:    stats.Min,
		Median: stats.Median,
	})
}

func (s *Server) GetEstateIdDronePlan(c echo.Context, estateID string, params generated.GetEstateIdDronePlanParams) error {
	ctx := c.Request().Context()
	var errResponse generated.ErrorResponse
	estate, err := s.Repository.FindEstate(ctx, &repository.FilterEstate{ID: estateID})
	if err != nil {
		errResponse.Message = fmt.Sprintf("estate %s not found", estateID)
		return c.JSON(http.StatusNotFound, errResponse)
	}

	filterEstateTree := repository.FilterEstateTree{
		Filter: repository.Filter{
			Page:    1,
			ShowAll: true,
		},
		EstateID: estateID,
	}
	estateTree, err := s.Repository.FindAllMapEstateTree(ctx, &filterEstateTree)
	if err != nil {
		errResponse.Message = err.Error()
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	var totalDistance, currentHeight int
	for y := 1; y <= estate.Length; y++ {
		var xStart, xEnd, xStep int
		if y%2 == 1 {
			xStart, xEnd, xStep = 1, estate.Width, 1
		} else {
			xStart, xEnd, xStep = estate.Width, 1, -1
		}
		for x := xStart; x != xEnd+xStep; x += xStep {
			tree := estateTree[repository.CoordinatePoint{X: x, Y: y}]
			targetHeight := tree.Height + 1
			totalDistance += int(math.Abs(float64(targetHeight - currentHeight)))
			currentHeight = targetHeight

			if params.MaxDistance != nil && totalDistance > *params.MaxDistance {
				return c.JSON(http.StatusOK, setResponseMaxDistance(totalDistance, x, y))
			}

			if x != xEnd {
				totalDistance += 10
				if params.MaxDistance != nil && totalDistance > *params.MaxDistance {
					return c.JSON(http.StatusOK, setResponseMaxDistance(totalDistance, x, y))
				}
			}
		}

		if y != estate.Length {
			totalDistance += 10
			if params.MaxDistance != nil && totalDistance > *params.MaxDistance {
				return c.JSON(http.StatusOK, setResponseMaxDistance(totalDistance, xStart, y+1))
			}
		}
	}

	totalDistance += currentHeight
	if params.MaxDistance != nil {
		return c.JSON(http.StatusOK, setResponseMaxDistance(totalDistance, estate.Width, estate.Length))
	}
	return c.JSON(http.StatusOK, generated.EstateDronePlanResponse{Distance: totalDistance})
}

func setResponseMaxDistance(distance, x, y int) generated.EstateDronePlanResponse {
	return generated.EstateDronePlanResponse{
		Distance: distance,
		Rest: &struct {
			X int `json:"x"`
			Y int `json:"y"`
		}{
			X: x,
			Y: y,
		},
	}
}
