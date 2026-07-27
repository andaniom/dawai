-- Fix: Seed data password hashes were incorrect (bcrypt compare failed)
-- Password: "password", cost: 12
-- Hash: $2a$12$YWvI08JumEIoDx9kHP8eZuVyXddedNCumY0a6QB9LEuSxSQhbQD6K

UPDATE users SET password_hash = '$2a$12$YWvI08JumEIoDx9kHP8eZuVyXddedNCumY0a6QB9LEuSxSQhbQD6K'
WHERE email IN ('qaadmin@test.local', 'teacher@test.local', 'student@test.local', 'qateacher@test.local', 'qastudent@test.local');
