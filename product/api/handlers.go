package api

import (
	"fmt"
	"net/http"

	"github.com/jis4nx/go-ecom/helpers"
)

type ProductApp struct{
  *helpers.App
}

func (p *ProductApp) SetProductApp(a *helpers.App){
  p.App = a
}

func (a *ProductApp) CreateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println(a.PGDB)
}
