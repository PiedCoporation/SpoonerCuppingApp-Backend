BEGIN;

ALTER TABLE events 
    DROP COLUMN IF EXISTS register_start_time,
    DROP COLUMN IF EXISTS register_end_time;

COMMIT;


