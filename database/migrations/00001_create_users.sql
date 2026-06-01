-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_role AS ENUM ('USER', 'ADMIN');
CREATE TYPE subscription_type AS ENUM ('FREE', 'PREMIUM');

CREATE TABLE users
(
    id                UUID PRIMARY KEY,
    email             TEXT              NOT NULL CONSTRAINT users_email_unique UNIQUE,
    username          TEXT              NOT NULL CONSTRAINT users_username_unique UNIQUE,
    password_hash     TEXT              NOT NULL,
    salt              TEXT              NOT NULL,
    role              user_role         NOT NULL DEFAULT 'USER',
    subscription_type subscription_type NOT NULL DEFAULT 'FREE',
    created_at        TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ       NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS subscription_type;
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
