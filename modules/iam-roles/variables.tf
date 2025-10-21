variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
}

variable "oidc_provider_arn" {
  description = "ARN of the OIDC provider for the EKS cluster"
  type        = string
}

variable "oidc_provider_url" {
  description = "URL of the OIDC provider for the EKS cluster"
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

variable "rds_instance_arn" {
  description = "ARN of the RDS instance"
  type        = string
  default     = null
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
