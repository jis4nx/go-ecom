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
	"github.com/jis4nx/go-ecom/helpers/rabbit"
)

type App struct {
	Router http.Handler
	Cfg    config.Config
	PGDB   *sql.DB
	Rabbit *rabbit.RabbitClient
}

func NewApp(c config.Config) *App {
	dbString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DB.DBUSER, c.DB.DBPASS, c.DB.DBHOST, c.DB.DBPORT, c.DB.DBNAME)

	conn, err := rabbit.ConnectRabbitMQ(c.RQ.USER, c.RQ.PASSWORD, c.RQ.HOST, c.RQ.VHOST)
	if err != nil {
		log.Fatal("Failed to connect Rabbitmq server", err.Error())
	}
	client, err := rabbit.NewRabbitClient(conn)
	if err != nil {
		log.Fatal("Failed to create Rabbitmq channel", err.Error())
	}
	pgdb, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatal("Failed to connect postgres server", err.Error())
	}

	if err = pgdb.Ping(); err != nil {
		log.Fatal("Failed to  ping the server", err.Error())
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

func (app *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", app.Cfg.Services.ProductServer.HOST, app.Cfg.Services.ProductServer.PORT),
		Handler: app.Router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	select {
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		fmt.Printf("Received Interrupt Signal, Closing Server")
		defer cancel()

		return server.Shutdown(timeout)
	}
}