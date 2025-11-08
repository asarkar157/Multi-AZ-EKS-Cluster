# Test variables for root module with LocalStack
# This module creates multi-region EKS clusters

primary_region   = "us-east-1"
secondary_region = "us-west-2"

cluster_name_prefix = "test-multi-region-eks"

# VPC IDs (these would be created or referenced from existing VPCs)
primary_vpc_id   = "vpc-primary-test"
secondary_vpc_id = "vpc-secondary-test"

primary_availability_zones   = ["us-east-1a", "us-east-1b", "us-east-1c"]
secondary_availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

environment = "test"

organizational_units = [
  {
    name        = "test-ops"
    ou_id       = "ou-test-ops-001"
    permissions = ["admin", "deploy", "view"]
  },
  {
    name        = "test-dev"
    ou_id       = "ou-test-dev-001"
    permissions = ["deploy", "view"]
  }
]

kubernetes_version = "1.28"

node_groups = {
  general = {
    desired_size   = 3
    min_size       = 1
    max_size       = 6
    instance_types = ["t3.medium"]
    capacity_type  = "ON_DEMAND"
    disk_size      = 20
  }
  spot = {
    desired_size   = 2
    min_size       = 0
    max_size       = 4
    instance_types = ["t3.medium"]
    capacity_type  = "SPOT"
    disk_size      = 20
  }
}

rds_config = {
  engine                  = "postgres"
  engine_version          = "15.4"
  instance_class          = "db.t3.micro"
  allocated_storage       = 20
  database_name           = "testdb"
  master_username         = "testadmin"
  backup_retention_period = 7
  multi_az                = false
  storage_encrypted        = true
}

tags = {
  Environment = "test"
  ManagedBy   = "terraform"
  TestRun     = "localstack"
}

