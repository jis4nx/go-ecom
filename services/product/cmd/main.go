package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/jis4nx/go-ecom/config"
	"github.com/jis4nx/go-ecom/helpers"
	"github.com/jis4nx/go-ecom/helpers/types"
	"github.com/jis4nx/go-ecom/pkg/logger"
	"github.com/jis4nx/go-ecom/services/product/api"
	"github.com/jis4nx/go-ecom/services/product/internals/rabbit/consumers"
	"github.com/jis4nx/go-ecom/services/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
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

	// Load and Sets the Config globally and passing the the Local app
	cfg := config.LoadConfig(envVars)
	config.SetGlobalConfig(&cfg)
	app := helpers.NewApp(cfg)

	serviceInfo := types.NewServiceInfo(
		"product",
		app.Cfg.Services.ProductServer.HOST,
		app.Cfg.Services.ProductServer.PORT,
	)
	userLog := logger.Logger{}
	userLog.SetLogFile(base, "gocom.log")
	userLog.InitLogger(serviceInfo)
	app.Logger = &userLog

	router := productApp.LoadRoutes()
	app.AddRoutes(router)

	utils.SetGlobalProductApp(app)
	productApp.SetProductApp(app)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go consumers.StartConsumer()
	err = app.Start(ctx)
	if err != nil {
		app.Logger.Fatal("Failed to start the server", zap.Error(err))
	}
}
