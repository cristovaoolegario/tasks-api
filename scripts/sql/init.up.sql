CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       role ENUM('manager', 'technician') NOT NULL,
                       UNIQUE INDEX idx_username (username)
);

CREATE TABLE tasks (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       summary VARCHAR(2500) NOT NULL,
                       performed_date DATE NOT NULL,
                       user_id INT,
                       FOREIGN KEY (user_id) REFERENCES users(id)
);
