package consumers

import (
	"context"
	"time"

	"github.com/jis4nx/go-ecom/pkg/rabbit"
	"github.com/jis4nx/go-ecom/services/product/internals/rabbit/publisher"
	"github.com/jis4nx/go-ecom/services/product/internals/rabbit/workers"
	"github.com/jis4nx/go-ecom/services/utils"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func StartConsumer() {
	app := utils.GetProductApp()
	conn, err := rabbit.ConnectRabbitMQ("guest", "guest", "rabbit:5672", "/")
	if err != nil {
		app.Logger.Fatal("Failed to connect rabbitmq", zap.Error(err))
	}

	client, err := rabbit.NewRabbitClient(conn)
	if err != nil {
		app.Logger.Fatal("Failed to create rabbitmq channel", zap.Error(err))
	}

	// Create Exchange
	if err = client.CreateExchange("product_events", "topic", true, false); err != nil {
		app.Logger.Fatal("Failed to create rabbitmq exchange", zap.Error(err))
	}

	// Decalrig New Queue
	if err = client.NewQueue("product_created", true, false); err != nil {
		app.Logger.Fatal("Failed to create rabbitmq queue", zap.Error(err))
	}

	// Create Binding with Queue & Routing key
	if err = client.CreateBinding("product_created", "product.created.*", "product_events"); err != nil {
		app.Logger.Fatal("Failed to create rabbitmq queue", zap.Error(err))
	}

	var blocking chan struct{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(10)

	consumer, err := rabbit.NewConsumer()
	if err != nil {
		app.Logger.Error("Failed to create Consumer", zap.Error(err))
	}

  // starting the Product Creation worker
	createWorker := workers.ProductWorkerFactory(publisher.ProductCreated)
	createWorker.Start(&consumer, g)

  // starting the Product Update worker
	updateWorker := workers.ProductWorkerFactory(publisher.ProductUpdated)
	updateWorker.Start(&consumer, g)

	// Wait for all workers to finish
	if err := g.Wait(); err != nil {
		app.Logger.Error("Error in worker", zap.Error(err))
	}

	defer cancel()
	<-blocking
}
