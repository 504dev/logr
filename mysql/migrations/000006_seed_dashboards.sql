-- +goose Up
-- +goose StatementBegin
INSERT INTO dashboards (id, owner_id, name) VALUES (1, 1, 'System') ON DUPLICATE KEY UPDATE name = name;
INSERT INTO dashboards (id, owner_id, name) VALUES (2, 1, 'Demo') ON DUPLICATE KEY UPDATE name = name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE dashboards;
-- +goose StatementEnd
