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
	// Create new server
	config := config.LoadConfigFromDefaultPath()
	db := db.NewDatabase(context.Background(), &config.Database)

	burrowRepo := repo.NewBurrowRepository(db)
	gopherApp := app.NewGopherApp(burrowRepo)
	api := controller.NewGopherController(gopherApp)

	server := server.NewServer(api)
	server.ServeHTTP()

	log.Println("Server exiting")
}
