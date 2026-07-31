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

-- ----------------------------------------------------------------
-- DUMMY DATA — development and testing only
-- password_hash below is structural only and will NOT verify
-- against any real password. To log in as a test user, register
-- a fresh account via POST /api/register after the server starts.
-- ----------------------------------------------------------------
INSERT OR IGNORE INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender) VALUES
    ('019400aa-0001-7000-a000-000000000001', 'alice_j',   'alice@test.com',   '$2a$12$KIXBsHWGgB4it6DI6NyNQOSXtbXGHMFuCTSWZTtDTNgOWtAiHXXfO', 'Alice',   'Johnson',  32, 'Female'),
    ('019400aa-0002-7000-a000-000000000002', 'bob_k',     'bob@test.com',     '$2a$12$KIXBsHWGgB4it6DI6NyNQOSXtbXGHMFuCTSWZTtDTNgOWtAiHXXfO', 'Bob',     'Khan',     45, 'Male'),
    ('019400aa-0003-7000-a000-000000000003', 'carol_m',   'carol@test.com',   '$2a$12$KIXBsHWGgB4it6DI6NyNQOSXtbXGHMFuCTSWZTtDTNgOWtAiHXXfO', 'Carol',   'Murphy',   27, 'Female'),
    ('019400aa-0004-7000-a000-000000000004', 'david_p',   'david@test.com',   '$2a$12$KIXBsHWGgB4it6DI6NyNQOSXtbXGHMFuCTSWZTtDTNgOWtAiHXXfO', 'David',   'Patel',    38, 'Male'),
    ('019400aa-0005-7000-a000-000000000005', 'eve_r',     'eve@test.com',     '$2a$12$KIXBsHWGgB4it6DI6NyNQOSXtbXGHMFuCTSWZTtDTNgOWtAiHXXfO', 'Eve',     'Roberts',  54, 'Female');