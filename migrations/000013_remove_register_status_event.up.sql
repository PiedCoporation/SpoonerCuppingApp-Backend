BEGIN;

ALTER TABLE events 
    DROP COLUMN IF EXISTS register_status;

COMMIT;


    