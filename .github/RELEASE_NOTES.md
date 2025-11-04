# Multi-Region EKS Cluster v1.0.0

We're excited to announce the initial release of the Multi-Region EKS Cluster Terraform module! 🎉

## 🚀 Overview

This release provides a production-ready, highly available multi-region EKS cluster infrastructure with comprehensive security, monitoring, and access control features.

## ✨ Key Features

### Multi-Region & Multi-AZ Architecture
- **Multi-Region Deployment**: Deploy EKS clusters across primary and secondary AWS regions
- **3 AZs per Region**: High availability with resources distributed across exactly 3 availability zones
- **VPC Peering**: Automated VPC peering between regions for seamless connectivity
- **Cross-Region RDS Replication**: Primary RDS in one region, read replica in another

### Production-Grade EKS Clusters
- **Latest Kubernetes**: Support for Kubernetes 1.28
- **Managed Add-ons**: VPC CNI, CoreDNS, kube-proxy, EBS CSI driver
- **KMS Encryption**: All cluster secrets encrypted at rest
- **OIDC Provider**: Built-in support for IAM Roles for Service Accounts (IRSA)
- **CloudWatch Logging**: Complete control plane logging

### Flexible Node Groups
- **Multi-AZ Distribution**: Nodes spread across all 3 availability zones
- **Mixed Capacity Types**: Support for both ON_DEMAND and SPOT instances
- **Custom Launch Templates**: Encrypted volumes, IMDSv2, detailed monitoring
- **Auto-scaling**: Configurable min/max/desired sizes per node group

### Secure RDS Deployment
- **Multi-AZ RDS**: High availability database deployment
- **Cross-Region Replicas**: Read replicas in secondary region
- **KMS Encryption**: Database encryption at rest
- **Secrets Manager**: Automatic credential management
- **Performance Insights**: Built-in performance monitoring
- **Automated Backups**: Configurable retention periods

### Advanced Access Control
- **Organizational Units**: Support for multiple production OUs
- **RBAC Integration**: EKS access entries for OU-based permissions
- **IRSA Roles**: Pre-configured roles for common Kubernetes services:
  - AWS Load Balancer Controller
  - EBS CSI Driver
  - Cluster Autoscaler
  - External DNS
  - OU-specific RDS access

### Comprehensive Testing
- **47 Unit Tests**: Full test coverage using Terratest
- **Integration Tests**: End-to-end multi-region testing
- **CI/CD Pipeline**: Automated testing with GitHub Actions
- **Security Scanning**: TFSec and Checkov integration
- **Zero-Cost Testing**: All tests use plan-only mode

## 📦 What's Included

### Terraform Modules
```
modules/
├── eks-cluster/         # EKS control plane
├── eks-node-groups/     # Worker nodes
├── rds/                 # Multi-AZ database
├── iam-roles/           # IRSA roles
├── regional-eks/        # Regional orchestration
└── vpc/                 # Network infrastructure (reference)
```

### Testing Suite
```
test/
├── vpc_test.go                         # VPC module tests (5 tests)
├── eks_cluster_test.go                 # EKS cluster tests (6 tests)
├── eks_node_groups_test.go             # Node groups tests (7 tests)
├── rds_test.go                         # RDS module tests (9 tests)
├── iam_roles_test.go                   # IAM roles tests (8 tests)
├── regional_eks_integration_test.go    # Regional integration (6 tests)
└── main_integration_test.go            # Multi-region integration (6 tests)
```

### Documentation
- 📖 Comprehensive README with architecture diagrams
- 📋 Testing guide (TESTING.md)
- 🔧 Example configurations
- 🛠️ Makefile with convenient commands
- 📝 Changelog
- 🔒 Security best practices

## 🎯 Quick Start

### Basic Usage

```hcl
module "multi_region_eks" {
  source = "github.com/your-org/multi-az-eks-cluster?ref=v1.0.0"

  primary_region   = "us-east-1"
  secondary_region = "us-west-2"

  primary_vpc_id   = "vpc-xxxxx"
  secondary_vpc_id = "vpc-yyyyy"

  primary_availability_zones   = ["us-east-1a", "us-east-1b", "us-east-1c"]
  secondary_availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

  cluster_name_prefix = "my-company"
  environment         = "production"

  organizational_units = [
    {
      name        = "production-ops"
      ou_id       = "ou-prod-ops-001"
      permissions = ["admin", "deploy", "view"]
    }
  ]
}
```

### Installation

```bash
# Clone the repository
git clone https://github.com/your-org/multi-az-eks-cluster.git
cd multi-az-eks-cluster

# Copy and update configuration
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values

# Initialize and apply
terraform init
terraform plan
terraform apply
```

## 📊 Resource Estimates

### Per Region
- **EKS Cluster**: ~$73/month
- **NAT Gateways**: ~$97/month each × 3 = ~$291/month
- **EKS Nodes**: Varies by instance type
- **RDS Instance**: Varies by instance class

### Total Estimated Cost
- **Minimum**: ~$1,000-2,000/month for complete multi-region setup
- Actual costs depend on instance types, storage, and data transfer

## 🔒 Security Features

- ✅ KMS encryption at rest (EKS + RDS)
- ✅ IMDSv2 enforcement
- ✅ Secrets Manager integration
- ✅ Security group least-privilege rules
- ✅ VPC Flow Logs
- ✅ IAM Roles for Service Accounts
- ✅ Deletion protection for critical resources

## ✅ Tested and Validated

All modules have been thoroughly tested:

- ✅ 47 automated tests (100% pass rate)
- ✅ Terraform validation passed
- ✅ Security scanning passed (TFSec + Checkov)
- ✅ Code formatting validated
- ✅ Integration tests passed

## 📋 Requirements

- **Terraform**: >= 1.0
- **AWS Provider**: ~> 5.0
- **Go**: >= 1.21 (for running tests)
- **Existing VPCs** with subnets tagged:
  - `Type=private` for EKS nodes
  - `Type=database` for RDS

## 🛠️ Breaking Changes

None - this is the initial release.

## 🐛 Known Issues

None at this time.

## 📚 Documentation

- [Main README](../README.md)
- [Testing Guide](../TESTING.md)
- [Changelog](../CHANGELOG.md)
- [Test Results](../TEST_RESULTS.md)

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines and submit pull requests.

## 📄 License

MIT License - See [LICENSE](../LICENSE) file for details.

## 🙏 Acknowledgments

Built with:
- [Terraform](https://www.terraform.io/)
- [Terratest](https://terratest.gruntwork.io/)
- [AWS EKS](https://aws.amazon.com/eks/)
- [AWS RDS](https://aws.amazon.com/rds/)

## 📞 Support

- 🐛 [Report Issues](https://github.com/your-org/multi-az-eks-cluster/issues)
- 💬 [Discussions](https://github.com/your-org/multi-az-eks-cluster/discussions)
- 📖 [Documentation](https://github.com/your-org/multi-az-eks-cluster)

---

**Full Changelog**: Initial Release v1.0.0

## 🎉 What's Next?

Check out our [roadmap](https://github.com/your-org/multi-az-eks-cluster/issues) for upcoming features!
