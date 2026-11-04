-- name: ListOrderItemsForOrder :many
SELECT * FROM order_items WHERE order_id = $1 ORDER BY id;
