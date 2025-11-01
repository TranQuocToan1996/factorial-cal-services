-- Migration: 000001_init
-- Description: Initial database schema for factorial calculations

CREATE TABLE factorial_calculations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    number VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    s3_key VARCHAR(512) NOT NULL,
    checksum VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_number (number),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
