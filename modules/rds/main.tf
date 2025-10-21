terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# Generate random password for RDS master user
resource "random_password" "master" {
  count   = var.replicate_source_db == null ? 1 : 0
  length  = 32
  special = true
}

# Store password in AWS Secrets Manager
resource "aws_secretsmanager_secret" "rds_password" {
  count                   = var.replicate_source_db == null ? 1 : 0
  name                    = "${var.identifier}-master-password"
  recovery_window_in_days = 7

  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "rds_password" {
  count     = var.replicate_source_db == null ? 1 : 0
  secret_id = aws_secretsmanager_secret.rds_password[0].id
  secret_string = jsonencode({
    username = var.master_username
    password = random_password.master[0].result
    engine   = var.engine
    host     = aws_db_instance.main[0].address
    port     = aws_db_instance.main[0].port
    dbname   = var.database_name
  })
}

# DB Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${var.identifier}-subnet-group"
  subnet_ids = var.subnet_ids

  tags = merge(
    var.tags,
    {
      Name = "${var.identifier}-subnet-group"
    }
  )
}

# Security Group for RDS
resource "aws_security_group" "rds" {
  name        = "${var.identifier}-rds-sg"
  description = "Security group for RDS instance"
  vpc_id      = var.vpc_id

  tags = merge(
    var.tags,
    {
      Name = "${var.identifier}-rds-sg"
    }
  )
}

# Allow inbound connections from EKS nodes
resource "aws_security_group_rule" "rds_ingress_eks" {
  for_each = toset(var.allowed_security_group_ids)

  description              = "Allow inbound from EKS cluster"
  type                     = "ingress"
  from_port                = local.port
  to_port                  = local.port
  protocol                 = "tcp"
  security_group_id        = aws_security_group.rds.id
  source_security_group_id = each.value
}

# Egress rule (generally not needed for RDS but included for completeness)
resource "aws_security_group_rule" "rds_egress" {
  description       = "Allow all outbound traffic"
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  security_group_id = aws_security_group.rds.id
  cidr_blocks       = ["0.0.0.0/0"]
}

# KMS Key for RDS Encryption
resource "aws_kms_key" "rds" {
  count                   = var.storage_encrypted && var.replicate_source_db == null ? 1 : 0
  description             = "KMS key for RDS instance ${var.identifier}"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = merge(
    var.tags,
    {
      Name = "${var.identifier}-rds-key"
    }
  )
}

resource "aws_kms_alias" "rds" {
  count         = var.storage_encrypted && var.replicate_source_db == null ? 1 : 0
  name          = "alias/${var.identifier}-rds"
  target_key_id = aws_kms_key.rds[0].key_id
}

# DB Parameter Group
resource "aws_db_parameter_group" "main" {
  name   = "${var.identifier}-params"
  family = local.parameter_group_family

  dynamic "parameter" {
    for_each = local.db_parameters
    content {
      name  = parameter.value.name
      value = parameter.value.value
    }
  }

  tags = merge(
    var.tags,
    {
      Name = "${var.identifier}-params"
    }
  )

  lifecycle {
    create_before_destroy = true
  }
}

# RDS Instance (Primary or Standalone)
resource "aws_db_instance" "main" {
  count = var.replicate_source_db == null ? 1 : 0

  identifier     = var.identifier
  engine         = var.engine
  engine_version = var.engine_version
  instance_class = var.instance_class

  allocated_storage = var.allocated_storage
  storage_type      = "gp3"
  storage_encrypted = var.storage_encrypted
  kms_key_id        = var.storage_encrypted ? aws_kms_key.rds[0].arn : null

  db_name  = var.database_name
  username = var.master_username
  password = random_password.master[0].result

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.main.name

  multi_az          = var.multi_az
  availability_zone = var.multi_az ? null : var.availability_zones[0]

  backup_retention_period = var.backup_retention_period
  backup_window           = "03:00-04:00"
  maintenance_window      = "mon:04:00-mon:05:00"

  enabled_cloudwatch_logs_exports = local.cloudwatch_logs_exports

  auto_minor_version_upgrade = true
  deletion_protection        = true
  skip_final_snapshot        = false
  final_snapshot_identifier  = "${var.identifier}-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"

  performance_insights_enabled          = true
  performance_insights_kms_key_id       = var.storage_encrypted ? aws_kms_key.rds[0].arn : null
  performance_insights_retention_period = 7

  monitoring_interval = 60
  monitoring_role_arn = aws_iam_role.rds_monitoring.arn

  tags = merge(
    var.tags,
    {
      Name = var.identifier
    }
  )

  lifecycle {
    ignore_changes = [
      final_snapshot_identifier
    ]
  }
}

# Read Replica (for secondary region)
resource "aws_db_instance" "replica" {
  count = var.replicate_source_db != null ? 1 : 0

  identifier          = var.identifier
  replicate_source_db = var.replicate_source_db
  instance_class      = var.instance_class

  storage_encrypted = var.storage_encrypted

  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.main.name

  multi_az          = var.multi_az
  availability_zone = var.multi_az ? null : var.availability_zones[0]

  backup_retention_period = 0 # Read replicas don't need backups
  skip_final_snapshot     = true

  auto_minor_version_upgrade = true

  performance_insights_enabled          = true
  performance_insights_retention_period = 7

  monitoring_interval = 60
  monitoring_role_arn = aws_iam_role.rds_monitoring.arn

  tags = merge(
    var.tags,
    {
      Name = "${var.identifier}-replica"
      Type = "ReadReplica"
    }
  )
}

# IAM Role for Enhanced Monitoring
resource "aws_iam_role" "rds_monitoring" {
  name = "${var.identifier}-rds-monitoring-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "monitoring.rds.amazonaws.com"
      }
    }]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "rds_monitoring" {
  role       = aws_iam_role.rds_monitoring.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}

# Locals for engine-specific configurations
locals {
  port = var.engine == "postgres" ? 5432 : var.engine == "mysql" ? 3306 : 5432

  parameter_group_family = var.engine == "postgres" ? "postgres${split(".", var.engine_version)[0]}" : var.engine == "mysql" ? "mysql${split(".", var.engine_version)[0]}" : "postgres15"

  cloudwatch_logs_exports = var.engine == "postgres" ? ["postgresql", "upgrade"] : var.engine == "mysql" ? ["error", "general", "slowquery"] : ["postgresql"]

  db_parameters = var.engine == "postgres" ? [
    {
      name  = "shared_preload_libraries"
      value = "pg_stat_statements"
    },
    {
      name  = "log_statement"
      value = "all"
    },
    {
      name  = "log_min_duration_statement"
      value = "1000"
    }
  ] : []
}
