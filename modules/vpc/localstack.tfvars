# Test variables for VPC module with LocalStack
region = "us-east-1"
vpc_cidr = "10.0.0.0/16"
availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]
cluster_name = "test-eks-cluster"
environment = "test"

tags = {
  Environment = "test"
  ManagedBy   = "terraform"
  TestRun     = "localstack"
}

