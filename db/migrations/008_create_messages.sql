CREATE TABLE IF NOT EXISTS messages (
    id          TEXT PRIMARY KEY,
    sender_id   TEXT NOT NULL,
    receiver_id TEXT NOT NULL,
    content     TEXT NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id)   REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_messages_sender   ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver ON messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_created  ON messages(created_at);

-- ── Dummy messages for development and testing ────────────────────
INSERT OR IGNORE INTO messages (id, sender_id, receiver_id, content) VALUES
    (
        '019400dd-0001-7000-a000-000000000001',
        '019400aa-0001-7000-a000-000000000001',
        '019400aa-0002-7000-a000-000000000002',
        'Hi Bob, did you get a chance to look at the deposit dispute guidance?'
    ),
    (
        '019400dd-0002-7000-a000-000000000002',
        '019400aa-0002-7000-a000-000000000002',
        '019400aa-0001-7000-a000-000000000001',
        'Yes! The TDS route looks promising. Have you submitted the dispute yet?'
    ),
    (
        '019400dd-0003-7000-a000-000000000003',
        '019400aa-0001-7000-a000-000000000001',
        '019400aa-0002-7000-a000-000000000002',
        'Not yet — waiting for the check-out inventory from my old landlord.'
    ),
    (
        '019400dd-0004-7000-a000-000000000004',
        '019400aa-0003-7000-a000-000000000003',
        '019400aa-0001-7000-a000-000000000001',
        'Alice, I saw your post about the Section 21. I went through the same thing last year.'
    ),
    (
        '019400dd-0005-7000-a000-000000000005',
        '019400aa-0001-7000-a000-000000000001',
        '019400aa-0003-7000-a000-000000000003',
        'Really? What did you do? Did you get legal advice?'
    );