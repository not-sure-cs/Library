-- +goose Up
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT,
    ph_no TEXT,
    role TEXT,
    PRIMARY KEY(id)
);

CREATE TABLE secrets (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    pass_hash VARCHAR(255) NOT NULL,
    PRIMARY KEY (user_id)
);

-- +goose Down
DROP TABLE secrets;

DROP TABLE users;

