output "vpc_id" {
  description = "ID of the VPC"
  value       = var.vpc_id
}

output "private_subnet_ids" {
  description = "IDs of private subnets"
  value       = data.aws_subnets.private.ids
}

output "cluster_name" {
  description = "Name of the EKS cluster"
  value       = module.eks.cluster_name
}

output "cluster_endpoint" {
  description = "Endpoint for EKS cluster"
  value       = module.eks.cluster_endpoint
}

output "cluster_security_group_id" {
  description = "Security group ID of the EKS cluster"
  value       = module.eks.cluster_security_group_id
}

output "cluster_oidc_issuer_url" {
  description = "OIDC issuer URL for the cluster"
  value       = module.eks.oidc_provider_url
}

output "cluster_certificate_authority_data" {
  description = "Certificate authority data for the cluster"
  value       = module.eks.cluster_certificate_authority_data
  sensitive   = true
}

output "rds_endpoint" {
  description = "RDS instance endpoint"
  value       = var.create_rds ? module.rds[0].endpoint : null
}

output "rds_instance_arn" {
  description = "ARN of RDS instance"
  value       = var.create_rds ? module.rds[0].instance_arn : null
}

output "node_groups" {
  description = "Information about created node groups"
  value       = module.node_groups.node_groups
}
