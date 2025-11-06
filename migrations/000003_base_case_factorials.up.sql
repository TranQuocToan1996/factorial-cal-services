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
INSERT INTO factorial_calculations (number, status, s3_key, checksum, size, bucket)
VALUES (
    0,
    'done',
    '0.txt',
    '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b',
    1 ,
    'factorial-calculator-service'
)
ON CONFLICT (number) DO NOTHING;

-- Factorial(1) = 1
INSERT INTO factorial_calculations (number, status, s3_key, checksum, size, bucket)
VALUES (
    1,
    'done',
    '1.txt',
    '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b',
    1 ,
    'factorial-calculator-service'
)
ON CONFLICT (number) DO NOTHING;

-- Factorial(2) = 2
INSERT INTO factorial_calculations (number, status, s3_key, checksum, size, bucket)
VALUES (
    2,
    'done',
    '2.txt',
    'd4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35',
    1 ,
    'factorial-calculator-service'
)
ON CONFLICT (number) DO NOTHING;

-- Factorial(3) = 6
INSERT INTO factorial_calculations (number, status, s3_key, checksum, size, bucket)
VALUES (
    3,
    'done',
    '3.txt',
    'e7f6c011776e8db7cd330b54174fd76f7d0216b612387a5ffcfb81e6f0919683',
    1 , 
    'factorial-calculator-service'
)
ON CONFLICT (number) DO NOTHING;

INSERT INTO factorial_current_calculated_numbers (next_number) VALUES (4);
INSERT INTO factorial_max_request_numbers (max_number) VALUES (3);
