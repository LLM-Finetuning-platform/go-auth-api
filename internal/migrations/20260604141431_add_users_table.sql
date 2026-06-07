-- +goose Up
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,       -- Assumes a UUID string or custom ID
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
