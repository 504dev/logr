CREATE TABLE IF NOT EXISTS users (
    id INT NOT NULL AUTO_INCREMENT,
    github_id INT,
    username VARCHAR(32),
    role INT,
    login_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (github_id),
    PRIMARY KEY (id)
);