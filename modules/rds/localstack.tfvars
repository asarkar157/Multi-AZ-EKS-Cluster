# Test variables for RDS module with LocalStack
# Note: These values assume a VPC and subnets have been created
# In a real scenario, you would reference outputs from the VPC module

identifier         = "test-db-instance"
vpc_id             = "vpc-test12345"
subnet_ids         = ["subnet-test1", "subnet-test2"]
availability_zones = ["us-east-1a", "us-east-1b"]
engine             = "postgres"
engine_version     = "15.4"
instance_class     = "db.t3.micro"
allocated_storage  = 20
database_name      = "testdb"
master_username    = "testadmin"

# Optional variables with defaults
backup_retention_period = 7
multi_az                = false  # Set to false for LocalStack testing
storage_encrypted        = true

tags = {
  Environment = "test"
  ManagedBy   = "terraform"
  TestRun     = "localstack"
}

