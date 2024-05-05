// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package productmodel

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Sku         sql.NullString `json:"sku"`
	Price       float64        `json:"price"`
	Available   sql.NullBool   `json:"available"`
	CreatedAt   time.Time      `json:"created_at"`
	Version     int32          `json:"version"`
}