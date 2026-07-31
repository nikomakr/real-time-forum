CREATE TABLE IF NOT EXISTS post_categories (
    post_id     TEXT NOT NULL,
    category_id TEXT NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id)     REFERENCES posts(id)      ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

-- ── Dummy post-category links for development and testing ─────────

-- Post 1 (deposit dispute) → Country: England + Legal Advice: Deposit Disputes
INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES
    ('019400bb-0001-7000-a000-000000000001', 'cat-001'), -- England
    ('019400bb-0001-7000-a000-000000000001', 'cat-032'); -- Deposit Disputes

-- Post 2 (Section 21) → Country: England + Legal Advice: Eviction & Possessions + Private Renters: AST
INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES
    ('019400bb-0002-7000-a000-000000000002', 'cat-001'), -- England
    ('019400bb-0002-7000-a000-000000000002', 'cat-030'), -- Eviction & Possessions
    ('019400bb-0002-7000-a000-000000000002', 'cat-062'); -- Assured Shorthold Tenancy (AST)

-- Post 3 (damp and mould) → Country: England + Housing Disrepair: Damp and Mould
INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES
    ('019400bb-0003-7000-a000-000000000003', 'cat-001'), -- England
    ('019400bb-0003-7000-a000-000000000003', 'cat-070'); -- Damp and Mould

-- Post 4 (Universal Credit) → Benefits: Universal Credit + Benefits: Local Housing Allowance
INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES
    ('019400bb-0004-7000-a000-000000000004', 'cat-016'), -- Universal Credit
    ('019400bb-0004-7000-a000-000000000004', 'cat-020'); -- Local Housing Allowance (LHA)

-- Post 5 (HMO licence) → Country: England + Shared Accommodation: HMO Licensing
INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES
    ('019400bb-0005-7000-a000-000000000005', 'cat-001'), -- England
    ('019400bb-0005-7000-a000-000000000005', 'cat-090'); -- HMO Licensing