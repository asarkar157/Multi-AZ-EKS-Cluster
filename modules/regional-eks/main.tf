terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Data sources for existing VPC resources
data "aws_vpc" "existing" {
  id = var.vpc_id
}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }

  filter {
    name   = "tag:Type"
    values = ["private"]
  }
}

data "aws_subnets" "database" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }

  filter {
    name   = "tag:Type"
    values = ["database"]
  }
}

# EKS Cluster Module
module "eks" {
  source = "../eks-cluster"

  cluster_name             = var.cluster_name
  kubernetes_version       = var.kubernetes_version
  vpc_id                   = var.vpc_id
  subnet_ids               = data.aws_subnets.private.ids
  control_plane_subnet_ids = data.aws_subnets.private.ids

  organizational_units = var.organizational_units
  environment          = var.environment

  tags = var.tags
}

# EKS Node Groups Module
module "node_groups" {
  source = "../eks-node-groups"

  cluster_name    = module.eks.cluster_name
  cluster_version = var.kubernetes_version
  vpc_id          = var.vpc_id
  subnet_ids      = data.aws_subnets.private.ids
  node_groups     = var.node_groups

  cluster_security_group_id         = module.eks.cluster_security_group_id
  cluster_primary_security_group_id = module.eks.cluster_primary_security_group_id

  tags = var.tags

  depends_on = [module.eks]
}

# RDS Instance Module
module "rds" {
  count  = var.create_rds ? 1 : 0
  source = "../rds"

  identifier         = "${var.cluster_name}-db"
  vpc_id             = var.vpc_id
  subnet_ids         = data.aws_subnets.database.ids
  availability_zones = var.availability_zones

  engine                  = var.rds_config.engine
  engine_version          = var.rds_config.engine_version
  instance_class          = var.rds_config.instance_class
  allocated_storage       = var.rds_config.allocated_storage
  database_name           = var.rds_config.database_name
  master_username         = var.rds_config.master_username
  backup_retention_period = var.rds_config.backup_retention_period
  multi_az                = var.rds_config.multi_az
  storage_encrypted       = var.rds_config.storage_encrypted

  # For read replicas in secondary region
  replicate_source_db = var.rds_primary_arn

  # Allow access from EKS cluster
  allowed_security_group_ids = [
    module.eks.cluster_security_group_id,
    module.node_groups.node_security_group_id
  ]

  tags = var.tags
}

# IAM Roles for Service Accounts (IRSA) for OU-based access
module "iam_roles" {
  source = "../iam-roles"

  cluster_name         = module.eks.cluster_name
  oidc_provider_arn    = module.eks.oidc_provider_arn
  oidc_provider_url    = module.eks.oidc_provider_url
  organizational_units = var.organizational_units
  rds_instance_arn     = var.create_rds ? module.rds[0].instance_arn : null

  tags = var.tags
}
