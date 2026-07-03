CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    nickname      TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name    TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    age           INTEGER NOT NULL,
    gender        TEXT NOT NULL,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);