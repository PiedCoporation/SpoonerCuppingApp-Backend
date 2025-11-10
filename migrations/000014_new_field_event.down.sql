BEGIN;

ALTER TABLE events
    RENAME COLUMN total_joined TO total_current;

ALTER TABLE events
    DROP COLUMN IF EXISTS auto_accept,
    DROP COLUMN IF EXISTS total_registered;

COMMIT;
