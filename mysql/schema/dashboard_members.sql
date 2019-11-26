CREATE TABLE IF NOT EXISTS dashboard_members (
    id INT NOT NULL AUTO_INCREMENT,
    dash_id INT,
    user_id INT,
    PRIMARY KEY (id),
    FOREIGN KEY (dash_id) REFERENCES dashboards(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);