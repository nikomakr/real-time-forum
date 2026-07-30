CREATE TABLE IF NOT EXISTS comments (
    id         TEXT PRIMARY KEY,
    post_id    TEXT NOT NULL,
    author_id  TEXT NOT NULL,
    content    TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id)   REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ── Dummy comments for development and testing ────────────────────

-- Comments on post 1 (deposit dispute)
INSERT OR IGNORE INTO comments (id, post_id, author_id, content) VALUES
    (
        '019400cc-0001-7000-a000-000000000001',
        '019400bb-0001-7000-a000-000000000001',
        '019400aa-0002-7000-a000-000000000002',
        'If your deposit is protected with TDS you can raise a dispute directly through their website. The landlord has to prove any deductions are justified. You do not need a solicitor for this.'
    ),
    (
        '019400cc-0002-7000-a000-000000000002',
        '019400bb-0001-7000-a000-000000000001',
        '019400aa-0003-7000-a000-000000000003',
        'I went through TDS last year for a similar amount. The whole process took about 6 weeks and I got most of it back. Make sure you have your check-in and check-out inventories ready as evidence.'
    ),
    (
        '019400cc-0003-7000-a000-000000000003',
        '019400bb-0001-7000-a000-000000000001',
        '019400aa-0004-7000-a000-000000000004',
        'Also worth sending a formal letter before action giving the landlord 14 days to return the deposit before you escalate. Citizens Advice has a template for this.'
    ),

-- Comments on post 2 (Section 21)
    (
        '019400cc-0004-7000-a000-000000000004',
        '019400bb-0002-7000-a000-000000000002',
        '019400aa-0001-7000-a000-000000000001',
        'Section 21 was formally abolished for new tenancies under the Renters'' Rights Act 2025. However if your tenancy predates the commencement date, transitional provisions may still apply. You need to check the exact date your tenancy started.'
    ),
    (
        '019400cc-0005-7000-a000-000000000005',
        '019400bb-0002-7000-a000-000000000002',
        '019400aa-0005-7000-a000-000000000005',
        'Even if a Section 21 were valid, the landlord still needs a court order to actually evict you. Do not leave without speaking to Shelter or a housing solicitor first.'
    ),

-- Comments on post 3 (damp and mould)
    (
        '019400cc-0006-7000-a000-000000000006',
        '019400bb-0003-7000-a000-000000000003',
        '019400aa-0002-7000-a000-000000000002',
        'Under the Homes (Fitness for Human Habitation) Act 2018 your landlord has a legal duty to ensure the property is free from serious hazards including damp and mould. This is especially important given you have a young child. Document everything with photos and dates.'
    ),
    (
        '019400cc-0007-7000-a000-000000000007',
        '019400bb-0003-7000-a000-000000000003',
        '019400aa-0004-7000-a000-000000000004',
        'Report it to your local council Environmental Health team if the landlord does not act. They can issue a formal improvement notice under the Housing Health and Safety Rating System (HHSRS). This usually gets landlords moving very quickly.'
    ),

-- Comments on post 4 (Universal Credit)
    (
        '019400cc-0008-7000-a000-000000000008',
        '019400bb-0004-7000-a000-000000000004',
        '019400aa-0003-7000-a000-000000000003',
        'The LHA rate for Birmingham varies by bedroom need and the Broad Rental Market Area (BRMA). Check the Valuation Office Agency website — you can search by postcode. The rate is reviewed annually.'
    ),
    (
        '019400cc-0009-7000-a000-000000000009',
        '019400bb-0004-7000-a000-000000000004',
        '019400aa-0001-7000-a000-000000000001',
        'Be aware there is a five-week wait for the first Universal Credit payment. You can apply for an advance payment on day one to cover this gap. Make sure you do this when you first apply.'
    ),

-- Comments on post 5 (HMO licence)
    (
        '019400cc-0010-7000-a000-000000000010',
        '019400bb-0005-7000-a000-000000000005',
        '019400aa-0005-7000-a000-000000000005',
        'A mandatory HMO licence is required in England for properties with five or more occupants forming two or more households. For four occupants your landlord may be correct at the national level, but many councils have additional licensing schemes that cover smaller properties. Check your council website.'
    ),
    (
        '019400cc-0011-7000-a000-000000000011',
        '019400bb-0005-7000-a000-000000000005',
        '019400aa-0002-7000-a000-000000000002',
        'If your landlord is operating an unlicensed HMO where a licence is required, you can apply for a Rent Repayment Order through the First-tier Tribunal and claim back up to 12 months of rent. It is worth looking into.'
    );