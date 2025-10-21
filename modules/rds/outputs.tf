output "instance_id" {
  description = "ID of the RDS instance"
  value       = var.replicate_source_db == null ? aws_db_instance.main[0].id : aws_db_instance.replica[0].id
}

output "instance_arn" {
  description = "ARN of the RDS instance"
  value       = var.replicate_source_db == null ? aws_db_instance.main[0].arn : aws_db_instance.replica[0].arn
}

output "endpoint" {
  description = "Connection endpoint"
  value       = var.replicate_source_db == null ? aws_db_instance.main[0].endpoint : aws_db_instance.replica[0].endpoint
}

output "address" {
  description = "Hostname of the RDS instance"
  value       = var.replicate_source_db == null ? aws_db_instance.main[0].address : aws_db_instance.replica[0].address
}

output "port" {
  description = "Port of the RDS instance"
  value       = var.replicate_source_db == null ? aws_db_instance.main[0].port : aws_db_instance.replica[0].port
}

output "database_name" {
  description = "Name of the database"
  value       = var.database_name
}

output "master_username" {
  description = "Master username"
  value       = var.master_username
  sensitive   = true
}

output "security_group_id" {
  description = "Security group ID for RDS"
  value       = aws_security_group.rds.id
}

output "secret_arn" {
  description = "ARN of the secret containing database credentials"
  value       = var.replicate_source_db == null ? aws_secretsmanager_secret.rds_password[0].arn : null
}
