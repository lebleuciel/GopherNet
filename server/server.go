package server

import (
	"context"
	"net/http"
	"time"

	controller "gophernet/pkg/controller"
	"gophernet/pkg/shutdown"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	environment string
	engine      *gin.Engine
	handler     controller.IGopherController
	srv         *http.Server
}

// NewServer creates a new Server instance
func NewServer(handler controller.IGopherController) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	// Add CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	// Add recovery middleware
	engine.Use(gin.Recovery())

	return &Server{
		environment: "release",
		engine:      engine,
		handler:     handler,
	}
}

func (s *Server) registerRoutes() {
	v1 := s.engine.Group("/api/v1")
	{
		v1.GET("/gopher", s.handler.GetGopher)

		// Burrow routes
		burrowRoutes := v1.Group("/burrows")
		{
			burrowRoutes.POST("/:id/rent", s.handler.RentBurrow)
			burrowRoutes.POST("/:id/release", s.handler.ReleaseBurrow)
			burrowRoutes.GET("/status", s.handler.GetBurrowStatus)
		}
	}
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP() {
	s.registerRoutes()
	s.srv = &http.Server{
		Addr:    ":8080",
		Handler: s.engine,
	}

	// Register shutdown handler
	shutdown.GetManager().Register("http-server", func(ctx context.Context) error {
		return s.srv.Shutdown(ctx)
	})
	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
