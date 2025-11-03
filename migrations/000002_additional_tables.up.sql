-- Migration: 000002_additional_tables
-- Description: Additional tables for factorial tracking
-- Cross-DB compatible (PostgreSQL/MySQL)

-- Table to track the maximum requested factorial number
CREATE TABLE IF NOT EXISTS factorial_max_request_numbers (
    id BIGINT NOT NULL,
    max_number VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Table to track the current calculated factorial number
CREATE TABLE IF NOT EXISTS factorial_current_calculated_numbers (
    id BIGINT NOT NULL,
    cur_number VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Create indexes
CREATE INDEX idx_max_request_number ON factorial_max_request_numbers (max_number);
CREATE INDEX idx_current_calculated_number ON factorial_current_calculated_numbers (cur_number);

