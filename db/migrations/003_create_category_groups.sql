CREATE TABLE IF NOT EXISTS category_groups (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO category_groups (id, name) VALUES
    ('grp-01', 'Country'),
    ('grp-02', 'Benefits'),
    ('grp-03', 'Legal Advice'),
    ('grp-04', 'Social and Council Housing'),
    ('grp-05', 'Private Renters'),
    ('grp-06', 'Housing Disrepair & Safety'),
    ('grp-07', 'Letting Agents & Fees'),
    ('grp-08', 'Shared Accommodation (HMO)'),
    ('grp-09', 'Student Housing'),
    ('grp-10', 'Short-Term & Holiday Lets'),
    ('grp-11', 'Utilities & Bills'),
    ('grp-12', 'Disability & Accessible Housing'),
    ('grp-13', 'Immigration & Housing Rights'),
    ('grp-14', 'Right to Buy / Right to Acquire'),
    ('grp-15', 'Discrimination in Housing'),
    ('grp-16', 'Other');