-- name: ListActivePriceRules :many
SELECT * FROM price_rules
WHERE starts_at <= now() AND ends_at > now()
  AND NOT EXISTS (SELECT 1 FROM coupons WHERE coupons.price_rule_id = price_rules.id)
ORDER BY priority, id;

-- name: GetActivePriceRuleByCouponCode :one
SELECT sqlc.embed(price_rules) FROM price_rules
JOIN coupons ON coupons.price_rule_id = price_rules.id
WHERE coupons.code = $1 AND price_rules.starts_at <= now() AND price_rules.ends_at > now();

-- name: GetPriceRuleByName :one
SELECT * FROM price_rules WHERE name = $1;
