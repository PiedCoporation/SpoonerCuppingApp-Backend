BEGIN;

-- Alter event_samples table to make Price and Rating nullable
ALTER TABLE event_samples 
ALTER COLUMN price DROP NOT NULL,
ALTER COLUMN rating DROP NOT NULL;

COMMIT;
