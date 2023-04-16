CREATE TABLE pins
(
    id         VARCHAR(64) PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    size       bigint,
    state      VARCHAR(255) NOT NULL DEFAULT 'new'
);