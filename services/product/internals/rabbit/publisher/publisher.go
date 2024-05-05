package publisher

import (
	"context"

	"github.com/jis4nx/go-ecom/pkg/rabbit"
	"github.com/jis4nx/go-ecom/services/product/internals/productmodel"
	"github.com/jis4nx/go-ecom/services/utils"
	"go.uber.org/zap"
)

func DispatchProductEvent(ctx context.Context, rc *rabbit.RabbitClient, eventType ProductEventType, p productmodel.Product) {
	var event string
	switch eventType {
	case ProductCreated:
		event = "product.created.key1"
	case ProductUpdated:
		event = "product.updated.key2"
	}

	app := utils.GetProductApp()
	if err := rc.PublishMsgWithContext(ctx, "product_events", event, p); err != nil {
		app.Logger.Info("Failed to publish product", zap.String("product", p.Name))
	}
}
