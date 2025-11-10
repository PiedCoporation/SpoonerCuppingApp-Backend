BEGIN;

DELETE FROM event_users
WHERE id IN (
    '1f4a6a20-4dc7-4b36-8bd2-f84e8a3b73df'::uuid,
    '2a6b9adf-3c47-4abf-9b8e-e3c62f7d1101'::uuid,
    '3c1d2e3f-4a5b-478c-9d0e-1f2a3b4c5d6e'::uuid,
    '4d5e6f70-8192-4a3b-bcde-1234567890ab'::uuid
);

DELETE FROM event_addresses
WHERE id = '5e6f7a80-9123-4b5c-cdef-2345678901bc'::uuid;

DELETE FROM events
WHERE id = '9c2c1a43-828d-4b3e-8f8b-df2e4d8438a1'::uuid;

COMMIT;

