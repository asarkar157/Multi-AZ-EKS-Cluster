# Release Monitor Workflow

This GitHub Action automatically runs whenever a new release is published in this IaC module monorepo.

## Purpose

The workflow parses release information, extracts semantic versioning details, identifies the affected submodule, and displays comprehensive information in the workflow logs.

## Trigger

```yaml
on:
  release:
    types: [published]
```

The workflow triggers when a GitHub release is **published** (not created as a draft).

## What It Does

### 1. Parse Release Tag
Extracts information from the release tag name following the convention:
```
<submodule-name>-v<semantic-version>
```

**Examples:**
- `eks-cluster-v1.0.0`
- `eks-node-groups-v1.0.1`
- `rds-v2.1.3`

### 2. Extract Components

From the tag, it extracts:
- **Submodule Name**: The module identifier (e.g., `eks-cluster`)
- **Semantic Version**: The version number (e.g., `1.0.1`)
- **Submodule Path**: The full path in the repo (e.g., `modules/eks-cluster`)

### 3. Validate Path
Checks if the submodule path exists in the repository.

### 4. Display Information

The workflow outputs detailed information in the logs:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📦 NEW RELEASE PUBLISHED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📌 Release Tag:       eks-node-groups-v1.0.1
📂 Submodule Name:    eks-node-groups
🔢 Semantic Version:  1.0.1
📁 Submodule Path:    modules/eks-node-groups
✓  Path Exists:       true

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📝 Release Details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Title:     EKS Node Groups Module v1.0.1
Author:    asarkar157
Created:   2025-10-29T10:30:00Z
Published: 2025-10-29T10:30:00Z
URL:       https://github.com/.../releases/tag/eks-node-groups-v1.0.1

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 5. Semantic Version Breakdown

Parses the semantic version into components:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔢 SEMANTIC VERSION BREAKDOWN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Full Version:  1.0.1
Major:         1
Minor:         0
Patch:         1

Release Type:  🔧 PATCH RELEASE (Bug Fixes)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Release Type Detection:**
- **🚀 MAJOR RELEASE** (e.g., v2.0.0) - Breaking changes
- **✨ MINOR RELEASE** (e.g., v1.1.0) - New features
- **🔧 PATCH RELEASE** (e.g., v1.0.1) - Bug fixes
- **🧪 PRE-RELEASE** (e.g., v0.x.x) - Development versions

### 6. List Submodule Contents

Shows the files in the released submodule:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📂 SUBMODULE CONTENTS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Path: modules/eks-node-groups

total 48
drwxr-xr-x  7 runner  staff   224 Oct 29 10:30 .
drwxr-xr-x  8 runner  staff   256 Oct 29 10:30 ..
-rw-r--r--  1 runner  staff  5234 Oct 29 10:30 CHANGELOG.md
-rw-r--r--  1 runner  staff 12456 Oct 29 10:30 README.md
-rw-r--r--  1 runner  staff  8765 Oct 29 10:30 main.tf
-rw-r--r--  1 runner  staff  1234 Oct 29 10:30 outputs.tf
-rw-r--r--  1 runner  staff  2345 Oct 29 10:30 variables.tf
```

### 7. Detect Changed Files

Compares the current release with the previous release for the same submodule:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📝 CHANGED FILES IN THIS RELEASE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Comparing: eks-node-groups-v1.0.0 → eks-node-groups-v1.0.1

Changed files in modules/eks-node-groups:

M       modules/eks-node-groups/CHANGELOG.md
M       modules/eks-node-groups/variables.tf
```

**Status Codes:**
- `M` - Modified
- `A` - Added
- `D` - Deleted
- `R` - Renamed

## Tag Naming Convention

The workflow expects release tags to follow this format:

```
<submodule-name>-v<major>.<minor>.<patch>
```

**Valid Examples:**
- ✅ `eks-cluster-v1.0.0`
- ✅ `eks-node-groups-v1.0.1`
- ✅ `rds-v2.1.3`
- ✅ `iam-roles-v1.2.0`

**Invalid Examples:**
- ❌ `v1.0.0` (missing submodule name)
- ❌ `eks-cluster-1.0.0` (missing 'v' prefix)
- ❌ `eks-cluster-v1.0` (missing patch version)

## Outputs

The workflow sets the following outputs from the `parse` step:

| Output | Description | Example |
|--------|-------------|---------|
| `submodule_name` | Name of the submodule | `eks-node-groups` |
| `semantic_version` | Version number (without 'v') | `1.0.1` |
| `submodule_path` | Path to submodule directory | `modules/eks-node-groups` |
| `path_exists` | Whether the path exists | `true` |

These outputs can be used by subsequent jobs if you extend the workflow.

## Usage

### Viewing Logs

1. Go to **Actions** tab in GitHub
2. Click on **Release Monitor** workflow
3. Select the workflow run for your release
4. Click on **parse-release** job
5. Expand the steps to see detailed output

### Extending the Workflow

You can add additional jobs that depend on the parse results:

```yaml
  notify-slack:
    needs: parse-release
    runs-on: ubuntu-latest
    steps:
      - name: Send Slack notification
        run: |
          echo "New release: ${{ needs.parse-release.outputs.submodule_name }} v${{ needs.parse-release.outputs.semantic_version }}"
          # Add Slack webhook call here

  update-documentation:
    needs: parse-release
    runs-on: ubuntu-latest
    steps:
      - name: Update docs
        run: |
          # Automatically update documentation sites
          # or trigger other downstream processes
```

## Example Output

Here's what the workflow logs look like for a real release:

```
Run: Parse release tag and extract information
Release tag: eks-node-groups-v1.0.1
✅ Submodule found at: modules/eks-node-groups

Run: Display release information
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📦 NEW RELEASE PUBLISHED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📌 Release Tag:       eks-node-groups-v1.0.1
📂 Submodule Name:    eks-node-groups
🔢 Semantic Version:  1.0.1
📁 Submodule Path:    modules/eks-node-groups
✓  Path Exists:       true

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Troubleshooting

### Tag Format Error

If you see:
```
❌ Error: Tag name does not match expected format '<submodule-name>-v<semantic-version>'
   Expected format: eks-cluster-v1.0.0
   Received: v1.0.0
```

**Solution**: Ensure your release tag includes the submodule name prefix.

### Path Not Found

If you see:
```
⚠️ Warning: Submodule path does not exist: modules/xyz
```

**Solution**:
- Check that the submodule name in the tag matches the directory name
- Verify the submodule exists at `modules/<submodule-name>`

### No Previous Tag Found

If you see:
```
This is the first release for this submodule
```

This is normal for the initial release of a submodule.

## Related Workflows

- **terraform-tests.yml** - Runs tests on PRs
- **terraform-docs.yml** - Generates documentation
- **security-scan.yml** - Scans for security issues

## Maintenance

This workflow requires minimal maintenance. Update the checkout action version periodically:

```yaml
- uses: actions/checkout@v4  # Check for newer versions
```

## Contributing

To modify this workflow:

1. Edit `.github/workflows/release-monitor.yml`
2. Test locally using [act](https://github.com/nektos/act)
3. Submit a PR with your changes

## License

Part of the Multi-AZ-EKS-Cluster project - MIT License
