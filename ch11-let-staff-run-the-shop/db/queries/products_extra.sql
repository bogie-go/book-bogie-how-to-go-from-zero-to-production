-- name: GetProductByName :one
SELECT * FROM products WHERE name = $1;

-- name: TakeStock :execrows
UPDATE products
SET stock = stock - sqlc.arg(quantity), updated_at = now()
WHERE id = sqlc.arg(id) AND stock >= sqlc.arg(quantity);
