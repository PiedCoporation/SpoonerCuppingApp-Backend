BEGIN;

INSERT INTO events (
    id,
    name,
    date_of_event,
    start_time,
    end_time,
    "limit",
    total_joined,
    number_samples,
    phone_contact,
    email_contact,
    is_public,
    user_id,
    register_date,
    register_start_time,
    register_end_time,
    auto_accept,
    total_registered
)
SELECT
    '9c2c1a43-828d-4b3e-8f8b-df2e4d8438a1'::uuid,
    'API Test Event Participants',
    NOW() + INTERVAL '7 day',
    NOW() + INTERVAL '7 day',
    NOW() + INTERVAL '7 day' + INTERVAL '2 hour',
    20,
    1,
    3,
    '+10000009999',
    'admin@coffee.local',
    true,
    u.id,
    NOW(),
    NOW() - INTERVAL '7 day',
    NOW() + INTERVAL '7 day',
    false,
    3
FROM users u
WHERE u.email = 'admin@coffee.local'
  AND NOT EXISTS (
      SELECT 1 FROM events WHERE id = '9c2c1a43-828d-4b3e-8f8b-df2e4d8438a1'::uuid
  );

INSERT INTO event_addresses (
    id,
    province,
    district,
    longitude,
    latitude,
    ward,
    street,
    phone,
    event_id
)
SELECT
    '5e6f7a80-9123-4b5c-cdef-2345678901bc'::uuid,
    'Ho Chi Minh City',
    'District 1',
    '106.7000',
    '10.7769',
    'Ben Nghe',
    '123 API Test St',
    '+10000009999',
    '9c2c1a43-828d-4b3e-8f8b-df2e4d8438a1'::uuid
WHERE NOT EXISTS (
    SELECT 1 FROM event_addresses WHERE id = '5e6f7a80-9123-4b5c-cdef-2345678901bc'::uuid
);

INSERT INTO event_users (
    id,
    is_accepted,
    is_invited,
    user_id,
    event_id,
    is_host
)
SELECT
    '1f4a6a20-4dc7-4b36-8bd2-f84e8a3b73df'::uuid,
    true,
    false,
    u.id,
    '9c2c1a43-828d-4b3e-8f8b-df2e4d8438a1'::uuid,
    true
FROM users u
WHERE u.email = 'admin@coffee.local'
  AND NOT EXISTS (
      SELECT 1 FROM event_users WHERE id = '1f4a6a20-4dc7-4b36-8bd2-f84e8a3b73df'::uuid
  );

INSERT INTO event_users (
    id,
    is_accepted,
    is_invited,
    user_id,
    event_id,
    is_host
)
SELECT
    '2a6b9adf-3c47-4abf-9b8e-e3c62f7d1101'::uuid,
    true,
    false,
    u.id,
    '9c2c1a43-828d-4b3e-8f8b-df2e4d8438a1'::uuid,
    false
FROM users u
WHERE u.email = 'quocthai@gmail.com'
  AND NOT EXISTS (
      SELECT 1 FROM event_users WHERE id = '2a6b9adf-3c47-4abf-9b8e-e3c62f7d1101'::uuid
  );

INSERT INTO event_users (
    id,
    is_accepted,
    is_invited,
    user_id,
    event_id,
    is_host
)
SELECT
    '3c1d2e3f-4a5b-478c-9d0e-1f2a3b4c5d6e'::uuid,
    false,
    false,
    u.id,
    '9c2c1a43-828d-4b3e-8f8b-df2e4d8438a1'::uuid,
    false
FROM users u
WHERE u.email = 'lediep@gmail.com'
  AND NOT EXISTS (
      SELECT 1 FROM event_users WHERE id = '3c1d2e3f-4a5b-478c-9d0e-1f2a3b4c5d6e'::uuid
  );

INSERT INTO event_users (
    id,
    is_accepted,
    is_invited,
    user_id,
    event_id,
    is_host
)
SELECT
    '4d5e6f70-8192-4a3b-bcde-1234567890ab'::uuid,
    false,
    true,
    u.id,
    '9c2c1a43-828d-4b3e-8f8b-df2e4d8438a1'::uuid,
    false
FROM users u
WHERE u.email = 'bentran@gmail.com'
  AND NOT EXISTS (
      SELECT 1 FROM event_users WHERE id = '4d5e6f70-8192-4a3b-bcde-1234567890ab'::uuid
  );

COMMIT;

