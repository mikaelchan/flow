CREATE TABLE IF NOT EXISTS events (
    stream_id VARCHAR(255) NOT NULL,
    event_id VARCHAR(255) NOT NULL PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL,
    payload BYTEA NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_stream_id ON events (stream_id);
CREATE INDEX IF NOT EXISTS idx_events_event_type ON events (event_type);