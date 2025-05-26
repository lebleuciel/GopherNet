package api

import (
	"gophernet/pkg/app"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IGopherController interface {
	GetGopher(c *gin.Context)
	RentBurrow(c *gin.Context)
	ReleaseBurrow(c *gin.Context)
	GetBurrowStatus(c *gin.Context)
}

type GopherController struct {
	gopherApp app.IGopherApp
}

func NewGopherController(gopherApp app.IGopherApp) *GopherController {
	return &GopherController{
		gopherApp: gopherApp,
	}
}

func (g *GopherController) GetGopher(c *gin.Context) {
	g.gopherApp.GetGopher()
	c.String(http.StatusOK, "Gopher says hi!")
}

func (g *GopherController) RentBurrow(c *gin.Context) {
	burrowIDStr := c.Param("id")
	burrowID, err := strconv.Atoi(burrowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid burrow ID"})
		return
	}

	burrow, err := g.gopherApp.RentBurrow(c.Request.Context(), burrowID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, burrow)
}

func (g *GopherController) ReleaseBurrow(c *gin.Context) {
	burrowIDStr := c.Param("id")
	burrowID, err := strconv.Atoi(burrowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid burrow ID"})
		return
	}

	burrow, err := g.gopherApp.ReleaseBurrow(c.Request.Context(), burrowID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, burrow)
}

func (g *GopherController) GetBurrowStatus(c *gin.Context) {
	burrows, err := g.gopherApp.GetBurrowStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get burrow status"})
		return
	}
	c.JSON(http.StatusOK, burrows)
}
