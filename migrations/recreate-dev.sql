DELETE FROM users;

INSERT INTO users (
    email,
    hashed_api_key,
    request_count,
    created_at_timestamp,
    last_used_at
) VALUES (
    'test@example.com',
    'test_key_123',
    0,
    strftime('%s', 'now'),             -- Current timestamp
    strftime('%s', 'now')              -- Current timestamp
);
