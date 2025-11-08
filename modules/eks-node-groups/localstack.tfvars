# Test variables for EKS Node Groups module with LocalStack
# Note: These values assume an EKS cluster has been created
# In a real scenario, you would reference outputs from the EKS cluster module

cluster_name   = "test-eks-cluster"
cluster_version = "1.28"
vpc_id         = "vpc-test12345"
subnet_ids     = ["subnet-test1", "subnet-test2", "subnet-test3"]
cluster_security_group_id         = "sg-test-cluster"
cluster_primary_security_group_id = "sg-test-primary"

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

tags = {
  Environment = "test"
  ManagedBy   = "terraform"
  TestRun     = "localstack"
}

