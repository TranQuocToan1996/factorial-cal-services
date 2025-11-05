-- Migration: 000003_base_case_factorials
-- Description: Insert base cases for factorial calculations (0, 1, 2, 3)
-- Note: S3 files will be uploaded manually
-- PostgreSQL

-- Insert factorial base cases (0, 1, 2, 3)
-- Factorial values:
--   0! = 1
--   1! = 1
--   2! = 2
--   3! = 6

-- Factorial(0) = 1
INSERT INTO factorial_calculations (number, status, s3_key, checksum, size)
VALUES (
    0,
    'done',
    '0.txt',
    '',
    1  -- "1" is 1 byte
)
ON CONFLICT (number) DO NOTHING;

-- Factorial(1) = 1
INSERT INTO factorial_calculations (number, status, s3_key, checksum, size)
VALUES (
    1,
    'done',
    '1.txt',
    '',
    1  -- "1" is 1 byte
)
ON CONFLICT (number) DO NOTHING;

-- Factorial(2) = 2
INSERT INTO factorial_calculations (number, status, s3_key, checksum, size)
VALUES (
    2,
    'done',
    '2.txt',
    '',
    1  -- "2" is 1 byte
)
ON CONFLICT (number) DO NOTHING;

-- Factorial(3) = 6
INSERT INTO factorial_calculations (number, status, s3_key, checksum, size)
VALUES (
    3,
    'done',
    '3.txt',
    '',
    1  -- "6" is 1 byte
)
ON CONFLICT (number) DO NOTHING;

-- Update current calculated number to 4 (since 0, 1, 2, 3 are already done)
-- This ensures the calculator starts calculating from 4
UPDATE factorial_current_calculated_numbers 
SET cur_number = 4 
WHERE id IN (SELECT id FROM factorial_current_calculated_numbers ORDER BY cur_number DESC LIMIT 1);

-- If no record exists, create one
INSERT INTO factorial_current_calculated_numbers (cur_number)
SELECT 4
WHERE NOT EXISTS (SELECT 1 FROM factorial_current_calculated_numbers);
