CREATE TABLE pins (
                       id  PRIMARY KEY,
                       user_id foreign key REFERENCES users(id),
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       size bigint,
                       state VARCHAR(255) NOT NULL DEFAULT 'new'
);