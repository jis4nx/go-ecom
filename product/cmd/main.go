package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/jis4nx/go-ecom/config"
	"github.com/jis4nx/go-ecom/helpers"
	"github.com/jis4nx/go-ecom/product/api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	base, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to retrieve current dir", err.Error())
	}
	abs := filepath.Join(base, "dev.env")

	envVars, err := godotenv.Read(abs)
	if err != nil {
		log.Fatal("Failed to load environment variables")
	}

	var productApp api.ProductApp

	cfg := config.LoadConfig(envVars)
	app := helpers.NewApp(cfg)

	router := productApp.LoadRoutes()
	app.AddRoutes(router)

  productApp.SetProductApp(app)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = app.Start(ctx)
	if err != nil {
		log.Fatal("Failed to close the server", err.Error())
	}
}
