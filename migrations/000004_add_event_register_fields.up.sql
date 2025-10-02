BEGIN;

-- Create register_status_enum type
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'register_status_enum') THEN
		CREATE TYPE register_status_enum AS ENUM (
			'PENDING',
			'ACCEPTED',
			'FULL'
		);
	END IF;
END $$;

-- Add new columns to events table
ALTER TABLE events 
ADD COLUMN register_date timestamptz NOT NULL DEFAULT now(),
ADD COLUMN register_status register_status_enum NOT NULL DEFAULT 'PENDING';

COMMIT;
