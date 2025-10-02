BEGIN;

-- Revert event_samples table to make Price and Rating NOT NULL
ALTER TABLE event_samples 
ALTER COLUMN price SET NOT NULL,
ALTER COLUMN rating SET NOT NULL;

COMMIT;
