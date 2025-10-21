output "primary_cluster_endpoint" {
  description = "Endpoint for primary EKS cluster"
  value       = module.primary_region.cluster_endpoint
}

output "primary_cluster_name" {
  description = "Name of the primary EKS cluster"
  value       = module.primary_region.cluster_name
}

output "primary_cluster_security_group_id" {
  description = "Security group ID attached to the primary EKS cluster"
  value       = module.primary_region.cluster_security_group_id
}

output "secondary_cluster_endpoint" {
  description = "Endpoint for secondary EKS cluster"
  value       = module.secondary_region.cluster_endpoint
}

output "secondary_cluster_name" {
  description = "Name of the secondary EKS cluster"
  value       = module.secondary_region.cluster_name
}

output "secondary_cluster_security_group_id" {
  description = "Security group ID attached to the secondary EKS cluster"
  value       = module.secondary_region.cluster_security_group_id
}

output "primary_rds_endpoint" {
  description = "Connection endpoint for primary RDS instance"
  value       = module.primary_region.rds_endpoint
  sensitive   = true
}

output "secondary_rds_endpoint" {
  description = "Connection endpoint for secondary RDS instance"
  value       = module.secondary_region.rds_endpoint
  sensitive   = true
}

output "primary_vpc_id" {
  description = "ID of the primary region VPC"
  value       = module.primary_region.vpc_id
}

output "secondary_vpc_id" {
  description = "ID of the secondary region VPC"
  value       = module.secondary_region.vpc_id
}

output "vpc_peering_connection_id" {
  description = "ID of the VPC peering connection between regions"
  value       = aws_vpc_peering_connection.primary_to_secondary.id
}

output "primary_cluster_oidc_issuer" {
  description = "OIDC issuer URL for primary cluster (for IAM roles for service accounts)"
  value       = module.primary_region.cluster_oidc_issuer_url
}

output "secondary_cluster_oidc_issuer" {
  description = "OIDC issuer URL for secondary cluster (for IAM roles for service accounts)"
  value       = module.secondary_region.cluster_oidc_issuer_url
}
