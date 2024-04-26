package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/jis4nx/go-ecom/config"
	"github.com/jis4nx/go-ecom/helpers"
	"github.com/jis4nx/go-ecom/product/api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	envVars, err := godotenv.Read("dev.env")
	if err != nil {
		log.Fatal("Failed to load environment variables")
	}

	var prodcutApp api.ProductApp

	cfg := config.LoadConfig(envVars)
	app := helpers.NewApp(cfg)

	router := prodcutApp.LoadRoutes()
	app.AddRoutes(router)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = app.Start(ctx)
	if err != nil {
		log.Fatal("Failed to close the server", err.Error())
	}
}
