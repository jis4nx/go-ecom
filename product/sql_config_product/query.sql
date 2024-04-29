-- name: InsertProduct :one
INSERT INTO product (name, description, price, available)
  VALUES ($1, $2, $3, $4)
  RETURNING *;


-- name: ListProduct :many
SELECT * FROM product;


-- name: GetProduct :one
SELECT * FROM product 
  WHERE id=$1;
