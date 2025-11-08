# Test configuration for terraform plan
primary_region   = "us-east-1"
secondary_region = "us-west-2"

# Mock VPC IDs for testing
primary_vpc_id   = "vpc-0123456789abcdef0"
secondary_vpc_id = "vpc-0fedcba9876543210"

# Availability Zones
primary_availability_zones   = ["us-east-1a", "us-east-1b", "us-east-1c"]
secondary_availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

# Cluster Configuration
cluster_name_prefix = "test-multi-region"
environment         = "test"
kubernetes_version  = "1.28"

# Node Groups
node_groups = {
  general = {
    desired_size   = 6
    min_size       = 3
    max_size       = 15
    instance_types = ["t3.large", "t3a.large"]
    capacity_type  = "ON_DEMAND"
    disk_size      = 50
  }
  spot = {
    desired_size   = 3
    min_size       = 0
    max_size       = 12
    instance_types = ["t3.large", "t3a.large", "t3.xlarge"]
    capacity_type  = "SPOT"
    disk_size      = 50
  }
}

# RDS Configuration
rds_config = {
  engine                  = "postgres"
  engine_version          = "15.4"
  instance_class          = "db.r6g.xlarge"
  allocated_storage       = 100
  database_name           = "testdb"
  master_username         = "dbadmin"
  backup_retention_period = 7
  multi_az                = true
  storage_encrypted       = true
}

# Organizational Units
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
  },
  {
    name        = "test-readonly"
    ou_id       = "ou-test-ro-001"
    permissions = ["view"]
  }
]

# Tags
tags = {
  Environment = "test"
  ManagedBy   = "terraform"
  Project     = "multi-region-eks-test"
}
