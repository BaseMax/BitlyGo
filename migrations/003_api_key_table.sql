-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS api_keys(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    key VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);
---- create above / drop below ----
DROP TABLE api_keys;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
