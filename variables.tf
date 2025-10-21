variable "primary_region" {
  description = "Primary AWS region"
  type        = string
  default     = "us-east-1"
}

variable "secondary_region" {
  description = "Secondary AWS region"
  type        = string
  default     = "us-west-2"
}

variable "cluster_name_prefix" {
  description = "Prefix for EKS cluster names"
  type        = string
  default     = "multi-region-eks"
}

variable "primary_vpc_id" {
  description = "ID of existing VPC in primary region"
  type        = string
}

variable "secondary_vpc_id" {
  description = "ID of existing VPC in secondary region"
  type        = string
}

variable "primary_availability_zones" {
  description = "Availability zones for primary region (must specify exactly 3)"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b", "us-east-1c"]
}

variable "secondary_availability_zones" {
  description = "Availability zones for secondary region (must specify exactly 3)"
  type        = list(string)
  default     = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

variable "environment" {
  description = "Environment name (e.g., production, staging)"
  type        = string
  default     = "production"
}

variable "organizational_units" {
  description = "List of organizational units for production access control"
  type = list(object({
    name        = string
    ou_id       = string
    permissions = list(string)
  }))
  default = [
    {
      name        = "production-ops"
      ou_id       = "ou-prod-ops-001"
      permissions = ["admin", "deploy", "view"]
    },
    {
      name        = "production-dev"
      ou_id       = "ou-prod-dev-001"
      permissions = ["deploy", "view"]
    },
    {
      name        = "production-readonly"
      ou_id       = "ou-prod-ro-001"
      permissions = ["view"]
    }
  ]
}

variable "kubernetes_version" {
  description = "Kubernetes version for EKS clusters"
  type        = string
  default     = "1.28"
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
  default = {
    general = {
      desired_size   = 6  # 2 per AZ
      min_size       = 3  # 1 per AZ
      max_size       = 15 # 5 per AZ
      instance_types = ["t3.large", "t3a.large"]
      capacity_type  = "ON_DEMAND"
      disk_size      = 50
    }
    spot = {
      desired_size   = 3 # 1 per AZ
      min_size       = 0
      max_size       = 12 # 4 per AZ
      instance_types = ["t3.large", "t3a.large", "t3.xlarge"]
      capacity_type  = "SPOT"
      disk_size      = 50
    }
  }
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

variable "tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default = {
    Terraform = "true"
    Project   = "multi-region-eks"
    ManagedBy = "terraform"
  }
}
