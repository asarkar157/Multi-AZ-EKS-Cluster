# Changelog - EKS Node Groups Module

All notable changes to the EKS Node Groups module will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2025-10-29

### Added
- Validation rule for `subnet_ids` variable to ensure at least 2 subnets are provided for high availability

### Changed
- Improved variable validation to catch configuration errors early

## [1.0.0] - 2025-10-29

### Added

#### Core Features
- EKS managed node groups with multi-AZ distribution
- Support for multiple node groups per cluster
- Configurable autoscaling per node group (min/max/desired)
- Mixed capacity types: ON_DEMAND and SPOT
- Multiple instance type support per node group

#### Launch Templates
- Custom launch templates for each node group
- Encrypted EBS volumes (gp3)
  - Configurable disk size
  - 3000 IOPS
  - 125 MB/s throughput
  - Encryption enabled by default
- IMDSv2 enforcement
  - HTTP endpoint enabled
  - HTTP tokens required
  - Hop limit: 2
  - Instance metadata tags enabled
- Detailed CloudWatch monitoring
- Network configuration
  - No public IP addresses
  - Security group attachment
  - Delete on termination
- Custom user data for node bootstrap
- Node labels (nodegroup=<name>)
- Tag specifications for instances and volumes

#### IAM Resources
- Shared IAM role for all node groups
- AWS managed policy attachments:
  - **AmazonEKSWorkerNodePolicy** - Basic node permissions
  - **AmazonEKS_CNI_Policy** - VPC networking
  - **AmazonEC2ContainerRegistryReadOnly** - ECR access
  - **AmazonSSMManagedInstanceCore** - Systems Manager
  - **AmazonEBSCSIDriverPolicy** - EBS volume management

#### Security Groups
- Dedicated security group for nodes
- Node-to-node communication rules (all traffic)
- Cluster-to-node communication rules (TCP 1025-65535)
- Node-to-cluster communication rules (HTTPS/443)
- Least-privilege access configuration

#### Node Group Configuration
- Independent configuration per node group
- Capacity type selection (ON_DEMAND/SPOT)
- Instance type flexibility (multiple types per group)
- Disk size customization
- Node labels and taints support
- Update configuration (max unavailable: 33%)
- Lifecycle management
  - Create before destroy
  - Ignore changes to desired_size (for autoscaler)

#### Multi-AZ Support
- Automatic distribution across provided subnets
- Support for 3 availability zones
- Even distribution of nodes across AZs

### Configuration

#### Required Inputs
- `cluster_name` - EKS cluster name
- `cluster_version` - Kubernetes version
- `vpc_id` - VPC ID for node groups
- `subnet_ids` - List of subnet IDs (multi-AZ)
- `cluster_security_group_id` - Cluster security group
- `cluster_primary_security_group_id` - EKS primary security group
- `node_groups` - Map of node group configurations

#### Optional Inputs
- `tags` - Resource tags (default: {})

#### Node Group Object Schema
```hcl
{
  desired_size   = number        # Desired node count
  min_size       = number        # Minimum nodes
  max_size       = number        # Maximum nodes
  instance_types = list(string)  # EC2 instance types
  capacity_type  = string        # "ON_DEMAND" or "SPOT"
  disk_size      = number        # EBS volume size (GB)
}
```

#### Outputs
- `node_groups` - Map of node group attributes
- `node_security_group_id` - Security group ID
- `node_iam_role_arn` - IAM role ARN
- `node_iam_role_name` - IAM role name

### Resources Created

#### Base Resources (10)
1. IAM role for nodes
2-6. IAM role policy attachments (5 policies)
7. Security group for nodes
8. Security group rule - node-to-node
9. Security group rule - cluster-to-node
10. Security group rule - node-to-cluster

#### Per Node Group (2 each)
- Launch template
- EKS node group

**Total:** 10 + (2 × number of node groups)

### Dependencies

#### Required Providers
- `hashicorp/aws` ~> 5.0

#### Terraform Version
- Terraform >= 1.0

#### Module Dependencies
- Requires EKS cluster to exist
- Requires VPC with subnets
- Requires cluster security groups

### Features in Detail

#### SPOT Instance Support
- Cost savings up to 90%
- Multiple instance types for availability
- 2-minute interruption warning
- Best for fault-tolerant workloads

#### Launch Template Customization
- gp3 volumes with custom IOPS
- IMDSv2 for enhanced security
- Detailed monitoring for observability
- Custom bootstrap scripts via user data

#### Security
- Encrypted EBS volumes
- IMDSv2 enforcement
- Least-privilege security groups
- Systems Manager access for debugging

#### High Availability
- Multi-AZ node distribution
- Independent node groups
- Rolling updates (33% max unavailable)
- Create before destroy lifecycle

### Best Practices

#### Sizing Recommendations
- Use node counts divisible by 3 for even AZ distribution
- Start with smaller instances, scale as needed
- Use SPOT for non-critical workloads

#### Example Configurations
```hcl
# Production: 6 nodes (2 per AZ)
desired_size = 6
min_size     = 3
max_size     = 15

# Development: 3 nodes (1 per AZ)
desired_size = 3
min_size     = 1
max_size     = 9
```

### Breaking Changes

None - this is the initial release.

### Deprecations

None - this is the initial release.

### Security Notes

- All EBS volumes encrypted by default
- IMDSv2 required for all instances
- Security groups follow least-privilege principle
- IAM policies limited to required permissions only

### Known Issues

None at this time.

### Upgrade Path

None - this is the initial release.

### Contributors

Initial release by the platform team.

---

## Future Roadmap

### Planned Features
- Taints and tolerations support
- Custom AMI support
- Node group-specific IAM roles
- GPU instance support
- ARM-based instance support
- Windows node support

### Potential Enhancements
- Karpenter integration
- Custom kubelet configurations
- Node problem detector integration
- Bottlerocket OS support
- Advanced scheduling configurations

---

## Semantic Versioning

This module follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality (backwards compatible)
- **PATCH** version for backwards compatible bug fixes

## Support

For issues, questions, or contributions:
- GitHub Issues: [Report bugs or request features](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/issues)
- Pull Requests: [Contribute improvements](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/pulls)

[1.0.0]: https://github.com/asarkar157/Multi-AZ-EKS-Cluster/releases/tag/eks-node-groups-v1.0.0
