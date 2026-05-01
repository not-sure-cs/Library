-- +goose Up

CREATE EXTENSION pgcrypto;

CREATE TABLE session (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    encrypted_session bytea,
    logged_at TIMESTAMP NOT NULL
);

-- +goose Down 

DROP TABLE session;

DROP EXTENSION pgcrypto;