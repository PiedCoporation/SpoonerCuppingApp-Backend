BEGIN;

ALTER TABLE events 
    ADD COLUMN register_status register_status_enum NOT NULL DEFAULT 'PENDING';

COMMIT;