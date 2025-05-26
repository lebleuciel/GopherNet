package controller

import (
	"net/http"
	"strconv"

	"gophernet/pkg/app"
	"gophernet/pkg/models"

	"github.com/gin-gonic/gin"
)

// @title GopherNet API
// @version 1.0
// @description This is the GopherNet API server.
// @host localhost:8080
// @BasePath /api/v1

type IGopherController interface {
	GetGopher(c *gin.Context)
	RentBurrow(c *gin.Context)
	ReleaseBurrow(c *gin.Context)
	GetBurrowStatus(c *gin.Context)
}

type GopherController struct {
	gopherApp *app.GopherApp
}

func NewGopherController(gopherApp *app.GopherApp) *GopherController {
	return &GopherController{
		gopherApp: gopherApp,
	}
}

// @Summary Get Gopher
// @Description Get a welcome message from Gopher
// @Tags gopher
// @Accept json
// @Produce json
// @Success 200 {string} string "Gopher says hi!"
// @Router /gopher [get]
func (g *GopherController) GetGopher(c *gin.Context) {
	g.gopherApp.GetGopher(c.Request.Context())
	c.String(http.StatusOK, "Gopher says hi!")
}

// @Summary Rent a Burrow
// @Description Rent a burrow by ID
// @Tags burrows
// @Accept json
// @Produce json
// @Param id path int true "Burrow ID"
// @Success 200 {object} models.Burrow
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /burrows/{id}/rent [post]
func (g *GopherController) RentBurrow(c *gin.Context) {
	burrowIDStr := c.Param("id")
	burrowID, err := strconv.Atoi(burrowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid burrow ID"})
		return
	}

	burrow, err := g.gopherApp.RentBurrow(c.Request.Context(), burrowID)
	if err != nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Burrow{
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
// @Success 200 {object} models.Burrow
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /burrows/{id}/release [post]
func (g *GopherController) ReleaseBurrow(c *gin.Context) {
	burrowIDStr := c.Param("id")
	burrowID, err := strconv.Atoi(burrowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid burrow ID"})
		return
	}

	burrow, err := g.gopherApp.ReleaseBurrow(c.Request.Context(), burrowID)
	if err != nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Burrow{
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
// @Success 200 {array} models.Burrow
// @Failure 500 {object} models.ErrorResponse
// @Router /burrows/status [get]
func (g *GopherController) GetBurrowStatus(c *gin.Context) {
	burrows, err := g.gopherApp.GetBurrowStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get burrow status"})
		return
	}

	var responseBurrows []models.Burrow
	for _, burrow := range burrows {
		responseBurrows = append(responseBurrows, models.Burrow{
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
