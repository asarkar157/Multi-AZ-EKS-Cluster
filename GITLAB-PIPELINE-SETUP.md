# GitLab CI/CD Pipeline Setup for StackGen Module Publisher

This document explains how to set up and use the GitLab CI/CD pipeline for publishing Terraform modules to StackGen.

## Overview

The `.gitlab-ci.yml` pipeline replicates the functionality of the GitHub Actions `module-publisher.yml` workflow. It automatically publishes Terraform modules to StackGen when you create a Git tag following the naming convention.

## Prerequisites

1. A GitLab repository (GitLab.com or self-hosted)
2. Docker available in your GitLab CI/CD environment (uses `docker:dind` service)
3. A StackGen account with API access

## Configuration Steps

### 1. Set Up CI/CD Variables

Navigate to your GitLab project: **Settings > CI/CD > Variables**

Add the following variables:

#### Required Variables:

- **STACKGEN_TOKEN** (Type: Variable, Protected: Yes, Masked: Yes)
  - Your StackGen authentication token
  - Generate from: https://cloud.stackgen.com/project/personal/account-settings/pat

- **STACKGEN_URL** (Type: Variable)
  - Base URL of your StackGen instance
  - Example: `https://acme.cloud.stackgen.com`

#### Optional Variables:

- **STACKGEN_PROVIDER** (Type: Variable)
  - Cloud provider: `aws`, `gcp`, or `azure`
  - If not set, the pipeline will auto-detect from your Terraform files

- **STACKGEN_CLI_VERSION** (Type: Variable)
  - Specific version of the StackGen CLI Docker image
  - Default: `latest`
  - Example: `v1.2.3` or `rc/v1.2.3-rc1`

- **STACKGEN_DEBUG** (Type: Variable)
  - Set to `true` to enable debug logging
  - Default: not set (debug disabled)

- **SCM_TYPE** (Type: Variable)
  - SCM type: `gitlab`, `github`, `ado`, or `bitbucket`
  - Default: auto-detected from repository URL

### 2. Tag Naming Convention

The pipeline triggers on Git tags following these formats:

- `<module-name>-v<semantic-version>` (e.g., `eks-cluster-v1.0.0`)
- `<module-name>/v<semantic-version>` (e.g., `regional-eks/v0.1.0`)

The module must be located at: `modules/<module-name>/`

### 3. Trigger Configuration

By default, the pipeline triggers automatically on any tag. To enable automatic triggering only for properly formatted tags, uncomment the `workflow:rules` section in `.gitlab-ci.yml`:

```yaml
workflow:
  rules:
    - if: $CI_COMMIT_TAG =~ /^.+[-\/]v[0-9]+\.[0-9]+\.[0-9]+$/
```

To require manual triggering, add `when: manual` to the job:

```yaml
publish-module:
  stage: publish
  when: manual  # Add this line
  # ... rest of configuration
```

## Usage

### Publishing a Module

1. Ensure your module is in the `modules/<module-name>/` directory
2. Create a Git tag with the appropriate format:

```bash
git tag eks-cluster-v1.0.0
git push origin eks-cluster-v1.0.0
```

Or using the slash format:

```bash
git tag regional-eks/v0.1.0
git push origin regional-eks/v0.1.0
```

3. The GitLab CI/CD pipeline will automatically trigger (or require manual trigger if configured)
4. Monitor the pipeline progress in **CI/CD > Pipelines**

### What the Pipeline Does

1. **Extract Module Data**: Parses the tag name to extract module name and version
2. **List Module Contents**: Shows the files in the module directory
3. **Prepare Upload**: Configures provider, SCM type, and validates settings
4. **Upload Module**: Uses the StackGen CLI Docker image to upload the module

## Cloud Provider Detection

If `STACKGEN_PROVIDER` is not set, the pipeline automatically detects the cloud provider by analyzing your Terraform files:

- **AWS**: Searches for `hashicorp/aws`, `resource "aws_*"`, `data "aws_*"`
- **GCP**: Searches for `hashicorp/google`, `resource "google_*"`, `data "google_*"`
- **Azure**: Searches for `hashicorp/azurerm`, `resource "azurerm_*"`, `data "azurerm_*"`

The detection uses a scoring system based on the number of files containing provider-specific syntax.

## Troubleshooting

### Pipeline Fails with "STACKGEN_TOKEN is not set"

- Verify the variable is set in **Settings > CI/CD > Variables**
- Ensure the variable name is exactly `STACKGEN_TOKEN`
- Check that the variable is available to the pipeline (not restricted to specific branches/tags)

### Pipeline Fails with "Module path does not exist"

- Verify your module is located at `modules/<module-name>/`
- Check that the tag name matches the module directory name
- Example: Tag `eks-cluster-v1.0.0` expects module at `modules/eks-cluster/`

### Docker Image Pull Fails

- Ensure your GitLab runner has internet access to pull from `ghcr.io`
- Check if you need to configure Docker registry authentication
- Verify the `STACKGEN_CLI_VERSION` (if set) is a valid version

### Module Upload Fails

- Check the pipeline logs for specific error messages
- Enable debug mode by setting `STACKGEN_DEBUG=true`
- Verify `STACKGEN_URL` is correct and accessible
- Ensure your StackGen token has the required permissions

## Differences from GitHub Actions

The GitLab pipeline has the following key differences:

1. **Docker Service**: Uses `docker:dind` service instead of native Docker support
2. **Environment Variables**: Uses GitLab CI/CD variables instead of GitHub secrets/vars
3. **SCM Token**: Uses `CI_JOB_TOKEN` instead of `GITHUB_TOKEN`
4. **Trigger**: Uses tag patterns instead of GitHub release events
5. **Paths**: Uses `$CI_PROJECT_DIR` instead of `${{ github.workspace }}`

## Example Module Structure

```
Multi-AZ-EKS-Cluster/
├── .gitlab-ci.yml
├── modules/
│   ├── eks-cluster/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── README.md
│   ├── vpc/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── eks-node-groups/
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
└── ...
```

To publish the `eks-cluster` module version 1.0.0:

```bash
git tag eks-cluster-v1.0.0
git push origin eks-cluster-v1.0.0
```

## Support

For issues with:
- **StackGen CLI**: Contact StackGen support or check documentation
- **GitLab CI/CD**: Check GitLab CI/CD documentation
- **This Pipeline**: Review pipeline logs and troubleshooting section above
