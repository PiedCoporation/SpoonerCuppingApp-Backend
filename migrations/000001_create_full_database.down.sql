BEGIN;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS tasting_notes;
DROP TABLE IF EXISTS user_tastings;
DROP TABLE IF EXISTS event_users;
DROP TABLE IF EXISTS event_samples;
DROP TABLE IF EXISTS user_samples;
DROP TABLE IF EXISTS event_addresses;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;

-- Drop custom types if exist
DO $$
BEGIN
	IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'roastinglever_enum') THEN
		DROP TYPE roastinglever_enum;
	END IF;
	IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'processing_enum') THEN
		DROP TYPE processing_enum;
	END IF;
END $$;

COMMIT;
