-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, github_id, username, role) VALUES (1, 0, 'logr', 1) ON DUPLICATE KEY UPDATE username = username;
INSERT INTO users (id, github_id, username, role) VALUES (2, 55717547, 'kidlog', 4) ON DUPLICATE KEY UPDATE username = username;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE users;
-- +goose StatementEnd
