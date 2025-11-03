# Project History

## 2025-11-03T08:58:32Z - PostgreSQL Migration and Code Review Fixes
- **Converted all database schemas to PostgreSQL**:
  - Updated migration files to use BIGSERIAL for auto-increment primary keys
  - Changed TIMESTAMP to TIMESTAMPTZ for timezone-aware timestamps
  - Removed MySQL-specific syntax (AUTO_INCREMENT, ENGINE=InnoDB)
  - Added IF NOT EXISTS to CREATE INDEX statements
  - Removed duplicate index on `number` field (already covered by UNIQUE constraint)
  
- **Added database connection pooling**:
  - Configured max idle connections: 10
  - Configured max open connections: 100
  - Set connection max lifetime: 1 hour
  - Added GORM logger for better debugging

- **Updated documentation and configuration**:
  - Updated context/database.md with proper PostgreSQL schema documentation
  - Updated context/arch.md to reflect PostgreSQL usage
  - Updated TODO.md to reflect PostgreSQL migration completion
  - Updated infrastructure/helm/values.yaml to use postgres-service instead of mysql-service
  - Changed default port from 3306 to 5432

- **Code quality improvements**:
  - All code now compiles without errors
  - No linter errors found
  - Proper PostgreSQL DSN format used throughout
  - Consistent use of postgres driver in pkg/db/gorm.go

## Previous Updates (2025-11-03)
- Moved Helm charts from helm/ to infrastructure/ directory
- Created comprehensive VPC network documentation
- Generated Terraform files for AWS VPC infrastructure (VPC, subnets, route tables, security groups, IGW)
- Fixed Terraform files to match actual AWS infrastructure from JSON exports
- Implemented batch processing for RabbitMQ consumer
- Added checksum calculation service (SHA256)
- Created repositories for max_request and current_calculated tracking
- Updated API response format to match design specification
- Fixed GET /factorial flow to return "calculating" status instead of 404
