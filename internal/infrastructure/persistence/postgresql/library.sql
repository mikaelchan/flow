CREATE TABLE IF NOT EXISTS libraries (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    media_type INT NOT NULL,
    location VARCHAR(255) NOT NULL,
    quality_preference JSONB NOT NULL,
    naming_template VARCHAR(255) NOT NULL,
    status INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_libraries_name ON libraries (name);
CREATE INDEX IF NOT EXISTS idx_libraries_location ON libraries (location);