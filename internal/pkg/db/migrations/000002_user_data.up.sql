BEGIN;

CREATE TYPE user_data_type AS ENUM ('AUTH', 'TEXT', 'BINARY', 'CREDIT_CARD');

CREATE TABLE IF NOT EXISTS user_data
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER,
    data       BYTEA       NOT NULL,
    metadata   VARCHAR,
    data_type  user_data_type       DEFAULT 'AUTH',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE

);

COMMIT;
