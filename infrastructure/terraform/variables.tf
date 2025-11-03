variable "region" {
  type        = string
  description = "AWS region"
  default     = "us-east-1"
}

variable "project" {
  type        = string
  description = "Project tag"
  default     = "express-nodejs-demo"
}

variable "environment" {
  type        = string
  description = "Environment tag"
  default     = "production"
}

variable "department" {
  type        = string
  description = "Department tag"
  default     = "dev"
}

variable "owner" {
  type        = string
  description = "Owner tag"
  default     = "MrT"
}

variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR block"
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  type = object({
    a = string
    b = string
  })
  description = "Public subnet CIDRs by AZ suffix"
  default = {
    a = "10.0.0.0/22"
    b = "10.0.4.0/22"
  }
}

variable "private_subnet_cidrs" {
  type = object({
    a = string
    b = string
  })
  description = "Private app subnet CIDRs by AZ suffix"
  default = {
    a = "10.0.8.0/22"
    b = "10.0.12.0/22"
  }
}

variable "rds_private_subnet_cidrs" {
  type        = list(string)
  description = "RDS private subnet CIDRs"
  default     = ["10.0.16.0/25", "10.0.16.128/25", "10.0.17.0/25", "10.0.17.128/25", "10.0.18.0/25"]
}

variable "availability_zones" {
  type = object({
    a = string
    b = string
    c = string
    d = string
    e = string
  })
  description = "AZ names mapping"
  default = {
    a = "us-east-1a"
    b = "us-east-1b"
    c = "us-east-1c"
    d = "us-east-1d"
    e = "us-east-1e"
  }
}

variable "rds_subnet_az_mapping" {
  type        = map(string)
  description = "Mapping of RDS subnet CIDR to AZ"
  default = {
    "10.0.16.0/25"   = "us-east-1a"
    "10.0.16.128/25" = "us-east-1d"
    "10.0.17.0/25"   = "us-east-1e"
    "10.0.17.128/25" = "us-east-1b"
    "10.0.18.0/25"   = "us-east-1c"
  }
}

variable "allowed_cidrs" {
  type        = list(string)
  description = "CIDRs allowed to access non-HTTPS service ports"
  default     = []
}

variable "admin_cidr" {
  type        = string
  description = "Admin CIDR for SSH and management access"
  default     = "183.80.39.60/32"
}


