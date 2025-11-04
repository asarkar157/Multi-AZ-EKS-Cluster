# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-21

### Added

#### Core Infrastructure
- Multi-region EKS cluster deployment (primary + secondary regions)
- Multi-AZ architecture with exactly 3 availability zones per region
- VPC peering between regions for cross-region connectivity
- Support for existing VPC infrastructure

#### EKS Cluster Module (`modules/eks-cluster`)
- EKS cluster with Kubernetes 1.28 support
- KMS encryption for cluster secrets
- OIDC provider for IAM Roles for Service Accounts (IRSA)
- CloudWatch logging for all control plane components
- EKS managed addons:
  - VPC CNI for networking
  - CoreDNS for service discovery
  - kube-proxy for networking
  - EBS CSI driver for persistent volumes
- Security group with least-privilege rules
- Multiple organizational unit (OU) support
- EKS access entries for OU-based RBAC
- Custom addon version support

#### EKS Node Groups Module (`modules/eks-node-groups`)
- Multi-AZ node group distribution
- Support for ON_DEMAND and SPOT capacity types
- Custom launch templates with:
  - Encrypted gp3 EBS volumes
  - IMDSv2 enforcement
  - Custom disk sizing
  - Detailed monitoring
- Security groups for node-to-node and node-to-control plane communication
- IAM roles with required AWS managed policies:
  - AmazonEKSWorkerNodePolicy
  - AmazonEKS_CNI_Policy
  - AmazonEC2ContainerRegistryReadOnly
  - AmazonSSMManagedInstanceCore
  - AmazonEBSCSIDriverPolicy
- Autoscaling configuration per node group
- Node labeling and taints support

#### RDS Module (`modules/rds`)
- Multi-AZ RDS deployment
- Cross-region read replica support
- KMS encryption at rest
- Random password generation
- AWS Secrets Manager integration
- DB parameter groups with engine-specific optimizations
- DB subnet groups spanning multiple AZs
- Security groups for EKS cluster access
- Performance Insights enabled
- Enhanced monitoring with CloudWatch
- Automated backups with configurable retention
- Support for PostgreSQL and MySQL engines
- Final snapshot on deletion

#### IAM Roles Module (`modules/iam-roles`)
- IRSA roles for Kubernetes service accounts
- OU-specific RDS access roles
- AWS Load Balancer Controller role and policy
- EBS CSI Driver role
- Cluster Autoscaler role and policy
- External DNS role and policy
- Proper OIDC trust relationships
- Least-privilege IAM policies

#### Regional EKS Module (`modules/regional-eks`)
- Complete regional EKS setup orchestration
- Integration of EKS cluster, node groups, and RDS
- Data sources for existing VPC subnet discovery
- Conditional RDS creation
- Read replica configuration for secondary regions

#### VPC Module (`modules/vpc`) - Reference
- Multi-AZ VPC with 3 availability zones
- Public, private, and database subnet tiers
- NAT gateway redundancy (one per AZ)
- Internet gateway for public subnets
- VPC Flow Logs to CloudWatch
- DB subnet groups for RDS
- Automatic subnet CIDR calculation

#### Testing & Quality
- 47 comprehensive unit and integration tests using Terratest
- GitHub Actions CI/CD pipeline
- Pre-commit hooks for code quality
- TFLint configuration
- TFSec security scanning
- Checkov compliance checking
- Makefile with convenient commands
- 100% test coverage

#### Documentation
- Comprehensive README with architecture diagrams
- Module-specific documentation
- Testing guide (TESTING.md)
- Example configurations (basic and advanced)
- Troubleshooting guide
- Cost estimation
- Post-deployment steps
- Best practices

#### Security Features
- KMS encryption for EKS secrets
- KMS encryption for RDS storage
- Encrypted EBS volumes with gp3
- IMDSv2 enforcement on EC2 instances
- Secrets Manager for database credentials
- Security group least-privilege rules
- VPC Flow Logs for network monitoring
- IAM Roles for Service Accounts (IRSA)
- Deletion protection for RDS

#### High Availability Features
- Multi-region deployment
- Multi-AZ within each region
- RDS Multi-AZ deployment
- Cross-region RDS read replicas
- NAT gateway redundancy
- VPC peering for cross-region communication

#### Operational Features
- CloudWatch logging for EKS control plane
- Performance Insights for RDS
- Enhanced monitoring for RDS
- VPC Flow Logs
- Automated RDS backups
- Point-in-time recovery for RDS
- Final snapshots before deletion

### Configuration
- Terraform >= 1.0 required
- AWS Provider ~> 5.0
- Kubernetes Provider ~> 2.0
- Support for custom Kubernetes versions
- Configurable node group parameters
- Configurable RDS parameters
- Tag propagation to all resources

### Examples
- Basic multi-region setup example
- Advanced configuration with multiple node groups
- Production-grade configuration example
- terraform.tfvars.example template

### CI/CD
- GitHub Actions workflow for automated testing
- Terraform validation on all modules
- Unit tests for individual modules
- Integration tests for complete setups
- Security scanning with TFSec and Checkov
- Automated test summary and notifications

## [Unreleased]

### Planned
- Support for additional AWS regions
- EKS cluster autoscaling policies
- Network policies examples
- Monitoring and alerting configurations
- Disaster recovery playbooks
- Blue-green deployment support

---

## Release Notes Template

### Version Numbering
- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality (backwards compatible)
- **PATCH** version for backwards compatible bug fixes

### Categories
- **Added** - New features
- **Changed** - Changes in existing functionality
- **Deprecated** - Soon-to-be removed features
- **Removed** - Removed features
- **Fixed** - Bug fixes
- **Security** - Security improvements

[1.0.0]: https://github.com/your-org/multi-az-eks-cluster/releases/tag/v1.0.0
