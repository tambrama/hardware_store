-- +goose Up
INSERT INTO apps (id, name)
VALUES ('e09906ec-980c-4dd0-9251-c87105817550', 'hardware_store')
ON CONFLICT DO NOTHING;