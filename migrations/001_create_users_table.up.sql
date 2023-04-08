CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    public_key VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);