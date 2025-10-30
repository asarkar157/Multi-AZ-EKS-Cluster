# Changelog - Regional EKS Module

All notable changes to the Regional EKS module will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-10-29

### Added
- Validation rule for `region` variable to ensure valid AWS region format
- Enhanced description for region variable

### Changed
- Improved variable validation to catch invalid region configurations early

## [0.0.0] - Initial Development

### Added

#### Core Features
- Regional EKS cluster orchestration
- Integration with EKS cluster module
- Integration with EKS node groups module
- Integration with RDS module
- Integration with IAM roles module
- Multi-AZ deployment support (exactly 3 AZs)

#### VPC Integration
- Data source for existing VPC lookup
- Automatic subnet discovery by tags
- Support for private, public, and database subnets
- VPC validation

#### EKS Cluster
- Kubernetes version configuration
- Multi-OU access control
- OIDC provider integration
- KMS encryption support
- CloudWatch logging

#### Node Groups
- Multiple node group support
- ON_DEMAND and SPOT capacity types
- Auto-scaling configuration
- Multi-AZ distribution
- Custom instance types

#### RDS Integration
- Optional RDS instance creation
- Multi-AZ RDS deployment
- Read replica support
- PostgreSQL and MySQL support
- Secrets Manager integration
- KMS encryption

#### IAM Roles
- Service account roles (IRSA)
- AWS Load Balancer Controller
- EBS CSI Driver
- Cluster Autoscaler
- External DNS
- Per-OU RDS access roles

### Configuration

#### Required Inputs
- `region` - AWS region (validated format)
- `cluster_name` - EKS cluster name
- `vpc_id` - Existing VPC ID
- `availability_zones` - List of 3 AZs
- `environment` - Environment name
- `organizational_units` - OU configurations
- `kubernetes_version` - Kubernetes version
- `node_groups` - Node group configurations

#### Optional Inputs
- `create_rds` - Create RDS instance (default: false)
- `rds_config` - RDS configuration object
- `rds_primary_arn` - Primary RDS ARN for replicas
- `tags` - Resource tags

#### Outputs
- `cluster_id` - EKS cluster ID
- `cluster_endpoint` - API server endpoint
- `cluster_security_group_id` - Cluster security group
- `oidc_provider_arn` - OIDC provider ARN
- `node_groups` - Node group details
- `rds_endpoint` - RDS endpoint (if created)
- `alb_controller_role_arn` - ALB controller IAM role
- `ebs_csi_driver_role_arn` - EBS CSI driver IAM role
- `cluster_autoscaler_role_arn` - Cluster autoscaler IAM role
- `external_dns_role_arn` - External DNS IAM role

### Resources Created

This module orchestrates multiple sub-modules:

#### EKS Cluster Module (~16 resources)
- EKS cluster
- KMS key
- Security group
- IAM role
- OIDC provider
- EKS add-ons
- CloudWatch log group
- Access entries (per OU)

#### Node Groups Module (10 + 2×N resources)
- IAM role + policies
- Security group + rules
- Launch templates (per node group)
- Node groups (per node group)

#### RDS Module (10-12 resources, if created)
- DB subnet group
- Security group
- KMS key
- DB parameter group
- RDS instance
- IAM monitoring role
- Secrets Manager secret

#### IAM Roles Module (8 + 2×N resources)
- Service account roles (4)
- RDS access roles (per OU)

**Total:** ~44+ base resources + additional per node group and OU

### Features in Detail

#### Subnet Discovery
Automatically discovers subnets by tags:
- **Private subnets**: `Type=private`
- **Public subnets**: `Type=public`
- **Database subnets**: `Type=database`

#### Multi-AZ Support
- Exactly 3 availability zones required
- Even distribution of nodes
- Multi-AZ RDS deployment
- High availability by default

#### Modular Design
- Loosely coupled sub-modules
- Independent lifecycle management
- Reusable components
- Clear separation of concerns

#### Security
- KMS encryption for EKS and RDS
- Security groups with least privilege
- OIDC-based pod authentication
- Secrets Manager for credentials
- IAM roles per organizational unit

### Best Practices

#### Region Selection
- Use regions with 3+ availability zones
- Consider latency to users
- Check service availability
- Plan for disaster recovery

#### Availability Zones
- Always use exactly 3 AZs
- Distribute across different zones
- Consider AZ-specific pricing
- Plan for AZ failures

#### Node Groups
- Separate system and application workloads
- Use SPOT for non-critical workloads
- Size for growth
- Monitor and adjust

#### RDS Configuration
- Enable multi-AZ for production
- Use appropriate instance class
- Configure backup retention
- Enable encryption

### Usage Example

```hcl
module "regional_eks" {
  source = "./modules/regional-eks"

  region             = "us-east-1"
  cluster_name       = "prod-eks-us-east-1"
  vpc_id             = "vpc-12345678"
  availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]
  environment        = "production"
  kubernetes_version = "1.28"

  organizational_units = [
    {
      name        = "platform-ops"
      ou_id       = "ou-ops-001"
      permissions = ["admin"]
    }
  ]

  node_groups = {
    general = {
      desired_size   = 6
      min_size       = 3
      max_size       = 15
      instance_types = ["t3.large"]
      capacity_type  = "ON_DEMAND"
      disk_size      = 50
    }
  }

  create_rds = true
  rds_config = {
    engine                  = "postgres"
    engine_version          = "15.4"
    instance_class          = "db.r6g.xlarge"
    allocated_storage       = 100
    database_name           = "appdb"
    master_username         = "dbadmin"
    backup_retention_period = 7
    multi_az                = true
    storage_encrypted       = true
  }

  tags = {
    Environment = "production"
    ManagedBy   = "terraform"
  }
}
```

### Known Limitations

- Requires existing VPC with proper subnet tags
- Exactly 3 AZs must be specified
- Subnet tags must follow naming convention
- RDS and EKS share subnet groups

### Dependencies

#### Required Providers
- `hashicorp/aws` ~> 5.0

#### Terraform Version
- Terraform >= 1.0

#### Module Dependencies
- `modules/eks-cluster`
- `modules/eks-node-groups`
- `modules/rds` (if create_rds = true)
- `modules/iam-roles`

#### External Dependencies
- Existing VPC with subnets
- Proper subnet tagging
- AWS Organizations (for OU-based access)

---

## Support

For issues, questions, or contributions:
- GitHub Issues: [Report bugs or request features](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/issues)
- Pull Requests: [Contribute improvements](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/pulls)

[0.1.0]: https://github.com/asarkar157/Multi-AZ-EKS-Cluster/releases/tag/regional-eks/v0.1.0
