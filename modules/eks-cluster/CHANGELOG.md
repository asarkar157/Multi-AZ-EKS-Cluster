# Changelog - EKS Cluster Module

All notable changes to the EKS Cluster module will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-21

### Added

#### Core Features
- EKS cluster creation with configurable Kubernetes version
- Support for Kubernetes 1.28 (configurable)
- VPC integration with custom subnet configuration
- Control plane and data plane subnet separation

#### Security
- KMS encryption for cluster secrets using customer-managed keys
- KMS key rotation enabled by default
- Security group with least-privilege ingress rules
- HTTPS-only access to cluster API
- IAM role for EKS cluster with required policies:
  - AmazonEKSClusterPolicy
  - AmazonEKSVPCResourceController

#### OIDC & IRSA
- OIDC provider creation for IAM Roles for Service Accounts (IRSA)
- TLS certificate fetching and fingerprint verification
- Automatic OIDC provider configuration

#### Monitoring & Logging
- CloudWatch log group for cluster logs
- All control plane log types enabled:
  - API server logs
  - Audit logs
  - Authenticator logs
  - Controller manager logs
  - Scheduler logs
- 7-day log retention for cost optimization

#### EKS Add-ons
- **VPC CNI**: Pod networking plugin
  - Configurable version support
  - OVERWRITE conflict resolution
- **CoreDNS**: DNS server for service discovery
  - Configurable version support
  - OVERWRITE conflict resolution
- **kube-proxy**: Network proxy
  - Configurable version support
  - OVERWRITE conflict resolution
- **EBS CSI Driver**: Persistent volume support
  - Configurable version support
  - OVERWRITE conflict resolution

#### Access Control
- **EKS Access Entries**: Modern EKS access management
- **Organizational Unit (OU) Support**: Multi-OU access control
- **IAM Roles per OU**: Separate roles for each organizational unit
- **Access Policy Associations**: Cluster-scoped access policies
- **Three Permission Levels**:
  - Admin: AmazonEKSClusterAdminPolicy
  - Deploy: AmazonEKSEditPolicy
  - View: AmazonEKSViewPolicy

#### Networking
- Endpoint configuration:
  - Private access enabled
  - Public access enabled (configurable)
  - Public access CIDR restriction support
- Security group rules for cluster communication

### Configuration

#### Required Inputs
- `cluster_name` - Unique cluster identifier
- `kubernetes_version` - Kubernetes version
- `vpc_id` - VPC where cluster will be created
- `subnet_ids` - Subnets for cluster networking
- `control_plane_subnet_ids` - Subnets for control plane
- `organizational_units` - OU configuration
- `environment` - Environment identifier

#### Optional Inputs
- `vpc_cni_version` - Custom VPC CNI version
- `coredns_version` - Custom CoreDNS version
- `kube_proxy_version` - Custom kube-proxy version
- `ebs_csi_driver_version` - Custom EBS CSI driver version
- `tags` - Resource tags

#### Outputs
- `cluster_id` - Cluster identifier
- `cluster_name` - Cluster name
- `cluster_endpoint` - API server endpoint
- `cluster_security_group_id` - Cluster security group ID
- `cluster_primary_security_group_id` - EKS-created security group ID
- `cluster_certificate_authority_data` - CA certificate
- `cluster_version` - Kubernetes version
- `oidc_provider_arn` - OIDC provider ARN
- `oidc_provider_url` - OIDC provider URL
- `cluster_iam_role_arn` - Cluster IAM role ARN
- `ou_access_roles` - Map of OU access roles

### Resources Created

The module creates approximately **16 base resources + 3 per organizational unit**:

#### Base Resources (16)
1. IAM role for cluster
2. IAM role policy attachment - EKS Cluster Policy
3. IAM role policy attachment - VPC Resource Controller
4. KMS key for encryption
5. KMS alias
6. Security group for cluster
7. Security group rule - HTTPS ingress
8. CloudWatch log group
9. EKS cluster
10. OIDC provider
11. EKS addon - VPC CNI
12. EKS addon - CoreDNS
13. EKS addon - kube-proxy
14. EKS addon - EBS CSI driver
15. Data source - TLS certificate
16. Data source - AWS caller identity

#### Per Organizational Unit (3 each)
- EKS access entry
- IAM role for OU access
- EKS access policy association

### Dependencies

#### Required Providers
- `hashicorp/aws` ~> 5.0
- `hashicorp/tls` ~> 4.0

#### Terraform Version
- Terraform >= 1.0

### Breaking Changes

None - this is the initial release.

### Deprecations

None - this is the initial release.

### Security Notes

- All secrets encrypted with KMS
- KMS key rotation enabled
- CloudWatch logging enabled for audit trail
- Security groups follow least-privilege principle
- OIDC provider properly configured for IRSA

### Upgrade Path

None - this is the initial release.

### Known Issues

None at this time.

### Contributors

Initial release by the platform team.

---

## Version History

### [1.0.0] - 2025-10-21
- Initial release of EKS Cluster module

## Future Roadmap

### Planned Features
- Private-only endpoint configuration option
- Custom IAM policy support for OUs
- Additional add-on support (Metrics Server, etc.)
- Fargate profile integration
- Pod identity support
- IPv6 cluster support

### Potential Enhancements
- Cluster autoscaling configuration
- Network policies
- Pod security policies/standards
- Service mesh integration guides
- Monitoring and alerting examples

---

## Semantic Versioning

This module follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality (backwards compatible)
- **PATCH** version for backwards compatible bug fixes

## Support

For issues, questions, or contributions:
- GitHub Issues: [Report bugs or request features](https://github.com/your-org/multi-az-eks-cluster/issues)
- Pull Requests: [Contribute improvements](https://github.com/your-org/multi-az-eks-cluster/pulls)

[1.0.0]: https://github.com/your-org/multi-az-eks-cluster/releases/tag/eks-cluster-v1.0.0
