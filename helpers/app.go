package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jis4nx/go-ecom/config"
	"github.com/jis4nx/go-ecom/pkg/logger"
	"github.com/jis4nx/go-ecom/pkg/rabbit"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type App struct {
	Router http.Handler
	Cfg    config.Config
	PGDB   *sql.DB
	Rabbit *rabbit.RabbitClient
  Logger *logger.Logger
}

// Wrapper Function to connect to Postgres DB
func ConnectPG(host, user, password, dbName, port string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, password)
	pgdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect postgres server", err.Error())
	}

	if err = pgdb.Ping(); err != nil {
		log.Fatal("Failed to  ping the server", err.Error())
	}
	return pgdb
}

func NewApp(c config.Config) *App {
	pgdb := ConnectPG(c.DB.DBHOST, c.DB.DBUSER, c.DB.DBPASS, c.DB.DBNAME, c.DB.DBPORT)

	conn, err := rabbit.ConnectRabbitMQ(c.RQ.USER, c.RQ.PASSWORD, c.RQ.HOST, c.RQ.VHOST)
	if err != nil {
		log.Fatal("Failed to connect Rabbitmq server", err.Error())
	}
	client, err := rabbit.NewRabbitClient(conn)
	if err != nil {
		log.Fatal("Failed to create Rabbitmq channel", err.Error())
	}
	app := &App{
		Cfg:    c,
		PGDB:   pgdb,
		Rabbit: &client,
	}

	return app
}

func (app *App) AddRoutes(router chi.Router) {
	app.Router = router
}

// Starts the Respected app server
func (app *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "product", app.Cfg.Services.ProductServer.PORT),
		Handler: app.Router,
	}

  app.Logger.Info("Server Started")

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
      app.Logger.Fatal("Failed to Start server", zap.Error(err))
		}
	}()

	select {
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
    app.Logger.Info("Received Interrupt Signal, Closing Server")
		defer cancel()

		return server.Shutdown(timeout)
	}
}
