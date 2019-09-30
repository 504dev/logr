CREATE TABLE IF NOT EXISTS dashboards (
    id INT NOT NULL AUTO_INCREMENT,
    owner_id INT,
    name VARCHAR(32),
    public_key VARCHAR(1024),
    private_key VARCHAR(1024),
    PRIMARY KEY (id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);