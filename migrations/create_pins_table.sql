CREATE TABLE users (
                       id  PRIMARY KEY,
                       user_id foreign key REFERENCES users(id),
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       pin_size bigint,
);