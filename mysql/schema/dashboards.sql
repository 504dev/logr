CREATE TABLE IF NOT EXISTS dashboards (
    id INT NOT NULL AUTO_INCREMENT,
    owner_id INT,
    name VARCHAR(32),
    PRIMARY KEY (id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);