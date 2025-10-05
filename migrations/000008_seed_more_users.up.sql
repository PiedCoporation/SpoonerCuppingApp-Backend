BEGIN;

-- Ensure USER role exists and get its id
WITH inserted_role AS (
    INSERT INTO roles (name, description)
    VALUES ('USER', 'Standard user role')
    ON CONFLICT (name) DO NOTHING
    RETURNING id, name
),
role_final AS (
    SELECT id, name FROM inserted_role
    UNION ALL
    SELECT r.id, r.name FROM roles r WHERE r.name = 'USER' AND NOT EXISTS (
        SELECT 1 FROM inserted_role ir WHERE ir.name = r.name
    )
)
INSERT INTO users (first_name, last_name, email, phone, password, is_verified, role_id)
SELECT 'Quoc', 'Thai', 'quocthai@gmail.com', NULL, crypt('User@123', gen_salt('bf')), true,
       (SELECT id FROM role_final WHERE name = 'USER')
UNION ALL
SELECT 'Le', 'Diep', 'lediep@gmail.com', NULL, crypt('User@123', gen_salt('bf')), true,
       (SELECT id FROM role_final WHERE name = 'USER')
UNION ALL
SELECT 'Ben', 'Tran', 'bentran@gmail.com', NULL, crypt('User@123', gen_salt('bf')), true,
       (SELECT id FROM role_final WHERE name = 'USER')
ON CONFLICT (email) DO NOTHING;

COMMIT;


