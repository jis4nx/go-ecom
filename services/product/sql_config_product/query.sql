-- name: InsertProduct :one
INSERT INTO product (name, description,sku, price, available)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *;


-- name: ListProduct :many
SELECT * FROM product;


-- name: GetProduct :one
SELECT * FROM product 
  WHERE id=$1;

-- name: UpdateProduct :one
UPDATE product
SET name = $1, description = $2, sku = $3, price = $4, available = $5, version = version + 1
  WHERE id = $6 AND version = $7
RETURNING version;


