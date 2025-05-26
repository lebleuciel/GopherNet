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
	"gophernet/server"
)

func main() {
	log.Println("Starting GopherNet server...")

	// Initialize database
	database := db.NewDatabase(context.Background(), &config.DefaultDatabase)
	defer database.Close()

	// Initialize repository
	burrowRepo := repo.NewBurrowRepository(database)

	// Initialize app
	gopherApp := app.NewGopherApp(burrowRepo)

	// Create context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the scheduler
	// Start the scheduler
	gopherApp.StartScheduler(ctx) // Wait for interrupt signal
	<-ctx.Done()
	log.Println("Shutting down...")

	// Stop the scheduler
	gopherApp.StopScheduler()

	log.Println("Starting HTTP server...") ///////
	server := server.NewServer(controller.NewGopherController(gopherApp))
	server.ServeHTTP()

	log.Println("Server exiting")
}
