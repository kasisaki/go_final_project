CREATE TABLE IF NOT EXISTS scheduler
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    VARCHAR(8),
    title   VARCHAR(128) NOT NULL,
    comment TEXT,
    repeat  VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
