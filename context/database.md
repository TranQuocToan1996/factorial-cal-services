# Database Schema - PostgreSQL

## Main Tables

### factorial_calculations
Stores factorial calculation metadata and results.

```sql
CREATE TABLE IF NOT EXISTS factorial_calculations (
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    s3_key VARCHAR(512) NOT NULL,
    checksum VARCHAR(64),
    size BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_factorial_status ON factorial_calculations (status);
CREATE INDEX IF NOT EXISTS idx_factorial_created_at ON factorial_calculations (created_at);
```

### factorial_max_request_numbers
Tracks the maximum requested factorial number for incremental calculation optimization.

```sql
CREATE TABLE IF NOT EXISTS factorial_max_request_numbers (
    id BIGSERIAL PRIMARY KEY,
    max_number VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_max_request_number ON factorial_max_request_numbers (max_number);
```

### factorial_current_calculated_numbers
Tracks the current calculated factorial number for incremental calculation.

```sql
CREATE TABLE IF NOT EXISTS factorial_current_calculated_numbers (
    id BIGSERIAL PRIMARY KEY,
    cur_number VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_current_calculated_number ON factorial_current_calculated_numbers (cur_number);
```

## Notes
- Using BIGSERIAL for auto-incrementing primary keys (PostgreSQL)
- Using TIMESTAMPTZ for timezone-aware timestamps
- Unique constraint on `number` field to prevent duplicate calculations
- Indexes on status and created_at for efficient querying