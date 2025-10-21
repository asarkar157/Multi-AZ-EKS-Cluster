variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
}

variable "cluster_version" {
  description = "Kubernetes version of the cluster"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where node groups will be created"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs for the node groups (should span 3 AZs)"
  type        = list(string)
}

variable "cluster_security_group_id" {
  description = "Security group ID of the EKS cluster"
  type        = string
}

variable "cluster_primary_security_group_id" {
  description = "The cluster primary security group ID created by EKS"
  type        = string
}

variable "node_groups" {
  description = "Map of node group configurations"
  type = map(object({
    desired_size   = number
    min_size       = number
    max_size       = number
    instance_types = list(string)
    capacity_type  = string
    disk_size      = number
  }))
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
