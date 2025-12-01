-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    bio TEXT,
    birth_date DATE,
    country TEXT,
    city TEXT,
    street TEXT,
    account_type TEXT,
    account_start TIMESTAMP,
    renewal_period INT,
    role TEXT,
    security_label TEXT,
    attributes JSONB DEFAULT '{}'::jsonb
);


CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS users_username_unique ON users(username);

-- +goose Down
DROP TABLE IF EXISTS users;
