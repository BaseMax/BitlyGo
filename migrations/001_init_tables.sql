-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

create table links(
    id SERIAL PRIMARY KEY,
    owner_id INTEGER,
    name VARCHAR(155) NOT NULL UNIQUE,
    link TEXT NOT NULL,
    visits INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT null default NOW(),
    updated_at TIMESTAMP,
    expired_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_owner_id FOREIGN KEY(owner_id) REFERENCES users(id)
);

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
DROP TABLE links;
DROP TABLE users;
DROP TABLE api_keys;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
