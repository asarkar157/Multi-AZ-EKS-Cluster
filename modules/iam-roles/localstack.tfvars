# Test variables for IAM Roles module with LocalStack
# Note: These values assume an EKS cluster with OIDC provider has been created
# In a real scenario, you would reference outputs from the EKS cluster module

cluster_name = "test-eks-cluster"

# OIDC provider values (these would come from EKS cluster outputs in real usage)
oidc_provider_arn = "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/TEST123456"
oidc_provider_url = "oidc.eks.us-east-1.amazonaws.com/id/TEST123456"

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

# Optional: RDS instance ARN if RDS access is needed
# rds_instance_arn = "arn:aws:rds:us-east-1:123456789012:db:test-db-instance"

tags = {
  Environment = "test"
  ManagedBy   = "terraform"
  TestRun     = "localstack"
}

