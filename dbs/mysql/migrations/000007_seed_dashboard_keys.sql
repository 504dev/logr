-- +goose Up
-- +goose StatementBegin
INSERT INTO dashboard_keys (id, dash_id, name, public_key, private_key) VALUES (1, 1, 'Default', MD5(RAND()), SHA2(RAND(), 256)) ON DUPLICATE KEY UPDATE name = name;
INSERT INTO dashboard_keys (id, dash_id, name, public_key, private_key) VALUES (2, 2, 'Default', MD5(RAND()), SHA2(RAND(), 256)) ON DUPLICATE KEY UPDATE name = name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE dashboard_keys;
-- +goose StatementEnd
