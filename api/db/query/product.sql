-- name: CreateProduct :one
INSERT INTO products (
  name, 
  description, 
  price, 
  stock
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateProductStock :one
UPDATE products
SET stock = stock - sqlc.arg(amount) -- <--- Pastikan ini MINUS (-)
WHERE id = sqlc.arg(id) AND stock >= sqlc.arg(amount)
RETURNING *;