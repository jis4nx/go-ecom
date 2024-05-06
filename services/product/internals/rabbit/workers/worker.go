package workers

import (
	"context"
	"encoding/json"

	"github.com/jis4nx/go-ecom/helpers"
	"github.com/jis4nx/go-ecom/pkg/rabbit"
	"github.com/jis4nx/go-ecom/services/product/internals/productmodel"
	"github.com/jis4nx/go-ecom/services/product/internals/rabbit/publisher"
	"github.com/jis4nx/go-ecom/services/utils"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type ProductWorker struct {
	Queue      string
	Consumer   string
	HandleFunc func(context.Context, amqp091.Delivery, *helpers.App, *productmodel.Queries) error
}

func NewProductWorker(queue, consumer string, handleFunc func(context.Context, amqp091.Delivery, *helpers.App, *productmodel.Queries) error) *ProductWorker {
	return &ProductWorker{
		Queue:      queue,
		Consumer:   consumer,
		HandleFunc: handleFunc,
	}
}

func getProductQueue(eventType publisher.ProductEventType) (string, string) {
	switch eventType {
	case publisher.ProductCreated:
		return "product_created", "createConsumer"
	case publisher.ProductUpdated:
		return "product_updated", "updateConsumer"
	default:
		panic("invalid event type")
	}
}

func (w *ProductWorker) Start(rc *rabbit.RabbitClient, g *errgroup.Group) {
	app := utils.GetProductApp()
	db := app.PGDB
	query := productmodel.New(db)

	messageBus, err := rc.Consume(w.Queue, w.Consumer, false)
	if err != nil {
		app.Logger.Error("Failed to consume message", zap.Error(err))
		return
	}

	for message := range messageBus {
		msg := message
		g.Go(func() error {
			return w.HandleFunc(context.Background(), msg, app, query)
		})
	}
}

func ProductWorkerFactory(eventType publisher.ProductEventType) *ProductWorker {
	queue, consumer := getProductQueue(eventType)
	var handleFunc func(context.Context, amqp091.Delivery, *helpers.App, *productmodel.Queries) error

	switch eventType {
	case publisher.ProductCreated:
		handleFunc = CreateProductWorker
	case publisher.ProductUpdated:
		handleFunc = UpdateProductWorker
	default:
		panic("invalid event type")
	}

	return NewProductWorker(queue, consumer, handleFunc)
}

func CreateProductWorker(ctx context.Context, msg amqp091.Delivery, app *helpers.App, query *productmodel.Queries) error {
	body := msg.Body
	p := productmodel.Product{}
	err := json.Unmarshal(body, &p)
	if err != nil {
		app.Logger.Error("Failed to Unmarshal product", zap.Error(err))
		return err
	}
	_, err = query.InsertProduct(ctx, productmodel.InsertProductParams{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Sku:         p.Sku,
		Available:   p.Available,
	})
	if err != nil {
		app.Logger.Error("Failed to Update product", zap.Error(err))
		return err
	}
	return nil
}

func UpdateProductWorker(ctx context.Context, msg amqp091.Delivery, app *helpers.App, query *productmodel.Queries) error {
	body := msg.Body
	p := productmodel.Product{}
	err := json.Unmarshal(body, &p)
	if err != nil {
		app.Logger.Error("Failed to Unmarshal product", zap.Error(err))
		return err
	}
	_, err = query.UpdateProduct(ctx, productmodel.UpdateProductParams{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Available:   p.Available,
		Sku:         p.Sku,
	})
	if err != nil {
		app.Logger.Error("Failed to Update product", zap.Error(err))
		return err
	}
	return nil
}
