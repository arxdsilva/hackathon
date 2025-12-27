-- Seed data for hackathon application

-- Insert users - pass hashed password for 'password'
INSERT INTO users (id, created_at, updated_at, name, email, role, password_hash) VALUES
('550e8400-e29b-41d4-a716-446655440000', NOW(), NOW(), 'Admin Owner', 'admin@example.com', 'owner', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6'),
('550e8400-e29b-41d4-a716-446655440001', NOW(), NOW(), 'John Developer', 'john@example.com', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6'),
('550e8400-e29b-41d4-a716-446655440002', NOW(), NOW(), 'Jane Designer', 'jane@example.com', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6'),
('550e8400-e29b-41d4-a716-446655440003', NOW(), NOW(), 'Bob Engineer', 'bob@example.com', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6'),
('550e8400-e29b-41d4-a716-446655440004', NOW(), NOW(), 'Alice Analyst', 'alice@example.com', 'hacker', '$2a$10$izUpv9c8SjPd2zLsIMaAmOtQdxOYE5QKPcyZjexgkn441Cmf22Wk6');

-- Insert hackathons
INSERT INTO hackathons (id, created_at, updated_at, title, description, start_date, end_date, owner_id, status, schedule) VALUES
(1, NOW(), NOW(), 'Winter Code Fest 2025', 'A month-long coding competition focusing on innovative solutions for real-world problems. Teams will build applications using modern technologies.', '2025-12-27 09:00:00', '2026-01-27 17:00:00', '550e8400-e29b-41d4-a716-446655440000', 'active', 'Day 1: Opening ceremony and team formation\nDay 7: Mid-point check-in\nDay 14: Mentor sessions\nDay 21: Final submissions\nDay 28: Presentations and awards'),
(2, NOW(), NOW(), 'AI Innovation Challenge', 'Explore the possibilities of artificial intelligence and machine learning. Create projects that demonstrate practical applications of AI technologies.', '2026-02-01 10:00:00', '2026-02-28 16:00:00', '550e8400-e29b-41d4-a716-446655440000', 'upcoming', 'Week 1: AI fundamentals workshop\nWeek 2: Project development\nWeek 3: AI ethics discussion\nWeek 4: Final presentations'),
(3, NOW(), NOW(), 'Green Tech Hackathon', 'Build sustainable technology solutions that address environmental challenges. Focus on renewable energy, waste reduction, and eco-friendly innovations.', '2026-03-15 08:00:00', '2026-03-29 18:00:00', '550e8400-e29b-41d4-a716-446655440000', 'upcoming', 'Opening: Environmental impact workshop\nMid-event: Sustainability expert talks\nClosing: Green tech showcase');

-- Insert projects
INSERT INTO projects (id, created_at, updated_at, hackathon_id, user_id, name, description, repository_url, demo_url, status) VALUES
(1, NOW(), NOW(), 1, '550e8400-e29b-41d4-a716-446655440001', 'EcoTracker', 'A mobile app that helps users track their carbon footprint and suggests ways to reduce environmental impact through gamification.', 'https://github.com/johndev/ecotracker', 'https://ecotracker-demo.herokuapp.com', 'in_progress'),
(2, NOW(), NOW(), 1, '550e8400-e29b-41d4-a716-446655440002', 'Smart Waste Sorter', 'An IoT device that automatically sorts household waste using computer vision and machine learning algorithms.', 'https://github.com/janedesign/smart-waste', 'https://smartwaste-demo.netlify.app', 'completed'),
(3, NOW(), NOW(), 1, '550e8400-e29b-41d4-a716-446655440003', 'Community Garden Planner', 'A web platform that helps communities plan and manage shared gardens, track plant growth, and coordinate volunteer activities.', 'https://github.com/bobeng/community-garden', 'https://gardenshare.vercel.app', 'in_progress'),
(4, NOW(), NOW(), 2, '550e8400-e29b-41d4-a716-446655440004', 'AI Recipe Generator', 'An AI-powered application that creates personalized recipes based on dietary preferences, available ingredients, and nutritional goals.', 'https://github.com/aliceai/ai-recipes', 'https://airecipes.streamlit.app', 'in_progress'),
(5, NOW(), NOW(), 2, '550e8400-e29b-41d4-a716-446655440001', 'Sentiment Analysis Dashboard', 'A real-time dashboard that analyzes social media sentiment about brands and products using natural language processing.', 'https://github.com/johndev/sentiment-dash', 'https://sentiment-analytics.fly.dev', 'completed'),
(6, NOW(), NOW(), 3, '550e8400-e29b-41d4-a716-446655440002', 'Solar Panel Optimizer', 'A smart system that maximizes solar panel efficiency by adjusting angles and cleaning schedules based on weather data and AI predictions.', 'https://github.com/janedesign/solar-opt', 'https://solaroptimizer.glitch.me', 'in_progress');

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