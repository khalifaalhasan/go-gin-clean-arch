-- name: CreateOrder :one
INSERT INTO orders (
    user_id, product_id, amount
) VALUES (
    $1, $2, $3
) RETURNING id, user_id, product_id, amount, created_at;



-- name: GetIdempotencyKey :one
SELECT * FROM idempotency_keys
WHERE key= $1 AND user_id = $2
FOR UPDATE;

-- name: CreateIdempotencyKey :one
INSERT INTO idempotency_keys (
    key, user_id, response_code, response_body
) VALUES (
    $1, $2, $3, $4
) RETURNING *;
