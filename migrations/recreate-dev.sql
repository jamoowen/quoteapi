DELETE FROM api_keys;

INSERT INTO api_keys (
    email,
    api_key,
    expiry,
    request_count,
    created_at,
    last_used_at
) VALUES (
    'test@example.com',
    'test_key_123',
    strftime('%s', 'now', '+1 year'),  -- One year from now
    0,
    strftime('%s', 'now'),             -- Current timestamp
    strftime('%s', 'now')              -- Current timestamp
);
