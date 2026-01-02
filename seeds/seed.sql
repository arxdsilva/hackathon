-- Seed data for hackathon after migrations

-- Insert users - password hashed for 'password'
INSERT INTO users (id, created_at, updated_at, name, email, company_team, role, password_hash, force_password_reset) VALUES
('550e8400-e29b-41d4-a716-446655440000', NOW(), NOW(), 'Admin Owner', 'admin@example.com', 'Platform', 'owner', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6', false),
('550e8400-e29b-41d4-a716-446655440001', NOW(), NOW(), 'John Developer', 'john@example.com', 'Developer Experience', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6', false),
('550e8400-e29b-41d4-a716-446655440002', NOW(), NOW(), 'Jane Designer', 'jane@example.com', 'Design', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6', false),
('550e8400-e29b-41d4-a716-446655440003', NOW(), NOW(), 'Bob Engineer', 'bob@example.com', 'Infra', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6', false),
('550e8400-e29b-41d4-a716-446655440004', NOW(), NOW(), 'Alice Analyst', 'alice@example.com', '', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6', false);

-- Insert hackathons
INSERT INTO hackathons (id, created_at, updated_at, title, description, start_date, end_date, owner_id, status) VALUES
('HACK25WINTER', NOW(), NOW(), 'Winter Code Fest 2025', 'A month-long coding competition focusing on innovative solutions for real-world problems. Teams will build applications using modern technologies.', '2025-12-27 09:00:00', '2026-01-27 17:00:00', '550e8400-e29b-41d4-a716-446655440000', 'active'),
('HACK26AIINNO', NOW(), NOW(), 'AI Innovation Challenge', 'Explore the possibilities of artificial intelligence and machine learning. Create projects that demonstrate practical applications of AI technologies.', '2026-02-01 10:00:00', '2026-02-28 16:00:00', '550e8400-e29b-41d4-a716-446655440000', 'upcoming'),
('HACK26GREENX', NOW(), NOW(), 'Green Tech Hackathon', 'Build sustainable technology solutions that address environmental challenges. Focus on renewable energy, waste reduction, and eco-friendly innovations.', '2026-03-15 08:00:00', '2026-03-29 18:00:00', '550e8400-e29b-41d4-a716-446655440000', 'upcoming');

-- Insert projects
INSERT INTO projects (id, user_id, created_at, updated_at, hackathon_id, name, description, repository_url, demo_url, status, team_members) VALUES
('PRJECO123456', '550e8400-e29b-41d4-a716-446655440001', NOW(), NOW(), 'HACK25WINTER', 'EcoTracker', 'Mobile app that helps users track carbon footprint and suggests reductions through gamification.', 'https://github.com/johndev/ecotracker', 'https://ecotracker-demo.herokuapp.com', 'in_progress', 2),
('PRJWASTE78901', '550e8400-e29b-41d4-a716-446655440002', NOW(), NOW(), 'HACK25WINTER', 'Smart Waste Sorter', 'IoT device that automatically sorts household waste using computer vision.', 'https://github.com/janedesign/smart-waste', 'https://smartwaste-demo.netlify.app', 'completed', 2),
('PRJGARDEN2345', '550e8400-e29b-41d4-a716-446655440003', NOW(), NOW(), 'HACK25WINTER', 'Community Garden Planner', 'Web platform to plan/manage shared gardens and coordinate volunteers.', 'https://github.com/bobeng/community-garden', 'https://gardenshare.vercel.app', 'in_progress', 1),
('PRJAIRECIPE67', '550e8400-e29b-41d4-a716-446655440004', NOW(), NOW(), 'HACK26AIINNO', 'AI Recipe Generator', 'AI-powered personalized recipes based on preferences and ingredients.', 'https://github.com/aliceai/ai-recipes', 'https://airecipes.streamlit.app', 'in_progress', 2),
('PRJSENTIMENT8', '550e8400-e29b-41d4-a716-446655440001', NOW(), NOW(), 'HACK26AIINNO', 'Sentiment Analysis Dashboard', 'Real-time dashboard analyzing social media sentiment using NLP.', 'https://github.com/johndev/sentiment-dash', 'https://sentiment-analytics.fly.dev', 'completed', 1),
('PRJSOLAR90123', '550e8400-e29b-41d4-a716-446655440002', NOW(), NOW(), 'HACK26GREENX', 'Solar Panel Optimizer', 'Smart system that maximizes solar panel efficiency via weather-driven adjustments.', 'https://github.com/janedesign/solar-opt', 'https://solaroptimizer.glitch.me', 'in_progress', 1);

-- Insert project memberships
INSERT INTO project_memberships (id, created_at, updated_at, project_id, user_id) VALUES
('650e8400-e29b-41d4-a716-446655440001', NOW(), NOW(), 'PRJECO123456', '550e8400-e29b-41d4-a716-446655440001'), -- John is founder
('650e8400-e29b-41d4-a716-446655440002', NOW(), NOW(), 'PRJECO123456', '550e8400-e29b-41d4-a716-446655440003'), -- Bob joins John's project
('650e8400-e29b-41d4-a716-446655440003', NOW(), NOW(), 'PRJWASTE78901', '550e8400-e29b-41d4-a716-446655440002'), -- Jane is founder
('650e8400-e29b-41d4-a716-446655440004', NOW(), NOW(), 'PRJWASTE78901', '550e8400-e29b-41d4-a716-446655440004'), -- Alice joins Jane's project
('650e8400-e29b-41d4-a716-446655440005', NOW(), NOW(), 'PRJGARDEN2345', '550e8400-e29b-41d4-a716-446655440003'), -- Bob is founder
('650e8400-e29b-41d4-a716-446655440006', NOW(), NOW(), 'PRJAIRECIPE67', '550e8400-e29b-41d4-a716-446655440004'), -- Alice is founder
('650e8400-e29b-41d4-a716-446655440007', NOW(), NOW(), 'PRJAIRECIPE67', '550e8400-e29b-41d4-a716-446655440002'), -- Jane joins Alice's project
('650e8400-e29b-41d4-a716-446655440008', NOW(), NOW(), 'PRJSENTIMENT8', '550e8400-e29b-41d4-a716-446655440001'), -- John is founder
('650e8400-e29b-41d4-a716-446655440009', NOW(), NOW(), 'PRJSOLAR90123', '550e8400-e29b-41d4-a716-446655440002'); -- Jane is founder

-- Insert some sample files (using dummy data for demonstration)
INSERT INTO files (id, created_at, updated_at, filename, data, content_type, size, user_id, hackathon_id, project_id) VALUES
('FILEECO12345', NOW(), NOW(), 'ecotracker-mockup.png', decode('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==', 'base64'), 'image/png', 68, '550e8400-e29b-41d4-a716-446655440001', 'HACK25WINTER', 'PRJECO123456'),
('FILEWASTE678', NOW(), NOW(), 'waste_sorter_specs.pdf', decode('JVBERi0xLjQKMSAwIG9iago8PC9UeXBlIC9DYXRhbG9nCi9QYWdlcyAyIDAgUgo+PgplbmRvYmoKMiAwIG9iago8PC9UeXBlIC9QYWdlcwovS2lkcyBbMyAwIFJdCi9Db3VudCAxCj4+CmVuZG9iagozIDAgb2JqCjw8L1R5cGUgL1BhZ2UKL1BhcmVudCAyIDAgUgovTWVkaWFCb3ggWzAgMCA2MTIgNzkyXQovQ29udGVudHMgNCAwIFIKPj4KZW5kb2JqCjQgMCBvYmoKPDwvTGVuZ3RoIDQ0Pj4Kc3RyZWFtCkJUCjEwIDUwIFRECi9GMSAxMiBUZgooSGVsbG8gV29ybGQpIFRqCkVUCmVuZHN0cmVhbQplbmRvYmoKeHJlZgowIDUKMDAwMDAwMDAwMCA2NTUzNSBmIAowMDAwMDAwMDA5IDAwMDAwIG4gCjAwMDAwMDAwNTkgMDAwMDAgbiAKMDAwMDAwMDEwMyAwMDAwMCBuIAp0cmFpbGVyCjw8L1NpemUgNQovUm9vdCAxIDAgUgo+PgpzdGFydHhyZWYKNzE5CmUlRU9GCg==', 'base64'), 'application/pdf', 719, '550e8400-e29b-41d4-a716-446655440002', 'HACK25WINTER', 'PRJWASTE78901');