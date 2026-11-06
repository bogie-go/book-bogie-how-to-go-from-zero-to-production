-- +goose Up
CREATE TABLE carts (
    id         bigserial   PRIMARY KEY,
    token      varchar     NOT NULL DEFAULT '' UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE carts;
