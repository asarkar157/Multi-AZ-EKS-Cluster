# Test variables for EKS Cluster module with LocalStack
# Note: These values assume a VPC and subnets have been created
# In a real scenario, you would reference outputs from the VPC module

cluster_name       = "test-eks-cluster"
kubernetes_version = "1.28"
vpc_id             = "vpc-test12345"
subnet_ids         = ["subnet-test1", "subnet-test2", "subnet-test3"]
control_plane_subnet_ids = ["subnet-test1", "subnet-test2", "subnet-test3"]
environment        = "test"

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

tags = {
  Environment = "test"
  ManagedBy   = "terraform"
  TestRun     = "localstack"
}

