locals {
  common_tags = {
    project    = var.project
    env        = var.environment
    department = var.department
    owner      = var.owner
  }

  rds_subnet_names = {
    "10.0.16.0/25"   = "RDS-Pvt-subnet-1"
    "10.0.16.128/25" = "RDS-Pvt-subnet-2"
    "10.0.17.0/25"   = "RDS-Pvt-subnet-3"
    "10.0.17.128/25" = "RDS-Pvt-subnet-4"
    "10.0.18.0/25"   = "RDS-Pvt-subnet-5"
  }
}


