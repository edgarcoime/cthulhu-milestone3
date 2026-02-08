CREATE TABLE IF NOT EXISTS lifecycle (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bucket_slug TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_lifecycle_bucket_slug ON lifecycle (bucket_slug);