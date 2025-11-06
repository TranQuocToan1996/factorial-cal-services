-- Migration: 000001_init
-- Description: Initial database schema for factorial calculations
-- PostgreSQL

CREATE TABLE IF NOT EXISTS factorial_calculations (
    id BIGSERIAL PRIMARY KEY,
    number BIGINT NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    s3_key VARCHAR(255),
    bucket VARCHAR(255),
    checksum VARCHAR(64),
    size BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT factorial_calculations_number_unique UNIQUE (number)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_factorial_status ON factorial_calculations (status);
CREATE INDEX IF NOT EXISTS idx_factorial_created_at ON factorial_calculations (created_at);
