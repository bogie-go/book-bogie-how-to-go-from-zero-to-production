-- +goose Up
CREATE TABLE coupons (
    id         bigserial   PRIMARY KEY,
    code       varchar     NOT NULL DEFAULT '' UNIQUE,
    price_rule_id bigint      NOT NULL REFERENCES price_rules (id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX coupons_price_rule_id_idx ON coupons (price_rule_id);

-- +goose Down
DROP TABLE coupons;
