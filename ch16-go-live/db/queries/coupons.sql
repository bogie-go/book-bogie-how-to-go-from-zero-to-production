-- Written by `bogie g model`: overwritten whole on the next run. A
-- hand-written query goes in coupons_extra.sql, not here.

-- name: GetCoupon :one
SELECT * FROM coupons WHERE id = $1;

-- name: ListCoupons :many
SELECT * FROM coupons ORDER BY created_at DESC, id DESC LIMIT $1;

-- name: CreateCoupon :one
INSERT INTO coupons (code, price_rule_id) VALUES ($1, $2) RETURNING *;

-- name: UpdateCoupon :one
UPDATE coupons SET code = $2, price_rule_id = $3, updated_at = now() WHERE id = $1 RETURNING *;

-- name: DeleteCoupon :execrows
DELETE FROM coupons WHERE id = $1;
