BEGIN;

-- Remove new columns from events table
ALTER TABLE events 
DROP COLUMN IF EXISTS register_date,
DROP COLUMN IF EXISTS register_status;

-- Drop register_status_enum type if it exists
DO $$
BEGIN
	IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'register_status_enum') THEN
		DROP TYPE register_status_enum;
	END IF;
END $$;

COMMIT;
