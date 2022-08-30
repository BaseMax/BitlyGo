-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
---- create above / drop below ----
DROP TABLE users;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
