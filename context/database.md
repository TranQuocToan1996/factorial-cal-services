# Describe the DB use
# TODO: add if not exist
CREATE TABLE factorial_calculations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    number VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    s3_key VARCHAR(512) NOT NULL,
    checksum VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_number (number),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
)

CREATE TABLE factorial_metadatas (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    number VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    s3_key VARCHAR(512) NOT NULL,
    checksum VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE factorial_max_request_numbers (
    max_number VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE factorial_current_calculated_numbers (
    cur_number VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)