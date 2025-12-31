-- Seed data for hackathon after migrations

-- Insert users - pass hashed password for 'password'
INSERT INTO users (id, created_at, updated_at, name, email, company_team, role, password_hash) VALUES
('550e8400-e29b-41d4-a716-446655440000', NOW(), NOW(), 'Admin Owner', 'admin@example.com', 'Platform', 'owner', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6'),
('550e8400-e29b-41d4-a716-446655440001', NOW(), NOW(), 'John Developer', 'john@example.com', 'Developer Experience', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6'),
('550e8400-e29b-41d4-a716-446655440002', NOW(), NOW(), 'Jane Designer', 'jane@example.com', 'Design', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6'),
('550e8400-e29b-41d4-a716-446655440003', NOW(), NOW(), 'Bob Engineer', 'bob@example.com', 'Infra', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6'),
('550e8400-e29b-41d4-a716-446655440004', NOW(), NOW(), 'Alice Analyst', 'alice@example.com', '', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6');

-- Insert hackathons
INSERT INTO hackathons (id, created_at, updated_at, title, description, start_date, end_date, owner_id, status) VALUES
(1, NOW(), NOW(), 'Winter Code Fest 2025', 'A month-long coding competition focusing on innovative solutions for real-world problems. Teams will build applications using modern technologies.', '2025-12-27 09:00:00', '2026-01-27 17:00:00', '550e8400-e29b-41d4-a716-446655440000', 'active'),
(2, NOW(), NOW(), 'AI Innovation Challenge', 'Explore the possibilities of artificial intelligence and machine learning. Create projects that demonstrate practical applications of AI technologies.', '2026-02-01 10:00:00', '2026-02-28 16:00:00', '550e8400-e29b-41d4-a716-446655440000', 'upcoming'),
(3, NOW(), NOW(), 'Green Tech Hackathon', 'Build sustainable technology solutions that address environmental challenges. Focus on renewable energy, waste reduction, and eco-friendly innovations.', '2026-03-15 08:00:00', '2026-03-29 18:00:00', '550e8400-e29b-41d4-a716-446655440000', 'upcoming');

-- Insert projects
INSERT INTO projects (id,  user_id, created_at, updated_at, hackathon_id, name, description, repository_url, demo_url, status, team_members) VALUES
(1, '550e8400-e29b-41d4-a716-446655440001', NOW(), NOW(), 1, 'EcoTracker', 'Mobile app that helps users track carbon footprint and suggests reductions through gamification.', 'https://github.com/johndev/ecotracker', 'https://ecotracker-demo.herokuapp.com', 'in_progress', 2),
(2, '550e8400-e29b-41d4-a716-446655440002', NOW(), NOW(), 1, 'Smart Waste Sorter', 'IoT device that automatically sorts household waste using computer vision.', 'https://github.com/janedesign/smart-waste', 'https://smartwaste-demo.netlify.app', 'completed', 2),
(3, '550e8400-e29b-41d4-a716-446655440003', NOW(), NOW(), 1, 'Community Garden Planner', 'Web platform to plan/manage shared gardens and coordinate volunteers.', 'https://github.com/bobeng/community-garden', 'https://gardenshare.vercel.app', 'in_progress', 1),
(4, '550e8400-e29b-41d4-a716-446655440004', NOW(), NOW(), 2, 'AI Recipe Generator', 'AI-powered personalized recipes based on preferences and ingredients.', 'https://github.com/aliceai/ai-recipes', 'https://airecipes.streamlit.app', 'in_progress', 2),
(5, '550e8400-e29b-41d4-a716-446655440001', NOW(), NOW(), 2, 'Sentiment Analysis Dashboard', 'Real-time dashboard analyzing social media sentiment using NLP.', 'https://github.com/johndev/sentiment-dash', 'https://sentiment-analytics.fly.dev', 'completed', 1),
(6, '550e8400-e29b-41d4-a716-446655440002', NOW(), NOW(), 3, 'Solar Panel Optimizer', 'Smart system that maximizes solar panel efficiency via weather-driven adjustments.', 'https://github.com/janedesign/solar-opt', 'https://solaroptimizer.glitch.me', 'in_progress', 1);


-- Insert project memberships (founders are automatically added, but let's add some team members)
INSERT INTO project_memberships (id, created_at, updated_at, project_id, user_id) VALUES
(1, NOW(), NOW(), 1, '550e8400-e29b-41d4-a716-446655440001'), -- John is founder
(2, NOW(), NOW(), 1, '550e8400-e29b-41d4-a716-446655440003'), -- Bob joins John's project
(3, NOW(), NOW(), 2, '550e8400-e29b-41d4-a716-446655440002'), -- Jane is founder
(4, NOW(), NOW(), 2, '550e8400-e29b-41d4-a716-446655440004'), -- Alice joins Jane's project
(5, NOW(), NOW(), 3, '550e8400-e29b-41d4-a716-446655440003'), -- Bob is founder
(6, NOW(), NOW(), 4, '550e8400-e29b-41d4-a716-446655440004'), -- Alice is founder
(7, NOW(), NOW(), 4, '550e8400-e29b-41d4-a716-446655440002'), -- Jane joins Alice's project
(8, NOW(), NOW(), 5, '550e8400-e29b-41d4-a716-446655440001'), -- John is founder
(9, NOW(), NOW(), 6, '550e8400-e29b-41d4-a716-446655440002'); -- Jane is founder

-- Insert some sample files (using dummy data for demonstration)
INSERT INTO files (id, created_at, updated_at, filename, data, content_type, size, user_id, hackathon_id, project_id) VALUES
(1, NOW(), NOW(), 'ecotracker-mockup.png', decode('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==', 'base64'), 'image/png', 68, '550e8400-e29b-41d4-a716-446655440001', 1, 1),
(2, NOW(), NOW(), 'waste_sorter_specs.pdf', decode('JVBERi0xLjQKMSAwIG9iago8PC9UeXBlIC9DYXRhbG9nCi9QYWdlcyAyIDAgUgo+PgplbmRvYmoKMiAwIG9iago8PC9UeXBlIC9QYWdlcwovS2lkcyBbMyAwIFJdCi9Db3VudCAxCj4+CmVuZG9iagozIDAgb2JqCjw8L1R5cGUgL1BhZ2UKL1BhcmVudCAyIDAgUgovTWVkaWFCb3ggWzAgMCA2MTIgNzkyXQovQ29udGVudHMgNCAwIFIKPj4KZW5kb2JqCjQgMCBvYmoKPDwvTGVuZ3RoIDQ0Pj4Kc3RyZWFtCkJUCjEwIDUwIFRECi9GMSAxMiBUZgooSGVsbG8gV29ybGQpIFRqCkVUCmVuZHN0cmVhbQplbmRvYmoKeHJlZgowIDUKMDAwMDAwMDAwMCA2NTUzNSBmIAowMDAwMDAwMDA5IDAwMDAwIG4gCjAwMDAwMDAwNTkgMDAwMDAgbiAKMDAwMDAwMDEwMyAwMDAwMCBuIAp0cmFpbGVyCjw8L1NpemUgNQovUm9vdCAxIDAgUgo+PgpzdGFydHhyZWYKNzE5CmUlRU9GCg==', 'base64'), 'application/pdf', 719, '550e8400-e29b-41d4-a716-446655440002', 1, 2);