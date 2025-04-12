CREATE TABLE IF NOT EXISTS api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    api_key TEXT NOT NULL UNIQUE,
    expiry INTEGER,           -- Unix timestamp
    request_count INTEGER DEFAULT 0,
    created_at INTEGER,       -- Unix timestamp
    last_used_at INTEGER     -- Unix timestamp
);

CREATE INDEX IF NOT EXISTS idx_api_keys_email ON api_keys(email);
CREATE INDEX IF NOT EXISTS idx_api_keys_expiry ON api_keys(expiry);
