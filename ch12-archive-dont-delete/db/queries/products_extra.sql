-- name: GetProductByName :one
SELECT * FROM products WHERE name = $1;

-- name: TakeStock :execrows
UPDATE products
SET stock = stock - sqlc.arg(quantity), updated_at = now()
WHERE id = sqlc.arg(id) AND stock >= sqlc.arg(quantity) AND archived_at IS NULL;

-- name: ArchiveProduct :one
UPDATE products SET archived_at = COALESCE(archived_at, now()), updated_at = now() WHERE id = $1 RETURNING *;

-- name: RestoreProduct :one
UPDATE products SET archived_at = NULL, updated_at = now() WHERE id = $1 RETURNING *;

-- name: GetProductForUpdate :one
SELECT * FROM products WHERE id = $1 FOR UPDATE;

-- name: ListCatalogProducts :many
SELECT * FROM products WHERE archived_at IS NULL
ORDER BY created_at DESC, id DESC LIMIT $1;

-- name: GetCatalogProduct :one
SELECT * FROM products WHERE id = $1 AND archived_at IS NULL;
