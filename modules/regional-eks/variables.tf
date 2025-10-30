variable "region" {
  description = "AWS region for EKS cluster deployment"
  type        = string

  validation {
    condition     = can(regex("^[a-z]{2}-[a-z]+-[0-9]{1}$", var.region))
    error_message = "Region must be a valid AWS region format (e.g., us-east-1, eu-west-2)."
  }
}

variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
}

variable "vpc_id" {
  description = "ID of existing VPC"
  type        = string
}

variable "availability_zones" {
  description = "List of availability zones (must be exactly 3)"
  type        = list(string)
  validation {
    condition     = length(var.availability_zones) == 3
    error_message = "Exactly 3 availability zones must be specified."
  }
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "organizational_units" {
  description = "List of organizational units for access control"
  type = list(object({
    name        = string
    ou_id       = string
    permissions = list(string)
  }))
}

variable "kubernetes_version" {
  description = "Kubernetes version for EKS cluster"
  type        = string
}

variable "node_groups" {
  description = "Configuration for EKS node groups"
  type = map(object({
    desired_size   = number
    min_size       = number
    max_size       = number
    instance_types = list(string)
    capacity_type  = string
    disk_size      = number
  }))
}

variable "create_rds" {
  description = "Whether to create RDS instance"
  type        = bool
  default     = false
}

variable "rds_config" {
  description = "RDS instance configuration"
  type = object({
    engine                  = string
    engine_version          = string
    instance_class          = string
    allocated_storage       = number
    database_name           = string
    master_username         = string
    backup_retention_period = number
    multi_az                = bool
    storage_encrypted       = bool
  })
  default = {
    engine                  = "postgres"
    engine_version          = "15.4"
    instance_class          = "db.r6g.xlarge"
    allocated_storage       = 100
    database_name           = "appdb"
    master_username         = "dbadmin"
    backup_retention_period = 7
    multi_az                = true
    storage_encrypted       = true
  }
}

variable "rds_primary_arn" {
  description = "ARN of primary RDS instance (for read replicas)"
  type        = string
  default     = null
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
