-- +goose Up
CREATE UNIQUE INDEX cart_items_cart_id_product_id_idx ON cart_items (cart_id, product_id);

-- +goose Down
DROP INDEX cart_items_cart_id_product_id_idx;
