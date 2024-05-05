package api

import (
	"context"
	"net/http"

	"github.com/jis4nx/go-ecom/helpers"
	"github.com/jis4nx/go-ecom/services/product/internals/productmodel"
	"github.com/jis4nx/go-ecom/services/product/internals/rabbit/publisher"
	"go.uber.org/zap"
)

type ProductApp struct {
	*helpers.App
}

func (p *ProductApp) SetProductApp(a *helpers.App) {
	p.App = a
}

func (a *ProductApp) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	input := productmodel.Product{}

	if err := a.ReadJson(w, r, &input); err != nil {
		a.Logger.Warn("Failed to Read json", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	publisher.DispatchProductEvent(ctx, a.Rabbit, publisher.ProductCreated, input)
}
