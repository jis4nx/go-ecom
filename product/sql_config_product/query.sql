-- name: InsertProduct :one
INSERT INTO product (name, description, price, available)
  VALUES ($1, $2, $3, $4)
  RETURNING *;


-- name: ListProduct :many
SELECT * FROM product;


-- name: GetProduct :one
SELECT * FROM product 
  WHERE id=$1;

-- name: UpdateProduct :one
UPDATE product
SET name = $1, description = $2, price = $3, available = $4, version = version + 1
  WHERE id = $5 AND version = $6
RETURNING version;


