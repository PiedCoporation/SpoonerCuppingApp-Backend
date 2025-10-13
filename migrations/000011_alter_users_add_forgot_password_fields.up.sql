BEGIN;

ALTER TABLE users 
    ADD COLUMN IF NOT EXISTS forgot_password_code text NULL,
    ADD COLUMN IF NOT EXISTS forgot_password_expires_at timestamptz NULL;

COMMIT;


