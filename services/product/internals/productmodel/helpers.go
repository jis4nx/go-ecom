package productmodel

import (
	"database/sql"

	"github.com/jis4nx/go-ecom/services/product/types"
)

func NewProduct(inp types.ProductInput) Product {
	// Function to create sql.NullString from a string
	nullString := func(s string) sql.NullString {
		if s == "" {
			return sql.NullString{Valid: false}
		}
		return sql.NullString{String: s, Valid: true}
	}

	// Function to create sql.NullBool from a bool
	nullBool := func(b bool) sql.NullBool {
		return sql.NullBool{Bool: b, Valid: true}
	}

	// Set fields based on input
	fDesc := nullString(inp.Description)
	fSku := nullString(inp.Sku)
	fAvailable := nullBool(inp.Available)

	// Create and return the Product
	return Product{
		Name:        inp.Name,
		Price:       inp.Price,
		Description: fDesc,
		Sku:         fSku,
		Available:   fAvailable,
	}
}
