# Test variables for Regional EKS module with LocalStack
# This module combines VPC, EKS cluster, node groups, and optionally RDS

region               = "us-east-1"
cluster_name         = "test-regional-eks"
vpc_id               = "vpc-test12345"
availability_zones   = ["us-east-1a", "us-east-1b", "us-east-1c"]
environment          = "test"
kubernetes_version   = "1.28"

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

# RDS configuration (optional - set create_rds = true to enable)
create_rds = false
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

