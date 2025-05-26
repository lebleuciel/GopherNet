package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	docs "gophernet/docs"
	controller "gophernet/pkg/controller"
	"gophernet/pkg/shutdown"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	// For development: set Gin to debug mode to catch errors
	// Use gin.ReleaseMode in production
	gin.SetMode(gin.DebugMode)

	engine := gin.Default()

	// Swagger documentation settings
	docs.SwaggerInfo.Title = "GopherNet API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	// CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	// Recovery middleware
	engine.Use(gin.Recovery())

	return &Server{
		environment: "debug",
		engine:      engine,
		handler:     handler,
	}
}

// registerRoutes sets up all HTTP routes including Swagger and API endpoints
func (s *Server) registerRoutes() {
	// Swagger route (not under /api/v1)
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := s.engine.Group("/api/v1")
	{
		v1.GET("/gopher", s.handler.GetGopher)

		burrowRoutes := v1.Group("/burrows")
		{
			burrowRoutes.POST("/:id/rent", s.handler.RentBurrow)
			burrowRoutes.POST("/:id/release", s.handler.ReleaseBurrow)
			burrowRoutes.GET("/status", s.handler.GetBurrowStatus)
		}
	}
}

// ServeHTTP starts the HTTP server
func (s *Server) ServeHTTP() {
	s.registerRoutes()

	s.srv = &http.Server{
		Addr:    ":8080",
		Handler: s.engine,
	}

	// Register graceful shutdown handler
	shutdown.GetManager().Register("http-server", func(ctx context.Context) error {
		return s.srv.Shutdown(ctx)
	})

	// Log where Swagger UI is available
	fmt.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")

	// Start server
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("HTTP server failed: %v", err))
	}
}
