BEGIN;

-- Remove seeded users by email
DELETE FROM users WHERE email IN ('admin@coffee.local', 'user@coffee.local');

-- Optionally remove roles if you want to revert roles too (uncomment if desired)
-- DELETE FROM roles WHERE name IN ('ADMIN', 'USER');

COMMIT;
