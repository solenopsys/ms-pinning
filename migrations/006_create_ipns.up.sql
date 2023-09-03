CREATE TABLE ipns
(
    id         VARCHAR(64) PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id),
    pin_id     VARCHAR(64) REFERENCES pins (id),
    name       text,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);