# StackGen Module Publisher Workflow Documentation

## Overview

The `module-publisher.yml` workflow automatically publishes Terraform modules to StackGen when a GitHub release is published. It uses Docker to run the StackGen CLI and uploads modules to your StackGen instance.

## How It Works

### Workflow Trigger

The workflow is triggered automatically when a GitHub release is **published** (not created as a draft):

```yaml
on:
  release:
    types: [published]
```

### Workflow Steps

1. **Checkout Repository** - Checks out the full repository history to access tags
2. **Extract Module Data** - Parses the release tag to extract module information
3. **List Module Contents** - Displays the files in the module directory
4. **Prepare Upload Configuration** - Validates configuration and detects provider/SCM type
5. **Upload Custom Module** - Uses Docker to run StackGen CLI and upload the module

## Tag Naming Convention

The workflow expects release tags to follow one of these formats:

- `<module-name>-v<semantic-version>` (e.g., `rds-v0.1.7`)
- `<module-name>/v<semantic-version>` (e.g., `regional-eks/v0.1.0`)

### Valid Examples

- ✅ `rds-v0.1.7`
- ✅ `eks-cluster-v1.0.0`
- ✅ `eks-node-groups-v1.0.1`
- ✅ `regional-eks/v0.1.0`
- ✅ `vpc-v2.1.3`

### Invalid Examples

- ❌ `v1.0.0` (missing module name)
- ❌ `rds-1.0.0` (missing 'v' prefix)
- ❌ `rds-v1.0` (missing patch version)

## Module Path Convention

The workflow assumes modules are located at:
```
modules/<module-name>/
```

For example:
- Tag: `rds-v0.1.7` → Path: `modules/rds/`
- Tag: `eks-cluster-v1.0.0` → Path: `modules/eks-cluster/`

## Automatic Provider Detection

The workflow automatically detects the cloud provider by analyzing Terraform files in the module:

1. **Counts provider indicators** in `.tf` files:
   - AWS: `aws`, `hashicorp/aws`, `aws_*` resources
   - GCP: `google`, `hashicorp/google`, `google_*` resources
   - Azure: `azurerm`, `hashicorp/azurerm`, `azurerm_*` resources

2. **Selects provider** based on highest count

3. **Falls back** to AWS if no provider detected

You can override this by setting the `STACKGEN_PROVIDER` repository variable.

## SCM Type Auto-Detection

The workflow automatically detects the SCM (Source Control Management) type from the repository URL:

- **GitHub**: `github.com`, `github.`, `github.io`
- **GitLab**: `gitlab.com`, `gitlab.`, `gitlab.io`
- **Azure DevOps**: `dev.azure.com`, `azure.com`, `visualstudio.com`
- **Bitbucket**: `bitbucket.org`, `bitbucket.`, `atlassian.com`

You can override this by setting the `SCM_TYPE` repository variable.

## Required Configuration

### Repository Secrets

1. **`STACKGEN_TOKEN`** (Required)
   - Personal Access Token for StackGen
   - Generate at: `https://cloud.stackgen.com/project/personal/account-settings/pat`
   - Used to authenticate with StackGen API

### Repository Variables

1. **`STACKGEN_URL`** (Required)
   - Base URL of your StackGen instance
   - Example: `https://acme.cloud.stackgen.com`
   - Used to determine which StackGen instance to upload to

2. **`STACKGEN_PROVIDER`** (Optional)
   - Override auto-detected provider
   - Values: `aws`, `gcp`, `azure`
   - If not set, provider is auto-detected from module code

3. **`SCM_TYPE`** (Optional)
   - Override auto-detected SCM type
   - Values: `github`, `gitlab`, `ado`, `bitbucket`, `auto`
   - If not set or set to `auto`, SCM type is auto-detected from repository URL

4. **`STACKGEN_DEBUG`** (Optional)
   - Enable debug logging
   - Values: `true`, `false`
   - When `true`, adds `--log 2` flag to StackGen CLI

## How to Use

### Step 1: Configure Repository Settings

1. Go to your repository **Settings** → **Secrets and variables** → **Actions**

2. Add the required secret:
   - **Name**: `STACKGEN_TOKEN`
   - **Value**: Your StackGen Personal Access Token

3. Add the required variable:
   - **Name**: `STACKGEN_URL`
   - **Value**: Your StackGen instance URL (e.g., `https://acme.cloud.stackgen.com`)

