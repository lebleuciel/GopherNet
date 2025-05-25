package main

import (
	"context"
	"log"

	"gophernet/pkg/app"
	"gophernet/pkg/config"
	controller "gophernet/pkg/controller"
	"gophernet/pkg/db"
	"gophernet/pkg/repo"
	"gophernet/server"
)

func main() {
	log.Println("Starting GopherNet server...")

	// Create new server
	log.Println("Loading configuration...")
	config := config.LoadConfigFromDefaultPath()

	log.Println("Initializing database connection...")
	db := db.NewDatabase(context.Background(), &config.Database)

	log.Println("Setting up repositories and controllers...")
	burrowRepo := repo.NewBurrowRepository(db)
	gopherApp := app.NewGopherApp(burrowRepo)
	api := controller.NewGopherController(gopherApp)

	log.Println("Starting HTTP server...")
	server := server.NewServer(api)
	server.ServeHTTP()

	log.Println("Server exiting")
}
