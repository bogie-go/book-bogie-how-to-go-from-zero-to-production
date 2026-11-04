-- name: ListActivePriceRules :many
SELECT * FROM price_rules
WHERE starts_at <= now() AND ends_at > now()
ORDER BY priority, id;

-- name: GetPriceRuleByName :one
SELECT * FROM price_rules WHERE name = $1;
