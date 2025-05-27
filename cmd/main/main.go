package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gophernet/pkg/app"
	"gophernet/pkg/config"
	controller "gophernet/pkg/controller"
	"gophernet/pkg/db"
	"gophernet/pkg/repo"
	"gophernet/pkg/shutdown"
	"gophernet/server"
)

func main() {
	log.Println("Starting GopherNet server...")
	bgCtx := context.Background()
	// Load configuration
	cfg := config.LoadConfigFromDefaultPath()

	// Initialize database
	database := db.NewDatabase(bgCtx, &cfg.Database)
	defer database.Close()

	// Check if database is initialized
	initialized, err := database.IsInitialized(bgCtx)
	if err != nil {
		log.Printf("Error checking database initialization: %v", err)
		os.Exit(1)
	}

	if !initialized {
		log.Println("Database not initialized. Please run database migrations first.")
		os.Exit(1)
	}

	// Initialize repository
	burrowRepo := repo.NewBurrowRepository(database)

	// Initialize app
	gopherApp := app.NewGopherApp(burrowRepo)
	scheduler := app.NewScheduler(burrowRepo, &cfg.Scheduler)
	scheduler.Start(bgCtx)
	shutdown.GetManager().Register("scheduler", func(ctx context.Context) error {
		scheduler.Stop()
		return nil
	})

	defer shutdown.GetManager().Shutdown(context.Background())
	// Create context that listens for the interrupt signal
	bgCtx, stop := signal.NotifyContext(bgCtx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize and start HTTP server
	server := server.NewServer(controller.NewGopherController(gopherApp))
	go server.ServeHTTP()

	// Wait for interrupt signal
	<-bgCtx.Done()
	log.Println("Shutting down...")

	log.Println("Server exiting")
}
