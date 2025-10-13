BEGIN;

ALTER TABLE users 
    DROP COLUMN IF EXISTS forgot_password_code,
    DROP COLUMN IF EXISTS forgot_password_expires_at;

COMMIT;


