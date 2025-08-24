-- +goose Up
CREATE TABLE users (
                       user_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       login VARCHAR(255) NOT NULL UNIQUE,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password VARCHAR(255) NOT NULL,
                       notification_methods JSONB DEFAULT '[]'::jsonb,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP
);

ALTER TABLE users ADD CONSTRAINT unique_login UNIQUE (login);
ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);

CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_notification_methods ON users USING gin (notification_methods);

-- +goose Down
DROP TABLE IF EXISTS users;
