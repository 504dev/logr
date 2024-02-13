-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dashboard_members (
    id INT NOT NULL AUTO_INCREMENT,
    dash_id INT,
    user_id INT,
    status INT,
    PRIMARY KEY (id),
    UNIQUE (dash_id, user_id),
    FOREIGN KEY (dash_id) REFERENCES dashboards(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dashboard_members;
-- +goose StatementEnd
