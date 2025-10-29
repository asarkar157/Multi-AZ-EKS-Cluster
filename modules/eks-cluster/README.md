# EKS Cluster Module

A production-ready Terraform module for deploying AWS EKS clusters with comprehensive security, monitoring, and organizational unit (OU) based access control.

## Features

- **EKS Cluster** with configurable Kubernetes version
- **KMS Encryption** for cluster secrets
- **OIDC Provider** for IAM Roles for Service Accounts (IRSA)
- **Managed Add-ons**: VPC CNI, CoreDNS, kube-proxy, EBS CSI driver
- **CloudWatch Logging** for all control plane components
- **Multi-OU Support** with EKS access entries and policies
- **Security Groups** with least-privilege rules
- **Custom Add-on Versions** support

## Usage

### Basic Example

```hcl
module "eks_cluster" {
  source = "github.com/your-org/multi-az-eks-cluster//modules/eks-cluster?ref=eks-cluster-v1.0.0"

  cluster_name             = "my-eks-cluster"
  kubernetes_version       = "1.28"
  vpc_id                   = "vpc-12345678"
  subnet_ids               = ["subnet-1", "subnet-2", "subnet-3"]
  control_plane_subnet_ids = ["subnet-1", "subnet-2", "subnet-3"]
  environment              = "production"

  organizational_units = [
    {
      name        = "platform-ops"
      ou_id       = "ou-ops-001"
      permissions = ["admin"]
    }
  ]

  tags = {
    Environment = "production"
    ManagedBy   = "terraform"
  }
}
```

### Advanced Example with Custom Add-on Versions

```hcl
module "eks_cluster" {
  source = "github.com/your-org/multi-az-eks-cluster//modules/eks-cluster?ref=eks-cluster-v1.0.0"

  cluster_name             = "my-eks-cluster"
  kubernetes_version       = "1.28"
  vpc_id                   = "vpc-12345678"
  subnet_ids               = ["subnet-1", "subnet-2", "subnet-3"]
  control_plane_subnet_ids = ["subnet-1", "subnet-2", "subnet-3"]
  environment              = "production"

  # Custom add-on versions
  vpc_cni_version         = "v1.15.0-eksbuild.1"
  coredns_version         = "v1.10.1-eksbuild.2"
  kube_proxy_version      = "v1.28.1-eksbuild.1"
  ebs_csi_driver_version  = "v1.25.0-eksbuild.1"

  # Multiple organizational units
  organizational_units = [
    {
      name        = "platform-ops"
      ou_id       = "ou-ops-001"
      permissions = ["admin", "deploy", "view"]
    },
    {
      name        = "engineering"
      ou_id       = "ou-eng-001"
      permissions = ["deploy", "view"]
    },
    {
      name        = "readonly"
      ou_id       = "ou-ro-001"
      permissions = ["view"]
    }
  ]

  tags = {
    Environment = "production"
    Team        = "platform"
    CostCenter  = "engineering"
  }
}
```

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0 |
| aws | ~> 5.0 |
| tls | ~> 4.0 |

## Providers

| Name | Version |
|------|---------|
| aws | ~> 5.0 |
| tls | ~> 4.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| cluster_name | Name of the EKS cluster | `string` | n/a | yes |
| kubernetes_version | Kubernetes version | `string` | n/a | yes |
| vpc_id | VPC ID where EKS cluster will be created | `string` | n/a | yes |
| subnet_ids | List of subnet IDs for the EKS cluster | `list(string)` | n/a | yes |
| control_plane_subnet_ids | List of subnet IDs for the EKS control plane | `list(string)` | n/a | yes |
| organizational_units | List of organizational units for access control | `list(object)` | n/a | yes |
| environment | Environment name | `string` | n/a | yes |
| vpc_cni_version | Version of VPC CNI addon | `string` | `null` | no |
| coredns_version | Version of CoreDNS addon | `string` | `null` | no |
| kube_proxy_version | Version of kube-proxy addon | `string` | `null` | no |
| ebs_csi_driver_version | Version of EBS CSI driver addon | `string` | `null` | no |
| tags | Tags to apply to resources | `map(string)` | `{}` | no |

### Organizational Units Object

```hcl
organizational_units = [
  {
    name        = string        # OU name
    ou_id       = string        # AWS Organizations OU ID
    permissions = list(string)  # ["admin", "deploy", "view"]
  }
]
```

**Permission Levels:**
- `admin` - Full cluster access (AmazonEKSClusterAdminPolicy)
- `deploy` - Edit/deploy permissions (AmazonEKSEditPolicy)
- `view` - Read-only access (AmazonEKSViewPolicy)

## Outputs

| Name | Description |
|------|-------------|
| cluster_id | EKS cluster ID |
| cluster_name | EKS cluster name |
| cluster_endpoint | Endpoint for EKS control plane |
| cluster_security_group_id | Security group ID attached to the EKS cluster |
| cluster_primary_security_group_id | The cluster primary security group ID created by EKS |
| cluster_certificate_authority_data | Base64 encoded certificate data |
| cluster_version | The Kubernetes server version |
| oidc_provider_arn | ARN of the OIDC Provider |
| oidc_provider_url | URL of the OIDC Provider |
| cluster_iam_role_arn | IAM role ARN of the EKS cluster |
| ou_access_roles | Map of OU IDs to IAM role ARNs |

## Resources Created

This module creates the following resources:

