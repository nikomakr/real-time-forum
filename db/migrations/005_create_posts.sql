CREATE TABLE IF NOT EXISTS posts (
    id         TEXT PRIMARY KEY,
    author_id  TEXT NOT NULL,
    title      TEXT NOT NULL,
    content    TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ── Dummy posts for development and testing ───────────────────────
INSERT OR IGNORE INTO posts (id, author_id, title, content) VALUES
    (
        '019400bb-0001-7000-a000-000000000001',
        '019400aa-0001-7000-a000-000000000001',
        'Landlord refusing to return my deposit after 6 weeks — what can I do?',
        'My tenancy ended on 1 June 2026. The landlord has not returned my £1,200 deposit and is not responding to emails. I have a receipt from the Tenancy Deposit Scheme (TDS). What are my options for getting this back?'
    ),
    (
        '019400bb-0002-7000-a000-000000000002',
        '019400aa-0002-7000-a000-000000000002',
        'Received a Section 21 notice — is it valid under the new rules?',
        'My landlord has issued a Section 21 notice dated July 2026. I thought Section 21 was being abolished under the Renters'' Rights Act. Does this notice have any legal force? I am in England on an Assured Shorthold Tenancy.'
    ),
    (
        '019400bb-0003-7000-a000-000000000003',
        '019400aa-0003-7000-a000-000000000003',
        'Severe damp and mould in my flat — landlord says it is my fault',
        'I have black mould covering an entire wall in my bedroom. My landlord is blaming condensation from my lifestyle. I have a two-year-old and I am very concerned about health. Who is responsible and what can I do?'
    ),
    (
        '019400bb-0004-7000-a000-000000000004',
        '019400aa-0004-7000-a000-000000000004',
        'Universal Credit housing element — how is it calculated for my area?',
        'I have just been made redundant and need to claim Universal Credit. I am a private renter in Birmingham paying £950 per month. Will the Local Housing Allowance cover this? How do I find out my LHA rate?'
    ),
    (
        '019400bb-0005-7000-a000-000000000005',
        '019400aa-0005-7000-a000-000000000005',
        'HMO licence — is my landlord required to have one for a 4-bed shared house?',
        'I live in a shared house with three other tenants. We have four separate bedrooms and share a kitchen and bathroom. Our landlord told us he does not need an HMO licence. I am not sure this is correct for our area.'
    );