4. (Optional) Add additional variables if needed:
   - `STACKGEN_PROVIDER` - Override provider detection
   - `SCM_TYPE` - Override SCM type detection
   - `STACKGEN_DEBUG` - Enable debug logging

### Step 2: Create a Release

#### Option A: Using GitHub CLI

```bash
# Create and push a tag
git tag -a "rds-v0.1.7" -m "Release rds/v0.1.7 - patch version update"
git push origin "rds-v0.1.7"

# Create a GitHub release
gh release create "rds-v0.1.7" \
  --title "RDS Module v0.1.7" \
  --notes "Patch version update for RDS module"
```

#### Option B: Using GitHub Web Interface

1. Go to your repository on GitHub
2. Click **Releases** → **Draft a new release**
3. **Choose a tag**: Create new tag `rds-v0.1.7` (or select existing)
4. **Release title**: `RDS Module v0.1.7`
5. **Describe this release**: Add release notes
6. Click **Publish release**

### Step 3: Monitor the Workflow

1. Go to **Actions** tab in GitHub
2. Find the **StackGen Module Publisher** workflow run
3. Click on the run to see detailed logs
4. The workflow will:
   - Parse the release tag
   - Extract module information
   - Detect provider and SCM type
   - Upload the module to StackGen

### Step 4: Verify Upload

Check the workflow logs for success message:

```
✅ Custom module uploaded successfully
  • Module: rds
  • Version:   v0.1.7
  • Path:      modules/rds
  • Provider:  aws
  • Status:    Published ✓
```

## Workflow Outputs

The workflow provides detailed information in the logs:

### Release Information

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📝 Terraform Release Details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📌 Release Tag:       rds-v0.1.7
📦 TF Module Name:    rds
🔢 Semantic Version:  0.1.7
📁 Module Path:    modules/rds
✓  Path Exists:       true
   URL:               https://github.com/.../releases/tag/rds-v0.1.7
```

### Module Contents

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📂 MODULE CONTENTS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Path: modules/rds

total 48
drwxr-xr-x  7 runner  staff   224 Nov 11 21:48 .
drwxr-xr-x  8 runner  staff   256 Nov 11 21:48 ..
-rw-r--r--  1 runner  staff  5234 Nov 11 21:48 CHANGELOG.md
-rw-r--r--  1 runner  staff 12456 Nov 11 21:48 README.md
-rw-r--r--  1 runner  staff  8765 Nov 11 21:48 main.tf
-rw-r--r--  1 runner  staff  1234 Nov 11 21:48 outputs.tf
-rw-r--r--  1 runner  staff  2345 Nov 11 21:48 variables.tf
```

### Upload Configuration

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📤 UPLOADING CUSTOM MODULE TO STACKGEN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Provider:     aws
StackGen URL: https://acme.cloud.stackgen.com
Module Dir:   modules/rds
SCM_REPO:     https://github.com/user/repo
SCM_TAG:      rds-v0.1.7
VERSION_NAME: v0.1.7
SCM_TYPE:     github
```

## Troubleshooting

### Error: Tag name does not match expected format

**Problem**: The release tag doesn't follow the required format.

**Solution**: Ensure your tag follows one of these formats:
- `<module-name>-v<semantic-version>` (e.g., `rds-v0.1.7`)
- `<module-name>/v<semantic-version>` (e.g., `regional-eks/v0.1.0`)

### Error: Module path does not exist

**Problem**: The module directory doesn't exist at the expected path.

**Solution**: 
- Verify the module exists at `modules/<module-name>/`
- Check that the module name in the tag matches the directory name
- Example: Tag `rds-v0.1.7` expects path `modules/rds/`

### Error: STACKGEN_TOKEN is not set

**Problem**: The required secret is missing.

**Solution**:
1. Go to repository **Settings** → **Secrets and variables** → **Actions**
2. Add secret named `STACKGEN_TOKEN`
3. Value should be your StackGen Personal Access Token
4. Generate token at: `https://cloud.stackgen.com/project/personal/account-settings/pat`

### Error: STACKGEN_URL is not set

**Problem**: The required variable is missing.

**Solution**:
1. Go to repository **Settings** → **Secrets and variables** → **Actions** → **Variables**
2. Add variable named `STACKGEN_URL`
3. Value should be your StackGen instance base URL
4. Example: `https://acme.cloud.stackgen.com`

