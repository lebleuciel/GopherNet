package api

import (
	"gophernet/pkg/app"

	"github.com/gin-gonic/gin"
)

type IGopherController interface {
	GetGopher(c *gin.Context)
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
}
