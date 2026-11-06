-- name: UpsertCartByToken :one
INSERT INTO carts (token) VALUES ($1)
ON CONFLICT (token) DO UPDATE SET token = EXCLUDED.token
RETURNING *;

-- name: GetCartByToken :one
SELECT * FROM carts WHERE token = $1;

-- name: LockCart :one
SELECT id FROM carts WHERE id = $1 FOR UPDATE;