### Upload Fails with Authentication Error

**Problem**: StackGen token is invalid or expired.

**Solution**:
1. Verify `STACKGEN_TOKEN` secret is correct
2. Generate a new token if needed
3. Ensure token has permissions to upload custom modules

### Wrong Provider Detected

**Problem**: Auto-detection selects wrong cloud provider.

**Solution**:
1. Set `STACKGEN_PROVIDER` repository variable
2. Values: `aws`, `gcp`, or `azure`
3. This overrides auto-detection

### Wrong SCM Type Detected

**Problem**: Auto-detection selects wrong SCM type.

**Solution**:
1. Set `SCM_TYPE` repository variable
2. Values: `github`, `gitlab`, `ado`, `bitbucket`
3. This overrides auto-detection

## Docker Image

The workflow uses the official StackGen Docker image:

```yaml
ghcr.io/stackgenhq/stackgen:latest
```

This image contains the StackGen CLI and is maintained by StackGen.

## Environment Variables Passed to Docker

The following environment variables are passed to the Docker container:

- `STACKGEN_PROVIDER` - Cloud provider (aws, gcp, azure)
- `SCM_TOKEN` - GitHub token (from `GITHUB_TOKEN`)
- `STACKGEN_TOKEN` - StackGen Personal Access Token
- `STACKGEN_URL` - StackGen instance URL
- `SCM_REPO` - Repository URL
- `SCM_TAG` - Release tag name
- `SCM_TYPE` - SCM type (github, gitlab, ado, bitbucket)

## StackGen CLI Command

The workflow runs the following StackGen CLI command inside Docker:

```bash
stackgen upload custom-modules \
  --dir modules/<module-name> \
  --provider <provider> \
  --overwrite-version \
  --version v<semantic-version> \
  [--log 2]  # If STACKGEN_DEBUG=true
```

## Version Overwriting

The workflow uses the `--overwrite-version` flag, which means:
- If a module version already exists in StackGen, it will be overwritten
- This allows you to republish the same version if needed
- Useful for fixing issues in a release without bumping the version

## Best Practices

1. **Semantic Versioning**: Always use semantic versioning (MAJOR.MINOR.PATCH)
2. **Consistent Tagging**: Use the same tag format across all modules
3. **Release Notes**: Always include meaningful release notes
4. **Test Before Release**: Test your module locally before publishing
5. **Monitor Workflows**: Check workflow runs after publishing to ensure success
6. **Version Management**: Don't reuse version numbers unless intentionally overwriting

## Example Workflow

Here's a complete example of publishing a module:

```bash
# 1. Make changes to your module
cd modules/rds
# ... make changes ...

# 2. Commit changes
git add .
git commit -m "feat(rds): Add new feature"

# 3. Create and push tag
git tag -a "rds-v0.1.7" -m "Release rds/v0.1.7 - patch version update"
git push origin "rds-v0.1.7"

# 4. Create GitHub release
gh release create "rds-v0.1.7" \
  --title "RDS Module v0.1.7" \
  --notes "Patch version update for RDS module"

# 5. Monitor workflow
gh run watch --workflow="StackGen Module Publisher"
```

## Related Workflows

- **`module-publisher-cli.yml`** - Alternative workflow using StackGen CLI directly (currently disabled)
- **`localstack-tests.yml`** - Tests modules with LocalStack before release
- **`terraform-tests.yml`** - Validates Terraform modules

## Maintenance

### Updating Docker Image

The workflow uses `ghcr.io/stackgenhq/stackgen:latest`. To pin to a specific version:

```yaml
ghcr.io/stackgenhq/stackgen:v1.2.3  # Use specific version
```

### Updating Actions

Periodically update GitHub Actions versions:

```yaml
- uses: actions/checkout@v4  # Check for newer versions
```

## Security Considerations

1. **Secrets**: Never commit `STACKGEN_TOKEN` to the repository
2. **Variables**: Repository variables are visible to all workflow runs
3. **Token Permissions**: Use tokens with minimal required permissions
4. **Docker**: The workflow runs Docker with access to workspace files

## Support

For issues or questions:
1. Check workflow logs in GitHub Actions
2. Verify repository secrets and variables are set correctly
3. Consult StackGen documentation: https://docs.stackgen.com
4. Check StackGen support channels

## License

Part of the Multi-AZ-EKS-Cluster project - MIT License

