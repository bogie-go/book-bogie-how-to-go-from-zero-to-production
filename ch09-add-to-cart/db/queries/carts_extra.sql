-- name: UpsertCartByToken :one
INSERT INTO carts (token) VALUES ($1)
ON CONFLICT (token) DO UPDATE SET token = EXCLUDED.token
RETURNING *;
