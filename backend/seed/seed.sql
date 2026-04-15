-- Password is 'password123' (bcrypt cost 12)
INSERT INTO users (id, name, email, password) VALUES 
    ('00000000-0000-0000-0000-000000000001',
     'Test User', 
     'test@example.com',
     '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5I5wG5Q1Fp2uO');

INSERT INTO projects (id, name, description, owner_id) VALUES
    ('00000000-0000-0000-0000-000000000010',
     'Demo Project',
     'A project for testing',
     '00000000-0000-0000-0000-000000000001');

INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id) VALUES
    ('00000000-0000-0000-0000-000000000100', 'Set up repository', 'Initialize the repo', 'done', 'high',
     '00000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000001'),
    ('00000000-0000-0000-0000-000000000101', 'Write API endpoints', 'Build REST API', 'in_progress', 'high',
     '00000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000001'),
    ('00000000-0000-0000-0000-000000000102', 'Write tests', 'Integration tests', 'todo', 'medium',
     '00000000-0000-0000-0000-000000000010', NULL)
ON CONFLICT DO NOTHING;