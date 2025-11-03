-- Migration: 000002_additional_tables (rollback)
-- Description: Rollback additional tables

DROP TABLE IF EXISTS factorial_current_calculated_numbers;
DROP TABLE IF EXISTS factorial_max_request_numbers;

