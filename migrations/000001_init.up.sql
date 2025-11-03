-- Migration: 000001_init
-- Description: Initial database schema for factorial calculations
-- Cross-DB compatible (PostgreSQL/MySQL)
-- Note: Auto-increment and updated_at are handled by GORM

CREATE TABLE IF NOT EXISTS factorial_calculations (
    id BIGINT NOT NULL,
    number VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    s3_key VARCHAR(512) NOT NULL,
    checksum VARCHAR(64),
    size BIGINT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (number)
);

-- Create indexes (works for both PostgreSQL and MySQL)
CREATE INDEX idx_factorial_number ON factorial_calculations (number);
CREATE INDEX idx_factorial_status ON factorial_calculations (status);
CREATE INDEX idx_factorial_created_at ON factorial_calculations (created_at);
