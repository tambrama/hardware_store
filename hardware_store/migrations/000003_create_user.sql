-- +goose Up
-- +goose StatementBegin
CREATE USER read_user WITH PASSWORD 'read_password';

GRANT CONNECT ON DATABASE hardwarestore TO read_user;

GRANT USAGE ON SCHEMA public TO read_user;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO read_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO read_user;
-- +goose StatementEnd