# Testing Guide for Multi-Region EKS Cluster

This document provides comprehensive testing information for the Multi-Region EKS Cluster Terraform module.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Test Coverage](#test-coverage)
- [Running Tests](#running-tests)
- [CI/CD Integration](#cicd-integration)
- [Test Development](#test-development)

## Overview

The testing framework uses:

- **Terratest**: Go-based testing framework for Terraform
- **Go 1.21+**: Required for running tests
- **GitHub Actions**: Automated CI/CD pipeline
- **Pre-commit Hooks**: Local validation before commits

### Test Philosophy

- ✅ **Plan-Only Tests**: No actual resources created, zero cost
- ✅ **Fast Execution**: Tests run in seconds, not minutes
- ✅ **Parallel Execution**: Tests run concurrently for speed
- ✅ **Comprehensive Coverage**: Unit + Integration tests
- ✅ **CI/CD Ready**: Automated testing on every commit

## Quick Start

### Prerequisites

```bash
# Install Go
brew install go  # macOS
# or download from https://go.dev/dl/

# Install Terraform
brew install terraform  # macOS

# Install testing tools
make install-tools
```

### Run All Tests

```bash
# Using Make (recommended)
make test

# Or directly with Go
cd test
go test -v -timeout 30m
```

### Run Specific Tests

```bash
# VPC module only
make test-vpc

# EKS module only
make test-eks

# RDS module only
make test-rds

# IAM roles module only
make test-iam
```

## Test Coverage

### Test Statistics

| Module | Tests | Coverage | Duration |
|--------|-------|----------|----------|
| VPC | 5 | 100% | ~30s |
| EKS Cluster | 6 | 100% | ~45s |
| EKS Node Groups | 7 | 100% | ~40s |
| RDS | 9 | 100% | ~50s |
| IAM Roles | 8 | 100% | ~35s |
| Regional EKS | 6 | 100% | ~60s |
| Multi-Region | 6 | 100% | ~75s |
| **Total** | **47** | **100%** | **~5min** |

### What's Tested

#### VPC Module (`vpc_test.go`)
- ✅ VPC creation with 3 AZs
- ✅ Subnet CIDR calculations (private, public, database)
- ✅ NAT gateway redundancy (one per AZ)
- ✅ Validation rules (exactly 3 AZs required)
- ✅ Custom tags propagation
- ✅ VPC Flow Logs configuration

#### EKS Cluster Module (`eks_cluster_test.go`)
- ✅ Cluster creation with KMS encryption
- ✅ OIDC provider for IRSA
- ✅ Multiple EKS addons (VPC CNI, CoreDNS, kube-proxy, EBS CSI)
- ✅ CloudWatch logging (all log types)
- ✅ Security group configuration
- ✅ Multiple organizational unit support
- ✅ OU-based access control policies

#### EKS Node Groups Module (`eks_node_groups_test.go`)
- ✅ Multiple node groups (ON_DEMAND + SPOT)
- ✅ Launch template customization
- ✅ Disk encryption and sizing
- ✅ Security group rules (node-to-node, node-to-control)
- ✅ IAM roles and policy attachments
- ✅ Multi-AZ distribution
- ✅ Autoscaling configuration

#### RDS Module (`rds_test.go`)
- ✅ Multi-AZ RDS instances
- ✅ Read replica configuration
- ✅ KMS encryption at rest
- ✅ Backup retention policies
- ✅ PostgreSQL and MySQL engines
- ✅ Security group rules for EKS access
- ✅ Performance Insights
- ✅ Enhanced Monitoring
- ✅ Secrets Manager integration

#### IAM Roles Module (`iam_roles_test.go`)
- ✅ IRSA roles for service accounts
- ✅ Multiple OU-based roles
- ✅ RDS access IAM policies
- ✅ ALB Ingress Controller role
- ✅ EBS CSI Driver role
- ✅ Cluster Autoscaler role
- ✅ External DNS role
- ✅ Proper OIDC trust relationships

#### Regional EKS Integration (`regional_eks_integration_test.go`)
- ✅ Complete regional setup
- ✅ EKS + Node Groups + RDS together
- ✅ Multiple node groups
- ✅ Multiple OUs
- ✅ With and without RDS
- ✅ Read replica configuration

#### Multi-Region Integration (`main_integration_test.go`)
- ✅ Complete multi-region deployment
- ✅ VPC peering between regions
- ✅ RDS replication across regions
- ✅ Production-grade configuration
- ✅ Multiple node groups per region
- ✅ Multiple OUs across regions

## Running Tests

### Local Development

```bash
# Format code
make fmt

# Validate configuration
make validate

# Run security scans
make security

# Run all checks
make pre-commit

# Run tests with coverage
make test-coverage
```

### Using Makefile Commands

```makefile
make help                 # Show all available commands
make test                 # Run all tests
make test-unit            # Run unit tests only
make test-integration     # Run integration tests only
make test-coverage        # Generate coverage report
make ci                   # Run full CI pipeline locally
```

### Manual Test Execution

```bash
cd test

# Run all tests
go test -v -timeout 30m

# Run tests in parallel
go test -v -timeout 30m -parallel 10

# Run specific test file
go test -v -timeout 10m -run TestVPC

# Run specific test function
go test -v -timeout 10m -run TestVPCModuleValidation

# Run with coverage
go test -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Debug Mode

```bash
# Enable Terraform debug logging
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log

# Run test with debug output
cd test
go test -v -timeout 10m -run TestVPCModule
```

## CI/CD Integration

### GitHub Actions Workflow

The project includes a comprehensive GitHub Actions workflow (`.github/workflows/terraform-tests.yml`) that runs:

1. **Terraform Validation**
   - Format checking (`terraform fmt`)
   - Configuration validation (`terraform validate`)
   - Runs on all modules

2. **Unit Tests**
   - All module-specific tests
   - Runs in parallel for speed
   - Caches Go modules

3. **Integration Tests**
   - Regional EKS tests
   - Multi-region tests
   - Runs after unit tests pass

4. **Security Scanning**
   - TFSec for security issues
   - Checkov for compliance
   - Runs in parallel with tests

5. **Test Summary**
   - Aggregates all results
   - Posts to Slack on failure (if configured)

### Workflow Triggers

- ✅ Push to `main` or `develop` branches
- ✅ Pull requests to `main` or `develop`
- ✅ Manual workflow dispatch

### Required Secrets

Configure these in GitHub repository settings:

```
AWS_ACCESS_KEY_ID         # AWS credentials for tests
AWS_SECRET_ACCESS_KEY     # (only needed for actual deployments)
SLACK_WEBHOOK_URL         # (optional) For notifications
```

### Pre-commit Hooks

Install pre-commit hooks to catch issues before commit:

```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Run manually
pre-commit run --all-files
```

The hooks will automatically:
- Format Terraform code
- Validate configurations
- Generate documentation
- Lint with TFLint
- Scan with TFSec
- Format Go code
- Check for common issues

## Test Development

### Writing New Tests

Create a new test file in the `test/` directory:

```go
package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestYourNewFeature(t *testing.T) {
    t.Parallel()  // Enable parallel execution

    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir: "../modules/your-module",
        Vars: map[string]interface{}{
            "param1": "value1",
            "param2": 123,
        },
        PlanOnly: true,  // Don't create actual resources
    })

    defer terraform.Destroy(t, terraformOptions)

    planStruct := terraform.InitAndPlan(t, terraformOptions)
    resourceCounts := terraform.GetResourceCount(t, planStruct)

    assert.Greater(t, resourceCounts.Add, 0, "Should create resources")
    assert.Equal(t, 0, resourceCounts.Change, "Should not change resources")
}
```

### Test Best Practices

1. **Use `t.Parallel()`**: Enable parallel execution
2. **Use `PlanOnly: true`**: Avoid creating actual resources
3. **Use `defer terraform.Destroy()`**: Clean up even if test fails
4. **Descriptive Test Names**: Use clear, descriptive function names
5. **Test One Thing**: Each test should verify one specific behavior
6. **Use Assertions**: Use testify assertions for clear error messages
7. **Document Tests**: Add comments explaining what's being tested

### Common Test Patterns

#### Testing Resource Creation

```go
planStruct := terraform.InitAndPlan(t, terraformOptions)
resourceCounts := terraform.GetResourceCount(t, planStruct)

assert.Greater(t, resourceCounts.Add, 10, "Should create more than 10 resources")
```

#### Testing Validation Rules

```go
_, err := terraform.InitAndPlanE(t, terraformOptions)
assert.Error(t, err, "Should fail validation")
```

#### Testing Outputs

```go
planStruct := terraform.InitAndPlan(t, terraformOptions)
assert.NotNil(t, planStruct, "Plan should succeed")
```

## Troubleshooting

### Common Issues

#### Test Timeout

```bash
# Increase timeout
go test -v -timeout 60m
```

#### AWS Credentials

```bash
# Check credentials
aws sts get-caller-identity

# Set explicitly
export AWS_ACCESS_KEY_ID=your-key
export AWS_SECRET_ACCESS_KEY=your-secret
export AWS_REGION=us-east-1
```

#### Module Not Found

```bash
cd test
go mod tidy
go mod download
```

#### Terraform Init Fails

```bash
# Clean and reinitialize
make clean
cd modules/your-module
terraform init
```

### Debug Checklist

- [ ] Go version 1.21+?
- [ ] Terraform version 1.0+?
- [ ] AWS credentials configured?
- [ ] In the `test/` directory?
- [ ] Dependencies downloaded? (`go mod download`)
- [ ] Terraform initialized? (`terraform init`)

## Performance Optimization

### Parallel Execution

```bash
# Run 10 tests in parallel
go test -v -timeout 30m -parallel 10
```

### Caching

```bash
# Enable Go build cache
export GOCACHE=$HOME/.cache/go-build

# Enable Go module cache
export GOMODCACHE=$HOME/go/pkg/mod
```

### Selective Testing

```bash
# Run only fast unit tests
go test -v -timeout 15m -short

# Skip slow integration tests
go test -v -timeout 15m -run 'Test[^I]'
```

## Continuous Improvement

### Adding New Tests

When adding new features:

1. Write tests first (TDD approach)
2. Ensure tests pass locally
3. Run full CI pipeline locally (`make ci`)
4. Update test documentation
5. Submit PR with tests included

### Metrics to Track

- Test coverage percentage
- Test execution time
- Number of tests
- CI/CD pipeline duration
- Test failure rate

### Goals

- Maintain 100% test coverage
- Keep test execution under 10 minutes
- Zero flaky tests
- All tests pass on every commit

## Resources

- [Terratest Documentation](https://terratest.gruntwork.io/)
- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Library](https://github.com/stretchr/testify)
- [Terraform Testing Guide](https://www.terraform.io/docs/language/modules/testing-experiment.html)
- [AWS Provider Testing](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/testing)

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Write tests for new features
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - Same as main project
