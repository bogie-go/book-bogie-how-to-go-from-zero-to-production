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

-- The WHERE of ListCatalogProducts and CountCatalogProducts must stay
-- identical, or the total won't match the pages.

-- name: ListCatalogProducts :many
SELECT * FROM products
WHERE archived_at IS NULL
  AND (sqlc.narg(category_id)::bigint IS NULL OR EXISTS (
        SELECT 1 FROM product_categories
        WHERE product_categories.product_id = products.id
          AND product_categories.category_id = sqlc.narg(category_id)))
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(row_limit) OFFSET sqlc.arg(row_offset);

-- name: CountCatalogProducts :one
SELECT count(*) FROM products
WHERE archived_at IS NULL
  AND (sqlc.narg(category_id)::bigint IS NULL OR EXISTS (
        SELECT 1 FROM product_categories
        WHERE product_categories.product_id = products.id
          AND product_categories.category_id = sqlc.narg(category_id)));

-- name: GetCatalogProduct :one
SELECT * FROM products WHERE id = $1 AND archived_at IS NULL;
