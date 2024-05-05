package models

type ProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Sku         string  `json:"sku"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
}
