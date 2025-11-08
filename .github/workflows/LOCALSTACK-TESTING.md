# LocalStack Testing Guide

This document describes the LocalStack testing workflow for Terraform modules.

## Overview

The LocalStack testing workflow (`localstack-tests.yml`) automatically tests changed Terraform modules on every pull request using LocalStack, a local AWS cloud stack emulator. This allows for fast, cost-free testing without requiring actual AWS credentials.

## How It Works

1. **Detects Changed Modules**: The workflow automatically detects which Terraform modules have been modified in the PR
2. **Starts LocalStack**: Launches a LocalStack container to emulate AWS services
3. **Runs Terraform Tests**: For each changed module, it runs:
   - `terraform fmt` - Format checking
   - `tflocal init` - Initialize Terraform
   - `tflocal validate` - Validate configuration
   - `tflocal plan` - Create execution plan
   - `tflocal apply` - Apply the configuration (simulated)
   - `tflocal destroy` - Clean up resources

## Test Variables

Each module has a `localstack.tfvars` file with default test values suitable for LocalStack testing:

### Module Test Variable Files

- `modules/vpc/localstack.tfvars` - VPC module test variables
- `modules/eks-cluster/localstack.tfvars` - EKS cluster module test variables
- `modules/eks-node-groups/localstack.tfvars` - EKS node groups module test variables
- `modules/rds/localstack.tfvars` - RDS module test variables
- `modules/iam-roles/localstack.tfvars` - IAM roles module test variables
- `modules/regional-eks/localstack.tfvars` - Regional EKS module test variables
- `localstack.tfvars` - Root module test variables

### Test Variable Values

The test variables use:
- **Region**: `us-east-1` (standard AWS region)
- **CIDR Blocks**: `10.0.0.0/16` for VPCs
- **Availability Zones**: `us-east-1a`, `us-east-1b`, `us-east-1c`
- **Instance Types**: Small instances like `t3.medium` and `t3.micro` for cost efficiency
- **Storage**: Minimal values (20GB) for testing
- **Environment**: `test` for all resources

### Important Notes

1. **Dummy Resource IDs**: Some modules (like `eks-cluster` and `eks-node-groups`) reference resources that would normally come from other modules. These use placeholder IDs like `vpc-test12345` and `subnet-test1`. In a real scenario, you would use module outputs.

2. **LocalStack Limitations**: Not all AWS services are fully supported by LocalStack. Some resources may not work exactly as they would in real AWS, but the workflow will still validate Terraform syntax and basic resource creation.

3. **Variable Dependencies**: Modules that depend on outputs from other modules (e.g., `eks-cluster` needs VPC outputs) will use placeholder values in the test files. The workflow focuses on validating Terraform configuration rather than end-to-end integration.

## Running Tests Locally

You can test the LocalStack workflow locally using the provided test script:

```bash
./test-localstack.sh
```

This script will:
1. Check prerequisites (Python, Docker)
2. Install `tflocal`
3. Start LocalStack
4. Run Terraform tests on the VPC module
5. Clean up

## Workflow Triggers

The workflow runs on:
- **Pull Requests** to `main` or `develop` branches
- **Manual trigger** via `workflow_dispatch`

The workflow is skipped for:
- Changes only to `.github/workflows/**`
- Changes only to markdown files (`**/*.md`, `README*`, `CHANGELOG*`)

## Customizing Test Variables

If you need to modify test variables for a specific module:

1. Edit the corresponding `localstack.tfvars` file in the module directory
2. Ensure values are appropriate for LocalStack (avoid real AWS resource IDs)
3. Test locally using `./test-localstack.sh` or by running the workflow manually

## Troubleshooting

### LocalStack Not Starting

If LocalStack fails to start:
- Check Docker is running: `docker ps`
- Check port 4566 is available: `lsof -i :4566`
- Review LocalStack logs: `docker logs localstack-test`

### Terraform Plan/Apply Failures

Some failures are expected:
- **Missing dependencies**: Modules that depend on other modules may fail if dependencies aren't created first
- **LocalStack limitations**: Some AWS services aren't fully supported
- **Resource validation**: Some validations may fail in LocalStack but work in real AWS

The workflow uses `|| true` to allow failures in plan/apply steps, focusing on syntax and configuration validation.

### Test Variables Not Found

If a module doesn't have a `localstack.tfvars` file, the workflow will attempt to run without variables. This may fail if required variables don't have defaults, but that's acceptable for testing purposes.

## Best Practices

1. **Keep test variables minimal**: Use the smallest instance types and storage sizes
2. **Use placeholder IDs**: Don't use real AWS resource IDs in test variables
3. **Update test variables**: When adding new required variables to a module, update the corresponding `localstack.tfvars`
4. **Test locally first**: Run `./test-localstack.sh` before pushing changes

## Related Files

- `.github/workflows/localstack-tests.yml` - Main workflow file
- `test-localstack.sh` - Local testing script
- `modules/*/localstack.tfvars` - Module-specific test variables

