-- +goose Up
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    full_name VARCHAR(150) NOT NULL,
    username VARCHAR(50) UNIQUE,
    email VARCHAR(150) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    role ENUM('user','admin','superadmin') DEFAULT 'user',
    status ENUM('active','inactive','banned') DEFAULT 'active',
    last_login DATETIME NULL,
    email_verified TINYINT(1) DEFAULT 0,
    phone_verified TINYINT(1) DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- +goose StatementBegin
INSERT INTO users (
    id, full_name, username, email, phone, password_hash,
    avatar_url, role, status, email_verified, phone_verified
)
VALUES
(
    UUID(),
    'Mashuri Mansur', 'mashuri', 'mashuri@example.com', '628111111111',
    '$2a$10$abcdefghijklmnopqrstuv',
    'https://example.com/avatar1.jpg',
    'admin', 'active', 1, 1
),
(
    UUID(),
    'Satria Nugraha', 'satria', 'satria@example.com', '628122222222',
    '$2a$10$abcdefghijklmnopqrstuv',
    'https://example.com/avatar2.jpg',
    'user', 'active', 0, 0
),
(
    UUID(),
    'Rizky Pratama', 'rizky', 'rizky@example.com', '628133333333',
    '$2a$10$abcdefghijklmnopqrstuv',
    'https://example.com/avatar3.jpg',
    'user', 'banned', 0, 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
