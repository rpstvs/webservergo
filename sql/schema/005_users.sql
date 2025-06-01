-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red boolean not NULL DEFAULT FALSE;
-- +goose Down
ALTER TABLE users DROP COLUMN is_chirpy_red;