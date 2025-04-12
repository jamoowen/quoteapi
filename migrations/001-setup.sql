CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    hashed_api_key TEXT NOT NULL UNIQUE,
    request_count INTEGER DEFAULT 0,
    created_at_timestamp INTEGER DEFAULT (strftime('%s', 'now')),  -- Unix timestamp in seconds
    last_used_at DATE     
);