### Core EKS Resources
- **aws_eks_cluster** - EKS cluster
- **aws_iam_role** - Cluster IAM role
- **aws_iam_role_policy_attachment** - Policy attachments (2)
- **aws_kms_key** - KMS key for encryption
- **aws_kms_alias** - KMS key alias
- **aws_security_group** - Cluster security group
- **aws_security_group_rule** - Security group rules
- **aws_cloudwatch_log_group** - CloudWatch logs

### OIDC & IRSA
- **aws_iam_openid_connect_provider** - OIDC provider

### EKS Add-ons
- **aws_eks_addon** - VPC CNI
- **aws_eks_addon** - CoreDNS
- **aws_eks_addon** - kube-proxy
- **aws_eks_addon** - EBS CSI driver

### OU-Based Access
- **aws_eks_access_entry** - Per OU (count = number of OUs)
- **aws_iam_role** - Per OU (count = number of OUs)
- **aws_eks_access_policy_association** - Per OU (count = number of OUs)

**Total Resources:** ~16 + (3 × number of OUs)

## Features in Detail

### Security

- **KMS Encryption**: All cluster secrets encrypted at rest using customer-managed KMS key
- **Security Groups**: Least-privilege security group rules
- **IAM Roles**: Separate IAM roles for cluster and organizational units
- **OIDC Provider**: Enables IAM Roles for Service Accounts (IRSA)

### Monitoring & Logging

- **CloudWatch Logs**: All control plane logs enabled:
  - API server
  - Audit
  - Authenticator
  - Controller manager
  - Scheduler
- **7-day retention** for cost optimization

### Add-ons

All managed add-ons are installed and configured:
- **VPC CNI**: Network plugin for pod networking
- **CoreDNS**: DNS server for service discovery
- **kube-proxy**: Network proxy
- **EBS CSI Driver**: Persistent volume support

Add-ons use `OVERWRITE` conflict resolution for automatic updates.

### Access Control

- **EKS Access Entries**: Modern RBAC integration
- **OU-Based Policies**: Different permissions per organizational unit
- **Cluster-Scoped**: Permissions apply to entire cluster
- **Automatic Role Creation**: IAM roles created per OU

## Post-Deployment

### Configure kubectl

```bash
aws eks update-kubeconfig --region <region> --name <cluster-name>
```

### Verify Cluster

```bash
kubectl cluster-info
kubectl get nodes
kubectl get pods -A
```

### Verify Add-ons

```bash
aws eks list-addons --cluster-name <cluster-name>
```

## Examples

See the [examples](../../examples) directory for complete examples:
- Basic EKS cluster
- Multi-OU configuration
- Custom add-on versions

## Upgrade Guide

### Kubernetes Version Upgrades

1. Update `kubernetes_version` variable
2. Update add-on versions if needed
3. Run `terraform plan` to review changes
4. Run `terraform apply`
5. Update node groups separately

### Add-on Upgrades

Add-ons can be upgraded by specifying new versions:

```hcl
vpc_cni_version = "v1.16.0-eksbuild.1"
```

## Troubleshooting

### Cluster Not Accessible

```bash
# Check cluster status
aws eks describe-cluster --name <cluster-name>

# Update kubeconfig
aws eks update-kubeconfig --name <cluster-name>

# Check IAM permissions
aws sts get-caller-identity
```

### Add-on Issues

```bash
# Check add-on status
aws eks describe-addon --cluster-name <cluster-name> --addon-name <addon-name>

# View add-on issues
aws eks describe-addon --cluster-name <cluster-name> --addon-name <addon-name> --query 'addon.health'
```

### OIDC Provider Issues

```bash
# Verify OIDC provider
aws iam list-open-id-connect-providers

# Check thumbprint
aws eks describe-cluster --name <cluster-name> --query 'cluster.identity.oidc.issuer'
```

## Security Considerations

- **Private Endpoints**: Consider using `endpoint_private_access = true` and `endpoint_public_access = false` for production
- **CIDR Blocks**: Restrict `public_access_cidrs` to known IPs
- **KMS Key Rotation**: Enabled by default
- **IAM Policies**: Follow least-privilege principle
- **Secrets Encryption**: Enabled by default with KMS

## Best Practices

1. **Use specific Kubernetes versions** instead of defaults
2. **Enable all control plane logs** for audit trail
3. **Use separate subnets** for control plane if possible
4. **Implement OU-based access** for multi-team environments
5. **Pin add-on versions** for reproducible deployments
6. **Use KMS encryption** for sensitive workloads
7. **Enable VPC Flow Logs** on associated VPC

## Cost Optimization

- **Control Plane**: ~$0.10/hour (~$73/month) - Fixed cost
- **CloudWatch Logs**: ~$0.50/GB ingested
- **KMS**: $1/month per key
- **Data Transfer**: Varies by usage

**Tip**: Use 7-day log retention to reduce costs.

## Known Limitations

- **Add-on Updates**: May require cluster restarts
- **Version Compatibility**: Ensure add-on versions match Kubernetes version
- **OU Requirements**: Requires AWS Organizations setup
- **OIDC Setup**: One-time per cluster, cannot be modified

## Contributing

Contributions welcome! Please submit issues and pull requests.

## License

MIT License - See [LICENSE](../../LICENSE) for details.

## Changelog

See [CHANGELOG](./CHANGELOG.md) for version history.

## Support

- GitHub Issues: [Report bugs](https://github.com/your-org/multi-az-eks-cluster/issues)
- Documentation: [Main README](../../README.md)

## Version

**Current Version**: 1.0.0

**Terraform Registry**: `github.com/your-org/multi-az-eks-cluster//modules/eks-cluster?ref=eks-cluster-v1.0.0`
