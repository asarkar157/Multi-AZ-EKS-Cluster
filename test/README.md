# Terraform Module Tests

This directory contains unit and integration tests for the Multi-Region EKS Cluster Terraform modules using [Terratest](https://terratest.gruntwork.io/).

## Test Structure

```
test/
├── go.mod                              # Go module dependencies
├── vpc_test.go                         # VPC module unit tests
├── eks_cluster_test.go                 # EKS cluster module unit tests
├── eks_node_groups_test.go             # EKS node groups module unit tests
├── rds_test.go                         # RDS module unit tests
├── iam_roles_test.go                   # IAM roles module unit tests
├── regional_eks_integration_test.go    # Regional EKS integration tests
├── main_integration_test.go            # Full multi-region integration tests
└── README.md                           # This file
```

## Prerequisites

### Required Software

1. **Go** (version 1.21 or later)
   ```bash
   # macOS
   brew install go

   # Linux
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```

2. **Terraform** (version 1.0 or later)
   ```bash
   # macOS
   brew install terraform

   # Linux
   wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
   unzip terraform_1.6.0_linux_amd64.zip
   sudo mv terraform /usr/local/bin/
   ```

3. **AWS CLI** (configured with credentials)
   ```bash
   aws configure
   ```

### Environment Setup

Set up environment variables for testing:

```bash
export AWS_REGION=us-east-1
export AWS_PROFILE=default  # or your specific profile
export TF_VAR_primary_vpc_id=vpc-xxxxx
export TF_VAR_secondary_vpc_id=vpc-yyyyy
```

## Running Tests

### Install Dependencies

```bash
cd test
go mod download
```

### Run All Tests

```bash
# Run all tests
go test -v -timeout 30m

# Run tests in parallel
go test -v -timeout 30m -parallel 10
```

### Run Specific Test Files

```bash
# Run only VPC tests
go test -v -timeout 10m -run TestVPC

# Run only EKS cluster tests
go test -v -timeout 10m -run TestEKSCluster

# Run only integration tests
go test -v -timeout 30m -run Integration
```

### Run Individual Tests

```bash
# Run a specific test function
go test -v -timeout 10m -run TestVPCModule

# Run a specific test with verbose output
go test -v -timeout 10m -run TestEKSClusterMultipleOUs
```

### Run Tests with Coverage

```bash
go test -v -cover -timeout 30m
```

## Test Categories

### Unit Tests

Unit tests validate individual modules in isolation using `PlanOnly: true` to avoid actually creating resources:

- **vpc_test.go**: Tests VPC module configuration, subnet calculations, and validation
- **eks_cluster_test.go**: Tests EKS cluster setup, encryption, addons, and OU access
- **eks_node_groups_test.go**: Tests node group configurations, launch templates, and autoscaling
- **rds_test.go**: Tests RDS instances, read replicas, encryption, and backup configurations
- **iam_roles_test.go**: Tests IRSA roles, service account permissions, and OU-based access

### Integration Tests

Integration tests validate how modules work together:

- **regional_eks_integration_test.go**: Tests complete regional EKS setup with all components
- **main_integration_test.go**: Tests full multi-region deployment with VPC peering and RDS replication

## Test Scenarios

### VPC Module Tests

- ✅ VPC creation with 3 AZs
- ✅ Subnet CIDR calculation
- ✅ Validation of exactly 3 AZs
- ✅ Custom tags
- ✅ NAT gateway redundancy

### EKS Cluster Module Tests

- ✅ Cluster creation with encryption
- ✅ Multiple organizational units
- ✅ EKS addons (VPC CNI, CoreDNS, kube-proxy, EBS CSI)
- ✅ CloudWatch logging
- ✅ OIDC provider setup
- ✅ OU-based access control

### EKS Node Groups Module Tests

- ✅ Multiple node groups
- ✅ ON_DEMAND and SPOT instances
- ✅ Custom launch templates
- ✅ Security group configuration
- ✅ IAM roles and policies
- ✅ Multi-AZ distribution

### RDS Module Tests

- ✅ Multi-AZ RDS instances
- ✅ Read replicas
- ✅ Storage encryption with KMS
- ✅ Backup retention
- ✅ PostgreSQL and MySQL engines
- ✅ Security group rules
- ✅ Performance Insights
- ✅ Secrets Manager integration

### IAM Roles Module Tests

- ✅ IRSA roles for service accounts
- ✅ Multiple organizational units
- ✅ RDS access roles
- ✅ ALB Ingress Controller role
- ✅ EBS CSI Driver role
- ✅ Cluster Autoscaler role
- ✅ External DNS role

### Integration Tests

- ✅ Complete regional EKS setup
- ✅ Multi-region deployment
- ✅ VPC peering
- ✅ RDS replication
- ✅ Production-grade configuration
- ✅ Multiple node groups and OUs

## Cost Optimization

All tests use `PlanOnly: true` by default, which means:

- ✅ No actual AWS resources are created
- ✅ No costs incurred during testing
- ✅ Fast test execution (seconds instead of minutes/hours)
- ✅ Safe to run in CI/CD pipelines

### Full Deployment Tests (Optional)

To run full deployment tests that actually create resources:

```bash
# WARNING: This will create real AWS resources and incur costs
export TERRATEST_FULL_DEPLOYMENT=true
go test -v -timeout 120m -run TestFullDeployment
```

**Note**: Full deployment tests are disabled by default and not included in the current test suite.

## CI/CD Integration

### GitHub Actions

```yaml
name: Terraform Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Set up Terraform
        uses: hashicorp/setup-terraform@v2
        with:
          terraform_version: '1.6.0'

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Run tests
        run: |
          cd test
          go mod download
          go test -v -timeout 30m -parallel 5
```

### GitLab CI

```yaml
test:
  image: golang:1.21
  before_script:
    - apt-get update && apt-get install -y wget unzip
    - wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
    - unzip terraform_1.6.0_linux_amd64.zip -d /usr/local/bin/
  script:
    - cd test
    - go mod download
    - go test -v -timeout 30m
  only:
    - merge_requests
    - main
```

## Troubleshooting

### Common Issues

1. **Timeout Errors**
   ```bash
   # Increase timeout
   go test -v -timeout 60m
   ```

2. **AWS Credentials Not Found**
   ```bash
   # Check AWS configuration
   aws sts get-caller-identity

   # Set credentials explicitly
   export AWS_ACCESS_KEY_ID=your-key
   export AWS_SECRET_ACCESS_KEY=your-secret
   ```

3. **Terraform Init Fails**
   ```bash
   # Clean up and reinitialize
   rm -rf ../.terraform
   terraform -chdir=.. init
   ```

4. **Module Not Found**
   ```bash
   # Ensure you're in the test directory
   cd test
   go mod tidy
   ```

### Debug Mode

Enable debug logging:

```bash
export TF_LOG=DEBUG
export TERRATEST_LOG_LEVEL=debug
go test -v -timeout 30m -run TestVPCModule
```

## Best Practices

1. **Run Tests Locally First**: Always run tests locally before pushing to CI/CD
2. **Use Plan-Only Tests**: Keep costs down by using plan-only tests for validation
3. **Parallel Execution**: Use `-parallel` flag to speed up test execution
4. **Specific Tests**: Run specific tests during development to save time
5. **Cleanup**: Tests with `defer terraform.Destroy()` will clean up resources automatically
6. **Version Control**: Keep `go.mod` and `go.sum` in version control

## Writing New Tests

Example test structure:

```go
func TestNewFeature(t *testing.T) {
    t.Parallel()  // Enable parallel execution

    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir: "../modules/your-module",
        Vars: map[string]interface{}{
            "param1": "value1",
        },
        PlanOnly: true,  // Don't create actual resources
    })

    defer terraform.Destroy(t, terraformOptions)

    planStruct := terraform.InitAndPlan(t, terraformOptions)
    resourceCounts := terraform.GetResourceCount(t, planStruct)

    assert.Greater(t, resourceCounts.Add, 0, "Should create resources")
}
```

## Resources

- [Terratest Documentation](https://terratest.gruntwork.io/)
- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Assertions](https://pkg.go.dev/github.com/stretchr/testify/assert)
- [Terraform Testing Best Practices](https://www.terraform.io/docs/language/modules/testing-experiment.html)

## Contributing

When adding new features or modules:

1. Write tests first (TDD approach)
2. Ensure tests pass locally
3. Add test documentation
4. Update this README if needed
5. Submit PR with tests included

## License

Same as the main project - MIT License
