BEGIN;

-- Seed roles (id auto-generated)
WITH inserted_roles AS (
	INSERT INTO roles (name, description)
	VALUES
		('ADMIN', 'Administrator role'),
		('USER', 'Standard user role')
	ON CONFLICT (name) DO NOTHING
	RETURNING id, name
),
roles_final AS (
	-- Ensure we have IDs for both roles even if they already existed
	SELECT id, name FROM inserted_roles
	UNION ALL
	SELECT r.id, r.name FROM roles r WHERE r.name IN ('ADMIN','USER') AND NOT EXISTS (
		SELECT 1 FROM inserted_roles ir WHERE ir.name = r.name
	)
)
-- Seed users referencing roles
INSERT INTO users (first_name, last_name, email, phone, password, is_verified, role_id)
SELECT 'System', 'Admin', 'admin@coffee.local', '+10000000001', crypt('Admin@123', gen_salt('bf')), true,
	(SELECT id FROM roles_final WHERE name = 'ADMIN')
UNION ALL
SELECT 'Test', 'User', 'user@coffee.local', '+10000000002', crypt('User@123', gen_salt('bf')), true,
	(SELECT id FROM roles_final WHERE name = 'USER')
ON CONFLICT (email) DO NOTHING;

COMMIT;
