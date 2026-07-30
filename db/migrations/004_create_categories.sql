CREATE TABLE IF NOT EXISTS categories (
    id         TEXT PRIMARY KEY,
    group_id   TEXT NOT NULL,
    name       TEXT NOT NULL,
    sub_group  TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES category_groups(id) ON DELETE RESTRICT
);

INSERT OR IGNORE INTO categories (id, group_id, name, sub_group) VALUES

    -- ── Country ───────────────────────────────────────────────────
    ('cat-001', 'grp-01', 'England',        NULL),
    ('cat-002', 'grp-01', 'Wales',          NULL),
    ('cat-003', 'grp-01', 'Scotland',       NULL),
    ('cat-004', 'grp-01', 'N. Ireland',     NULL),
    ('cat-005', 'grp-01', 'Other',          NULL),

    -- ── Benefits ──────────────────────────────────────────────────
    ('cat-010', 'grp-02', 'Child Benefit',                          NULL),
    ('cat-011', 'grp-02', 'Carer''s Allowance',                    NULL),
    ('cat-012', 'grp-02', 'Employment and Support Allowance (ESA)', NULL),
    ('cat-013', 'grp-02', 'Personal Independence Payment (PIP)',    NULL),
    ('cat-014', 'grp-02', 'Pension Credit',                         NULL),
    ('cat-015', 'grp-02', 'Jobseeker''s Allowance (JSA)',          NULL),
    ('cat-016', 'grp-02', 'Universal Credit',                       NULL),
    ('cat-017', 'grp-02', 'Discretionary Housing Payments (DHP)',   NULL),
    ('cat-018', 'grp-02', 'Support for Mortgage Interest (SMI)',    NULL),
    ('cat-019', 'grp-02', 'Housing Benefit',                        NULL),
    ('cat-020', 'grp-02', 'Local Housing Allowance (LHA)',          NULL),
    ('cat-021', 'grp-02', 'Other',                                  NULL),

    -- ── Legal Advice ──────────────────────────────────────────────
    ('cat-030', 'grp-03', 'Eviction & Possessions',                 NULL),
    ('cat-031', 'grp-03', 'Disrepair & Safety',                     NULL),
    ('cat-032', 'grp-03', 'Deposit Disputes',                       NULL),
    ('cat-033', 'grp-03', 'Harassment and Quiet Enjoyment',         NULL),
    ('cat-034', 'grp-03', 'Service Charge Disputes',                NULL),
    ('cat-035', 'grp-03', 'Lease Extensions',                       NULL),
    ('cat-036', 'grp-03', 'Enfranchisement',                        NULL),
    ('cat-037', 'grp-03', 'Staircasing (Shared Owners)',            NULL),
    ('cat-038', 'grp-03', 'Boundary and Neighbour Disputes',        NULL),
    ('cat-039', 'grp-03', 'Mortgage Arrears and Repossession',      NULL),
    ('cat-040', 'grp-03', 'Planning and Restrictive Covenants',     NULL),
    ('cat-041', 'grp-03', 'Section 21 Notice',                      NULL),
    ('cat-042', 'grp-03', 'Section 8 Notice',                       NULL),
    ('cat-043', 'grp-03', 'Rent Increase Disputes',                 NULL),
    ('cat-044', 'grp-03', 'Other',                                  NULL),

    -- ── Social and Council Housing ────────────────────────────────
    ('cat-050', 'grp-04', 'Social Rent',        'Rent Type'),
    ('cat-051', 'grp-04', 'Affordable Rent',    'Rent Type'),
    ('cat-052', 'grp-04', 'Other',              'Rent Type'),
    ('cat-053', 'grp-04', 'Secured Tenancy',    'Tenancy Type'),
    ('cat-054', 'grp-04', 'Assured Tenancy',    'Tenancy Type'),
    ('cat-055', 'grp-04', 'Introductory Tenancy', 'Tenancy Type'),
    ('cat-056', 'grp-04', 'Flexible Tenancy',   'Tenancy Type'),
    ('cat-057', 'grp-04', 'Other',              'Tenancy Type'),

    -- ── Private Renters ───────────────────────────────────────────
    ('cat-060', 'grp-05', 'Market Rent',                         'Rent Type'),
    ('cat-061', 'grp-05', 'Other',                               'Rent Type'),
    ('cat-062', 'grp-05', 'Assured Shorthold Tenancy (AST)',     'Tenancy Type'),
    ('cat-063', 'grp-05', 'Assured Periodic Tenancy (APT)',      'Tenancy Type'),
    ('cat-064', 'grp-05', 'Excluded Tenancy (Lodging)',          'Tenancy Type'),
    ('cat-065', 'grp-05', 'Non-AST (Contractual Tenancy)',       'Tenancy Type'),
    ('cat-066', 'grp-05', 'Company Let',                         'Tenancy Type'),
    ('cat-067', 'grp-05', 'Other',                               'Tenancy Type'),

    -- ── Housing Disrepair & Safety ────────────────────────────────
    ('cat-070', 'grp-06', 'Damp and Mould',                     NULL),
    ('cat-071', 'grp-06', 'Structural Defects',                  NULL),
    ('cat-072', 'grp-06', 'Heating and Hot Water Failure',       NULL),
    ('cat-073', 'grp-06', 'Electrical Safety',                   NULL),
    ('cat-074', 'grp-06', 'Gas Safety',                          NULL),
    ('cat-075', 'grp-06', 'Fire Safety',                         NULL),
    ('cat-076', 'grp-06', 'Pest Infestation',                    NULL),
    ('cat-077', 'grp-06', 'Asbestos',                            NULL),
    ('cat-078', 'grp-06', 'Carbon Monoxide Risk',                NULL),
    ('cat-079', 'grp-06', 'Other',                               NULL),

    -- ── Letting Agents & Fees ─────────────────────────────────────
    ('cat-080', 'grp-07', 'Unfair Fees',                         NULL),
    ('cat-081', 'grp-07', 'Agent Conduct',                       NULL),
    ('cat-082', 'grp-07', 'Referencing Issues',                  NULL),
    ('cat-083', 'grp-07', 'Inventory Disputes',                  NULL),
    ('cat-084', 'grp-07', 'Contract Transparency',               NULL),
    ('cat-085', 'grp-07', 'Other',                               NULL),

    -- ── Shared Accommodation (HMO) ────────────────────────────────
    ('cat-090', 'grp-08', 'HMO Licensing',                       NULL),
    ('cat-091', 'grp-08', 'Shared Bills',                        NULL),
    ('cat-092', 'grp-08', 'Housemate Disputes',                  NULL),
    ('cat-093', 'grp-08', 'Common Areas',                        NULL),
    ('cat-094', 'grp-08', 'Overcrowding',                        NULL),
    ('cat-095', 'grp-08', 'Other',                               NULL),

    -- ── Student Housing ───────────────────────────────────────────
    ('cat-100', 'grp-09', 'University Accommodation',            NULL),
    ('cat-101', 'grp-09', 'Private Student Lets',                NULL),
    ('cat-102', 'grp-09', 'Purpose-Built Student Accommodation', NULL),
    ('cat-103', 'grp-09', 'Guarantor Issues',                    NULL),
    ('cat-104', 'grp-09', 'Joint Tenancy Disputes',              NULL),
    ('cat-105', 'grp-09', 'Other',                               NULL),

    -- ── Short-Term & Holiday Lets ─────────────────────────────────
    ('cat-110', 'grp-10', 'Airbnb & Platform Lets',              NULL),
    ('cat-111', 'grp-10', 'Serviced Apartments',                 NULL),
    ('cat-112', 'grp-10', 'Corporate Lets',                      NULL),
    ('cat-113', 'grp-10', 'Licensing Requirements',              NULL),
    ('cat-114', 'grp-10', 'Neighbour Complaints',                NULL),
    ('cat-115', 'grp-10', 'Other',                               NULL),

    -- ── Utilities & Bills ─────────────────────────────────────────
    ('cat-120', 'grp-11', 'Electricity',                         NULL),
    ('cat-121', 'grp-11', 'Gas',                                 NULL),
    ('cat-122', 'grp-11', 'Water',                               NULL),
    ('cat-123', 'grp-11', 'Internet & Broadband',                NULL),
    ('cat-124', 'grp-11', 'Council Tax',                         NULL),
    ('cat-125', 'grp-11', 'Prepayment Meters',                   NULL),
    ('cat-126', 'grp-11', 'Energy Efficiency & EPC',             NULL),
    ('cat-127', 'grp-11', 'Other',                               NULL),

    -- ── Disability & Accessible Housing ──────────────────────────
    ('cat-130', 'grp-12', 'Adaptations and Adjustments',        NULL),
    ('cat-131', 'grp-12', 'Disabled Facilities Grant',           NULL),
    ('cat-132', 'grp-12', 'Wheelchair Accessible Housing',       NULL),
    ('cat-133', 'grp-12', 'Care and Support Housing',            NULL),
    ('cat-134', 'grp-12', 'Other',                               NULL),

    -- ── Immigration & Housing Rights ──────────────────────────────
    ('cat-140', 'grp-13', 'Right to Rent Checks',               NULL),
    ('cat-141', 'grp-13', 'No Recourse to Public Funds (NRPF)', NULL),
    ('cat-142', 'grp-13', 'Asylum Seeker Housing',              NULL),
    ('cat-143', 'grp-13', 'Refugee Housing',                    NULL),
    ('cat-144', 'grp-13', 'Visa and Tenancy',                   NULL),
    ('cat-145', 'grp-13', 'Other',                              NULL),

    -- ── Right to Buy / Right to Acquire ──────────────────────────
    ('cat-150', 'grp-14', 'Right to Buy Eligibility',           NULL),
    ('cat-151', 'grp-14', 'Right to Acquire Eligibility',       NULL),
    ('cat-152', 'grp-14', 'Valuation Disputes',                 NULL),
    ('cat-153', 'grp-14', 'Delays and Refusals',                NULL),
    ('cat-154', 'grp-14', 'Preserved Right to Buy',             NULL),
    ('cat-155', 'grp-14', 'Other',                              NULL),

    -- ── Discrimination in Housing ─────────────────────────────────
    ('cat-160', 'grp-15', 'Race Discrimination',                NULL),
    ('cat-161', 'grp-15', 'Disability Discrimination',          NULL),
    ('cat-162', 'grp-15', 'Gender Discrimination',              NULL),
    ('cat-163', 'grp-15', 'Age Discrimination',                 NULL),
    ('cat-164', 'grp-15', 'DSS / Benefits Discrimination',      NULL),
    ('cat-165', 'grp-15', 'Families with Children',             NULL),
    ('cat-166', 'grp-15', 'Sexual Orientation',                 NULL),
    ('cat-167', 'grp-15', 'Religion or Belief',                 NULL),
    ('cat-168', 'grp-15', 'Other',                              NULL),

    -- ── Other ─────────────────────────────────────────────────────
    ('cat-999', 'grp-16', 'Other',                              NULL);