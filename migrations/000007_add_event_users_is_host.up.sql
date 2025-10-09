BEGIN;

ALTER TABLE event_users 
	ADD COLUMN is_host boolean NOT NULL DEFAULT false;

COMMIT;
