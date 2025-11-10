BEGIN;

ALTER TABLE events
    ADD COLUMN IF NOT EXISTS auto_accept boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS total_registered integer NOT NULL DEFAULT 0;

ALTER TABLE events
    RENAME COLUMN total_current TO total_joined;

COMMIT;
