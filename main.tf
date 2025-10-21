terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

# Primary Region Provider
provider "aws" {
  region = var.primary_region
  alias  = "primary"
}

# Secondary Region Provider
provider "aws" {
  region = var.secondary_region
  alias  = "secondary"
}

# Primary Region Infrastructure
module "primary_region" {
  source = "./modules/regional-eks"

  providers = {
    aws = aws.primary
  }

  region               = var.primary_region
  cluster_name         = "${var.cluster_name_prefix}-primary"
  vpc_id               = var.primary_vpc_id
  availability_zones   = var.primary_availability_zones
  environment          = var.environment
  organizational_units = var.organizational_units

  # EKS Configuration
  kubernetes_version = var.kubernetes_version
  node_groups        = var.node_groups

  # RDS Configuration
  create_rds = true
  rds_config = var.rds_config

  tags = merge(
    var.tags,
    {
      Region = var.primary_region
      Type   = "Primary"
    }
  )
}

# Secondary Region Infrastructure
module "secondary_region" {
  source = "./modules/regional-eks"

  providers = {
    aws = aws.secondary
  }

  region               = var.secondary_region
  cluster_name         = "${var.cluster_name_prefix}-secondary"
  vpc_id               = var.secondary_vpc_id
  availability_zones   = var.secondary_availability_zones
  environment          = var.environment
  organizational_units = var.organizational_units

  # EKS Configuration
  kubernetes_version = var.kubernetes_version
  node_groups        = var.node_groups

  # RDS Configuration - Read replica
  create_rds      = true
  rds_config      = var.rds_config
  rds_primary_arn = module.primary_region.rds_instance_arn

  tags = merge(
    var.tags,
    {
      Region = var.secondary_region
      Type   = "Secondary"
    }
  )
}

# VPC Peering between regions for RDS access
resource "aws_vpc_peering_connection" "primary_to_secondary" {
  provider    = aws.primary
  vpc_id      = module.primary_region.vpc_id
  peer_vpc_id = module.secondary_region.vpc_id
  peer_region = var.secondary_region
  auto_accept = false

  tags = merge(
    var.tags,
    {
      Name = "primary-to-secondary-peering"
    }
  )
}

resource "aws_vpc_peering_connection_accepter" "secondary" {
  provider                  = aws.secondary
  vpc_peering_connection_id = aws_vpc_peering_connection.primary_to_secondary.id
  auto_accept               = true

  tags = merge(
    var.tags,
    {
      Name = "primary-to-secondary-peering-accepter"
    }
  )
}
