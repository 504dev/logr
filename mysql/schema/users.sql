CREATE TABLE IF NOT EXISTS users (
    id INT NOT NULL AUTO_INCREMENT,
    github_id INT,
    username VARCHAR(32),
    PRIMARY KEY (id)
);