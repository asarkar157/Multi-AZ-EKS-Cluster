# Changelog - IAM Roles Module

All notable changes to the IAM Roles module will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.1] - 2025-10-29

### Added
- Validation rule for `cluster_name` variable to ensure valid length (1-100 characters)

### Changed
- Improved variable validation to catch configuration errors early

## [0.0.0] - Initial Development

### Added

#### Core Features
- IAM Roles for Service Accounts (IRSA) support
- Automatic role creation for common EKS add-ons
- Organizational unit-based RDS access roles
- OIDC provider integration for secure pod authentication

#### Service Account Roles
- **AWS Load Balancer Controller** - Manages ALBs/NLBs for Ingress
- **EBS CSI Driver** - Persistent volume management
- **Cluster Autoscaler** - Automatic node scaling
- **External DNS** - Automatic DNS record management

#### RDS Access Roles
- Per-organizational-unit IAM roles
- RDS access permissions
- Secrets Manager read access
- Assume role from IRSA principals

#### Security
- Least-privilege IAM policies
- OIDC-based authentication
- Trust policies for service accounts
- Namespace-scoped access

### Resources Created

#### Base Service Account Roles (4)
1. AWS Load Balancer Controller role + policy
2. EBS CSI Driver role + policy
3. Cluster Autoscaler role + policy
4. External DNS role + policy

#### Per Organizational Unit
- IAM role for RDS access
- IAM policy for RDS/Secrets Manager

**Total:** 8 base resources + (2 × number of OUs with RDS access)

### Configuration

#### Required Inputs
- `cluster_name` - EKS cluster name (1-100 characters)
- `oidc_provider_arn` - OIDC provider ARN
- `oidc_provider_url` - OIDC provider URL
- `organizational_units` - OU configuration list

#### Optional Inputs
- `rds_instance_arn` - RDS instance ARN (for RDS access roles)
- `tags` - Resource tags

#### Outputs
- `alb_controller_role_arn` - ALB Controller role ARN
- `ebs_csi_driver_role_arn` - EBS CSI Driver role ARN
- `cluster_autoscaler_role_arn` - Cluster Autoscaler role ARN
- `external_dns_role_arn` - External DNS role ARN
- `ou_rds_access_roles` - Map of OU IDs to RDS access role ARNs

### Features in Detail

#### IRSA (IAM Roles for Service Accounts)
- Pod-level IAM permissions
- No node-level credentials required
- Automatic credential rotation
- Audit trail in CloudTrail

#### Service Account Roles

**AWS Load Balancer Controller:**
- Create/manage ALB and NLB
- Manage target groups
- Configure listeners and rules
- Integration with Ingress resources

**EBS CSI Driver:**
- Create/delete EBS volumes
- Attach/detach volumes to nodes
- Snapshot management
- Volume encryption

**Cluster Autoscaler:**
- Describe Auto Scaling Groups
- Modify desired capacity
- Terminate instances
- Scale based on pod demands

**External DNS:**
- Manage Route53 records
- Automatic DNS updates
- Support for multiple hosted zones
- TTL configuration

#### RDS Access Roles
- Per-OU access control
- RDS connection permissions
- Secrets Manager integration
- Cross-service authentication

### Best Practices

#### Service Account Setup
- Deploy corresponding Kubernetes ServiceAccount
- Annotate with IAM role ARN
- Use proper namespaces
- Follow principle of least privilege

#### Trust Policies
- Scope to specific namespaces
- Limit to specific service accounts
- Review trust relationships regularly

#### Security
- Enable CloudTrail logging
- Monitor IAM role usage
- Rotate credentials regularly
- Audit policy permissions

### Usage Examples

#### AWS Load Balancer Controller
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: aws-load-balancer-controller
  namespace: kube-system
  annotations:
    eks.amazonaws.com/role-arn: <alb_controller_role_arn>
```

#### EBS CSI Driver
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ebs-csi-controller-sa
  namespace: kube-system
  annotations:
    eks.amazonaws.com/role-arn: <ebs_csi_driver_role_arn>
```

### Known Limitations
- OIDC provider must be configured first
- Service accounts must be created separately
- Role names must be unique per account
- Trust policy changes require pod restart

### Dependencies

#### Required Providers
- `hashicorp/aws` ~> 5.0

#### Terraform Version
- Terraform >= 1.0

#### Module Dependencies
- EKS cluster with OIDC provider
- Kubernetes service accounts (created separately)

---

## Support

For issues, questions, or contributions:
- GitHub Issues: [Report bugs or request features](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/issues)
- Pull Requests: [Contribute improvements](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/pulls)

[0.0.1]: https://github.com/asarkar157/Multi-AZ-EKS-Cluster/releases/tag/iam-roles-v0.0.1
