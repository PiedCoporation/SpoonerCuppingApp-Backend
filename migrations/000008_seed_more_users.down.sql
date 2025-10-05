BEGIN;

DELETE FROM users
WHERE email IN (
    'quocthai@gmail.com',
    'lediep@gmail.com',
    'bentran@gmail.com'
);

COMMIT;


