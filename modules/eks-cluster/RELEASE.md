# EKS Cluster Module v1.0.0

Production-ready Terraform module for AWS EKS cluster deployment with comprehensive security, monitoring, and multi-organizational unit access control.

## 🚀 Quick Start

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

## ✨ Features

### Core Capabilities
- ✅ **EKS Cluster** with configurable Kubernetes version (1.28+)
- ✅ **KMS Encryption** for all cluster secrets at rest
- ✅ **OIDC Provider** for IAM Roles for Service Accounts (IRSA)
- ✅ **CloudWatch Logging** for all control plane components
- ✅ **Security Groups** with least-privilege access rules

### Managed Add-ons
- ✅ **VPC CNI** - Pod networking plugin
- ✅ **CoreDNS** - DNS server for service discovery
- ✅ **kube-proxy** - Network proxy for services
- ✅ **EBS CSI Driver** - Persistent volume support

All add-ons support custom version specification.

### Access Control
- ✅ **Multi-OU Support** - Multiple organizational unit access
- ✅ **EKS Access Entries** - Modern RBAC integration
- ✅ **Three Permission Levels**:
  - **Admin** - Full cluster access
  - **Deploy** - Edit and deploy permissions
  - **View** - Read-only access
- ✅ **Automatic IAM Role Creation** per organizational unit

## 📦 What's Included

### Module Files
```
modules/eks-cluster/
├── main.tf          # Main resource definitions
├── variables.tf     # Input variables
├── outputs.tf       # Module outputs
├── README.md        # Complete documentation
└── CHANGELOG.md     # Version history
```

### Resources Created

**Base Resources (~16):**
- 1 EKS Cluster
- 1 KMS Key + Alias
- 1 Security Group + Rules
- 1 IAM Role + 2 Policy Attachments
- 1 OIDC Provider
- 4 EKS Add-ons
- 1 CloudWatch Log Group

**Per Organizational Unit (+3 each):**
- 1 EKS Access Entry
- 1 IAM Role
- 1 Access Policy Association

**Total:** ~16 + (3 × number of OUs)

## 🔒 Security Features

- **KMS Encryption**: Customer-managed KMS key with automatic rotation
- **Secrets Encryption**: All Kubernetes secrets encrypted at rest
- **Security Groups**: Least-privilege HTTPS-only access
- **OIDC Provider**: Secure IAM role assumption for pods
- **CloudWatch Logs**: Complete audit trail (7-day retention)
- **IAM Policies**: Separate roles per organizational unit

## 📊 Requirements

| Component | Version |
|-----------|---------|
| Terraform | >= 1.0 |
| AWS Provider | ~> 5.0 |
| TLS Provider | ~> 4.0 |

## 🎯 Use Cases

### Single Organization
```hcl
organizational_units = [
  {
    name        = "platform-team"
    ou_id       = "ou-platform-001"
    permissions = ["admin"]
  }
]
```

### Multiple Teams
```hcl
organizational_units = [
  {
    name        = "ops"
    ou_id       = "ou-ops-001"
    permissions = ["admin", "deploy", "view"]
  },
  {
    name        = "developers"
    ou_id       = "ou-dev-001"
    permissions = ["deploy", "view"]
  },
  {
    name        = "auditors"
    ou_id       = "ou-audit-001"
    permissions = ["view"]
  }
]
```

### Custom Add-on Versions
```hcl
vpc_cni_version        = "v1.15.0-eksbuild.1"
coredns_version        = "v1.10.1-eksbuild.2"
kube_proxy_version     = "v1.28.1-eksbuild.1"
ebs_csi_driver_version = "v1.25.0-eksbuild.1"
```

## 📝 Inputs

### Required

| Name | Description | Type |
|------|-------------|------|
| cluster_name | Name of the EKS cluster | `string` |
| kubernetes_version | Kubernetes version | `string` |
| vpc_id | VPC ID | `string` |
| subnet_ids | Subnet IDs for cluster | `list(string)` |
| control_plane_subnet_ids | Subnets for control plane | `list(string)` |
| organizational_units | OU configurations | `list(object)` |
| environment | Environment name | `string` |

### Optional

| Name | Description | Type | Default |
|------|-------------|------|---------|
| vpc_cni_version | VPC CNI version | `string` | `null` (latest) |
| coredns_version | CoreDNS version | `string` | `null` (latest) |
| kube_proxy_version | kube-proxy version | `string` | `null` (latest) |
| ebs_csi_driver_version | EBS CSI version | `string` | `null` (latest) |
| tags | Resource tags | `map(string)` | `{}` |

## 📤 Outputs

| Name | Description |
|------|-------------|
| cluster_id | EKS cluster ID |
| cluster_endpoint | API server endpoint |
| cluster_security_group_id | Cluster security group ID |
| oidc_provider_arn | OIDC provider ARN (for IRSA) |
| oidc_provider_url | OIDC provider URL |
| cluster_certificate_authority_data | Cluster CA certificate |
| ou_access_roles | Map of OU IDs to IAM role ARNs |

## 🛠️ Post-Deployment

### Configure kubectl
```bash
aws eks update-kubeconfig \
  --region us-east-1 \
  --name my-eks-cluster
```

### Verify Cluster
```bash
kubectl cluster-info
kubectl get nodes
kubectl get pods -A
```

### Check Add-ons
```bash
aws eks list-addons --cluster-name my-eks-cluster
```

## 💰 Cost Estimate

- **EKS Control Plane**: ~$73/month (fixed)
- **CloudWatch Logs**: ~$0.50/GB ingested
- **KMS Key**: $1/month
- **Add-ons**: Free (included with EKS)

**Note:** Node costs are separate (handled by node groups module).

## 🔄 Upgrade Guide

### Kubernetes Version
```hcl
# Update version
kubernetes_version = "1.29"

# Plan and apply
terraform plan
terraform apply
```

### Add-on Versions
```hcl
# Specify new version
vpc_cni_version = "v1.16.0-eksbuild.1"

# Apply
terraform apply
```

## 🐛 Known Limitations

- OIDC provider cannot be modified after creation (requires cluster recreation)
- Add-on version must be compatible with Kubernetes version
- Requires AWS Organizations setup for OU-based access

## 📚 Documentation

- [README.md](README.md) - Complete module documentation
- [CHANGELOG.md](CHANGELOG.md) - Version history
- [Main Repo README](../../README.md) - Full project documentation

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📄 License

MIT License - See [LICENSE](../../LICENSE) for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/your-org/multi-az-eks-cluster/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/multi-az-eks-cluster/discussions)
- **Documentation**: [Module README](README.md)

## 🎉 What's Next?

- Node groups for worker nodes
- RDS for database
- IAM roles for service accounts
- Complete multi-region setup

See the [main repository](../../) for complete multi-region examples.

---

**Module Version**: 1.0.0
**Release Date**: October 29, 2025
**Terraform Registry**: `github.com/your-org/multi-az-eks-cluster//modules/eks-cluster?ref=eks-cluster-v1.0.0`

**Status**: ✅ Production Ready
