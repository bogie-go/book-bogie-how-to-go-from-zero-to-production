-- name: ListCartItemsForCart :many
SELECT cart_items.id, cart_items.product_id, cart_items.quantity,
       products.name, products.image_url, products.price_cents
FROM cart_items
JOIN products ON products.id = cart_items.product_id
WHERE cart_items.cart_id = $1
ORDER BY cart_items.created_at, cart_items.id;

-- name: AddToCart :one
INSERT INTO cart_items (cart_id, product_id, quantity)
VALUES (sqlc.arg(cart_id), sqlc.arg(product_id), sqlc.arg(quantity))
ON CONFLICT (cart_id, product_id) DO UPDATE
SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = now()
WHERE cart_items.quantity + EXCLUDED.quantity <= sqlc.arg(stock)::integer
RETURNING *;

-- name: GetCartItemInCart :one
SELECT * FROM cart_items WHERE id = $1 AND cart_id = $2;

-- name: DeleteCartItemInCart :execrows
DELETE FROM cart_items WHERE id = $1 AND cart_id = $2;
