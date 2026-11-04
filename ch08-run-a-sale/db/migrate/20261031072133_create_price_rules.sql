-- +goose Up
CREATE TABLE price_rules (
    id         bigserial   PRIMARY KEY,
    name       varchar     NOT NULL DEFAULT '' UNIQUE,
    kind       varchar     NOT NULL DEFAULT '',
    properties jsonb       NOT NULL DEFAULT '{}',
    priority   integer     NOT NULL DEFAULT 0,
    starts_at  timestamptz NOT NULL,
    ends_at    timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE price_rules;
