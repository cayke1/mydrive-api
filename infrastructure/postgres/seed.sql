-- MyDrive Quick Seed Data
-- Insert sample users, folders, and files for development/testing

BEGIN;

-- Users
INSERT INTO users (id, email, password_hash, session_token, CSRF_token, created_at, updated_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'john@example.com', '$2a$12$abc123def456ghi789jkl', 'session_john_001', 'csrf_john_001', NOW(), NOW()),
    ('550e8400-e29b-41d4-a716-446655440002', 'jane@example.com', '$2a$12$xyz789uvw012pqr345stu', 'session_jane_001', 'csrf_jane_001', NOW(), NOW()),
    ('550e8400-e29b-41d4-a716-446655440003', 'alice@example.com', '$2a$12$vwx012yza345bcd678efg', 'session_alice_001', 'csrf_alice_001', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Root folders for each user
INSERT INTO folders (id, name, created_at, updated_at, parent_id, owner_id) VALUES
    -- John's folders
    ('650e8400-e29b-41d4-a716-446655440001', 'My Drive', NOW(), NOW(), NULL, '550e8400-e29b-41d4-a716-446655440001'),
    ('650e8400-e29b-41d4-a716-446655440002', 'Documents', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001'),
    ('650e8400-e29b-41d4-a716-446655440003', 'Photos', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001'),
    ('650e8400-e29b-41d4-a716-446655440004', 'Projects', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001'),

    -- Jane's folders
    ('650e8400-e29b-41d4-a716-446655440005', 'My Drive', NOW(), NOW(), NULL, '550e8400-e29b-41d4-a716-446655440002'),
    ('650e8400-e29b-41d4-a716-446655440006', 'Work', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440002'),

    -- Alice's folders
    ('650e8400-e29b-41d4-a716-446655440007', 'My Drive', NOW(), NOW(), NULL, '550e8400-e29b-41d4-a716-446655440003')
ON CONFLICT DO NOTHING;

-- Sample files
INSERT INTO files (id, name, size, mime_type, storage_key, checksum, created_at, updated_at, folder_id, owner_id) VALUES
    -- John's files
    ('750e8400-e29b-41d4-a716-446655440001', 'Resume.pdf', 125456, 'application/pdf', 'john/resume.pdf', 'abc123def456', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001'),
    ('750e8400-e29b-41d4-a716-446655440002', 'Cover Letter.docx', 45678, 'application/vnd.openxmlformats-officedocument.wordprocessingml.document', 'john/cover_letter.docx', 'xyz789uvw012', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001'),
    ('750e8400-e29b-41d4-a716-446655440003', 'Vacation.jpg', 2097152, 'image/jpeg', 'john/vacation.jpg', 'img123456789', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440001'),
    ('750e8400-e29b-41d4-a716-446655440004', 'Architecture.png', 854321, 'image/png', 'john/architecture.png', 'arch123456789', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440001'),
    ('750e8400-e29b-41d4-a716-446655440005', 'Project Proposal.xlsx', 98765, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'john/project_proposal.xlsx', 'xl123456789', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440001'),

    -- Jane's files
    ('750e8400-e29b-41d4-a716-446655440006', 'Quarterly Report.pdf', 234567, 'application/pdf', 'jane/quarterly_report.pdf', 'rep123456789', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440002'),
    ('750e8400-e29b-41d4-a716-446655440007', 'Meeting Notes.txt', 5432, 'text/plain', 'jane/meeting_notes.txt', 'notes123456789', NOW(), NOW(), '650e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440002')
ON CONFLICT DO NOTHING;

COMMIT;
