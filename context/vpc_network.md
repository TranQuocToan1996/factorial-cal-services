# VPC Network Overview

This VPC is a non-default VPC designed for a production environment with segregated public and private subnets across multiple Availability Zones in `us-east-1`. It supports application workloads (public/private), RabbitMQ, RDS (Postgres) in private subnets, and standard internet egress via an Internet Gateway for public subnets.

## VPC
- CIDR: `10.0.0.0/16`
- VPC ID (source data): `vpc-0956e9914efe15691`
- Tenancy: default
- Tags indicate production environment and project context (e.g., department=dev, env=production, owner=MrT)

## Internet Gateway
- IGW ID (source data): `igw-03ec83faa21e2b617`
- Attached to the VPC and used by the public route table for 0.0.0.0/0 internet access

## Subnets
The VPC contains multiple subnets, including application public/private subnets and RDS private subnets distributed across AZs.

- Public subnets (application ingress):
  - us-east-1a: `10.0.0.0/22` (Name: express-nodejs-demo-subnet-public1-us-east-1a)
  - us-east-1b: `10.0.4.0/22` (Name: express-nodejs-demo-subnet-public2-us-east-1b)

- Private subnets (application/private services):
  - us-east-1a: `10.0.8.0/22` (Name: express-nodejs-demo-subnet-private1-us-east-1a)
  - us-east-1b: `10.0.12.0/22` (Name: express-nodejs-demo-subnet-private2-us-east-1b)

- RDS private subnets:
  - `10.0.16.0/25` (RDS-Pvt-subnet-1)
  - `10.0.16.128/25` (RDS-Pvt-subnet-2)
  - `10.0.17.0/25` (RDS-Pvt-subnet-3)
  - `10.0.17.128/25` (RDS-Pvt-subnet-4)
  - `10.0.18.0/25` (RDS-Pvt-subnet-5)

## Route Tables
- Public Route Table:
  - Routes: local (10.0.0.0/16), default route `0.0.0.0/0` via IGW
  - Associated with both public subnets

- Private Route Tables:
  - Routes: local (10.0.0.0/16) only (per source data). Suitable for RDS and private application subnets. If NAT egress is required for private subnets, add a NAT Gateway and default route to it.

## Security Groups (intended usage summary)
- Application/Public SG: inbound 80/443/8080/3000 from controlled CIDRs; egress open
- RabbitMQ SG: inbound 5671/5672/15672 from controlled CIDRs; egress open
- Postgres SG: inbound 5432 from controlled CIDRs or from application SG; egress open
- Admin/Dev SG: inbound 22 from office IPs

Note: The source JSON references AWS account-specific Security Group IDs in `UserIdGroupPairs`. In Terraform, reference SGs by resource, or model ingress via CIDRs/peer SG attachments rather than hardcoding IDs.

## Notes
- Public subnets are internet-routable through the IGW; private/RDS subnets are isolated. Ensure NAT Gateways if private subnets need outbound internet.
- Tagging is standardized across resources (department, env, project, owner, name). Terraform definitions keep tags centralized.