-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE urls
ADD COLUMN user_id BIGINT REFERENCES users (id) ON DELETE SET NULL;

CREATE INDEX idx_urls_user_id ON urls (user_id);

-- +goose Down
ALTER TABLE urls
DROP COLUMN user_id;

DROP TABLE users;
