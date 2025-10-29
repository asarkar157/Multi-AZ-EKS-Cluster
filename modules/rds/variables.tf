variable "identifier" {
  description = "Identifier for the RDS instance"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where RDS will be created"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs for RDS (should span multiple AZs)"
  type        = list(string)

  validation {
    condition     = length(var.subnet_ids) >= 2
    error_message = "At least 2 subnet IDs must be provided for high availability in multi-AZ deployments."
  }
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
}

variable "engine" {
  description = "Database engine (postgres, mysql, etc.)"
  type        = string
  default     = "postgres"
}

variable "engine_version" {
  description = "Database engine version"
  type        = string
}

variable "instance_class" {
  description = "RDS instance class"
  type        = string
}

variable "allocated_storage" {
  description = "Allocated storage in GB"
  type        = number
}

variable "database_name" {
  description = "Name of the database to create"
  type        = string
}

variable "master_username" {
  description = "Master username for the database"
  type        = string
}

variable "backup_retention_period" {
  description = "Backup retention period in days"
  type        = number
  default     = 7
}

variable "multi_az" {
  description = "Enable multi-AZ deployment"
  type        = bool
  default     = true
}

variable "storage_encrypted" {
  description = "Enable storage encryption"
  type        = bool
  default     = true
}

variable "replicate_source_db" {
  description = "ARN of source database for read replica"
  type        = string
  default     = null
}

variable "allowed_security_group_ids" {
  description = "List of security group IDs allowed to access RDS"
  type        = list(string)
  default     = []
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
