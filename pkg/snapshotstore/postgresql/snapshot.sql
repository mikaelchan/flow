CREATE TABLE IF NOT EXISTS snapshots
(
    stream_id VARCHAR(255) PRIMARY KEY,
    type      VARCHAR(255),
    version   BIGINT,
    state     BYTEA
);

CREATE INDEX idx_snapshots_type ON snapshots (type);