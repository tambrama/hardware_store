-- +goose Up
-- +goose StatementBegin
ALTER TABLE product DROP CONSTRAINT IF EXISTS product_image_id_fkey;
-- +goose StatementEnd