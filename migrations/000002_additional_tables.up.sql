-- Migration: 000002_additional_tables
-- Description: Additional tables for factorial tracking
-- PostgreSQL

-- Table to track the maximum requested factorial number
CREATE TABLE IF NOT EXISTS factorial_max_request_numbers (
    id BIGSERIAL PRIMARY KEY,
    max_number BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table to track the current calculated factorial number
CREATE TABLE IF NOT EXISTS factorial_current_calculated_numbers (
    id BIGSERIAL PRIMARY KEY,
    next_number BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_max_request_number ON factorial_max_request_numbers (max_number);
CREATE INDEX IF NOT EXISTS idx_current_calculated_number ON factorial_current_calculated_numbers (next_number);

