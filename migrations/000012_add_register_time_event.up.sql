BEGIN;

ALTER TABLE events 
    ADD COLUMN IF NOT EXISTS register_start_time timestamptz NULL,
    ADD COLUMN IF NOT EXISTS register_end_time timestamptz NULL;

COMMIT;


