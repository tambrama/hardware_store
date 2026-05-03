-- +goose Up
INSERT INTO apps (id, name)
VALUES ('e09906ec-980c-4dd0-9251-c87105817550', 'hardware_store'),
       ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'store_python')
ON CONFLICT DO NOTHING;