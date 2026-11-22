-- +goose Up
ALTER TABLE users ADD COLUMN role varchar NOT NULL DEFAULT '' CHECK (role IN ('', 'staff'));

-- +goose Down
ALTER TABLE users DROP COLUMN role;
