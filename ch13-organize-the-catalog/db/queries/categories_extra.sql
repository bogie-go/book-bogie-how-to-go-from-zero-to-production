-- name: ListCategoriesForProduct :many
SELECT categories.* FROM categories
JOIN product_categories ON product_categories.category_id = categories.id
WHERE product_categories.product_id = $1
ORDER BY categories.name;

-- name: ListCategorySummaries :many
SELECT categories.id, categories.name, count(products.id) AS product_count
FROM categories
LEFT JOIN product_categories ON product_categories.category_id = categories.id
LEFT JOIN products
  ON products.id = product_categories.product_id AND products.archived_at IS NULL
GROUP BY categories.id
ORDER BY categories.name;
