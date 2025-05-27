package controller

import (
	"net/http"
	"strconv"

	"gophernet/pkg/app"
	"gophernet/pkg/dto"
	"gophernet/pkg/errors"
	"gophernet/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title GopherNet API
// @version 1.0
// @description This is the GopherNet API server.
// @host localhost:8080
// @BasePath /api/v1

type IGopherController interface {
	RentBurrow(c *gin.Context)
	ReleaseBurrow(c *gin.Context)
	GetBurrowStatus(c *gin.Context)
	GetBurrow(c *gin.Context)
}

type GopherController struct {
	gopherApp *app.GopherApp
	log       *zap.Logger
}

func NewGopherController(gopherApp *app.GopherApp) *GopherController {
	return &GopherController{
		gopherApp: gopherApp,
		log:       logger.Get(),
	}
}

func (g *GopherController) handleError(c *gin.Context, err error) {
	g.log.Error("Request error", zap.Error(err))

	var statusCode int
	var message string

	switch err {
	case errors.ErrBurrowNotFound:
		statusCode = http.StatusNotFound
		message = "Burrow not found"
	case errors.ErrBurrowOccupied:
		statusCode = http.StatusConflict
		message = "Burrow is already occupied"
	case errors.ErrBurrowNotOccupied:
		statusCode = http.StatusConflict
		message = "Burrow is not occupied"
	case errors.ErrInvalidBurrowID:
		statusCode = http.StatusBadRequest
		message = "Invalid burrow ID"
	default:
		statusCode = http.StatusInternalServerError
		message = "Internal server error"
	}

	c.JSON(statusCode, dto.ErrorResponse{Error: message})
}

// @Summary Get a Burrow
// @Description Get a burrow by ID
// @Tags burrows
// @Accept json
// @Produce json
// @Param id path int true "Burrow ID"
// @Success 200 {object} dto.BurrowResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /burrows/{id} [get]
func (g *GopherController) GetBurrow(c *gin.Context) {
	burrowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		g.handleError(c, errors.ErrInvalidBurrowID)
		return
	}

	burrow, err := g.gopherApp.GetBurrow(c.Request.Context(), burrowID)
	if err != nil {
		g.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.BurrowResponse{
		ID:         burrow.ID,
		Name:       burrow.Name,
		Depth:      burrow.Depth,
		Width:      burrow.Width,
		IsOccupied: burrow.IsOccupied,
		Age:        burrow.Age,
	})
}

// @Summary Rent a Burrow
// @Description Rent a burrow by ID
// @Tags burrows
// @Accept json
// @Produce json
// @Param id path int true "Burrow ID"
// @Success 200 {object} dto.BurrowResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /burrows/{id}/rent [post]
func (g *GopherController) RentBurrow(c *gin.Context) {
	burrowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		g.handleError(c, errors.ErrInvalidBurrowID)
		return
	}

	burrow, err := g.gopherApp.RentBurrow(c.Request.Context(), burrowID)
	if err != nil {
		g.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.BurrowResponse{
		ID:         burrow.ID,
		Name:       burrow.Name,
		Depth:      burrow.Depth,
		Width:      burrow.Width,
		IsOccupied: burrow.IsOccupied,
		Age:        burrow.Age,
	})
}

// @Summary Release a Burrow
// @Description Release a burrow by ID
// @Tags burrows
// @Accept json
// @Produce json
// @Param id path int true "Burrow ID"
// @Success 200 {object} dto.BurrowResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /burrows/{id}/release [post]
func (g *GopherController) ReleaseBurrow(c *gin.Context) {
	burrowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		g.handleError(c, errors.ErrInvalidBurrowID)
		return
	}

	burrow, err := g.gopherApp.ReleaseBurrow(c.Request.Context(), burrowID)
	if err != nil {
		g.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.BurrowResponse{
		ID:         burrow.ID,
		Name:       burrow.Name,
		Depth:      burrow.Depth,
		Width:      burrow.Width,
		IsOccupied: burrow.IsOccupied,
		Age:        burrow.Age,
	})
}

// @Summary Get Burrow Status
// @Description Get the status of all burrows
// @Tags burrows
// @Accept json
// @Produce json
// @Success 200 {array} dto.BurrowResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /burrows/status [get]
func (g *GopherController) GetBurrowStatus(c *gin.Context) {
	burrows, err := g.gopherApp.GetBurrowStatus(c.Request.Context())
	if err != nil {
		g.handleError(c, err)
		return
	}

	var responseBurrows []dto.BurrowResponse
	for _, burrow := range burrows {
		responseBurrows = append(responseBurrows, dto.BurrowResponse{
			ID:         burrow.ID,
			Name:       burrow.Name,
			Depth:      burrow.Depth,
			Width:      burrow.Width,
			IsOccupied: burrow.IsOccupied,
			Age:        burrow.Age,
		})
	}
	c.JSON(http.StatusOK, responseBurrows)
}
