package utils

import "github.com/jis4nx/go-ecom/helpers"

var productApp *helpers.App

func SetGlobalProductApp(a *helpers.App) {
	productApp = a
}

func GetProductApp() *helpers.App {
	return productApp
}
