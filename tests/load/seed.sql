INSERT INTO users (id, email, password_hash, role, created_at, updated_at) VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'organizer@loadtest.com', '$2a$12$YwRyt8wAQrp2TX9rNDLSye9TZ2pgILU7cMVbi4ecYAqA6EsIPRQge', 'organizer', NOW(), NOW()) ON CONFLICT (id) DO NOTHING;

INSERT INTO events (id, name, capacity, available_spots,price, start_at, end_at, created_at, updated_at) VALUES 
('a0000000-0000-0000-0000-000000000001', 'Load Test Event', 10000, 10000, 100, '2026-03-01 12:00:00+00', '2026-03-03 12:00:00+00', NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
