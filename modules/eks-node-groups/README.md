# EKS Node Groups Module

A production-ready Terraform module for deploying AWS EKS node groups with multi-AZ distribution, custom launch templates, and support for both ON_DEMAND and SPOT capacity types.

## Features

- **Multi-AZ Distribution** - Nodes automatically distributed across multiple availability zones
- **Mixed Capacity Types** - Support for ON_DEMAND and SPOT instances
- **Custom Launch Templates** - Encrypted EBS volumes, IMDSv2, monitoring
- **Security Groups** - Least-privilege rules for node communication
- **IAM Roles** - Pre-configured with all required AWS managed policies
- **Autoscaling** - Configurable min/max/desired sizes per node group
- **Multiple Node Groups** - Deploy multiple node groups with different configurations

## Usage

### Basic Example

```hcl
module "eks_node_groups" {
  source = "github.com/asarkar157/Multi-AZ-EKS-Cluster//modules/eks-node-groups?ref=eks-node-groups-v1.0.0"

  cluster_name                      = "my-eks-cluster"
  cluster_version                   = "1.28"
  vpc_id                            = "vpc-12345678"
  subnet_ids                        = ["subnet-1", "subnet-2", "subnet-3"]
  cluster_security_group_id         = "sg-cluster123"
  cluster_primary_security_group_id = "sg-primary456"

  node_groups = {
    general = {
      desired_size   = 6
      min_size       = 3
      max_size       = 15
      instance_types = ["t3.large", "t3a.large"]
      capacity_type  = "ON_DEMAND"
      disk_size      = 50
    }
  }

  tags = {
    Environment = "production"
    ManagedBy   = "terraform"
  }
}
```

### Multiple Node Groups with SPOT Instances

```hcl
module "eks_node_groups" {
  source = "github.com/asarkar157/Multi-AZ-EKS-Cluster//modules/eks-node-groups?ref=eks-node-groups-v1.0.0"

  cluster_name                      = "my-eks-cluster"
  cluster_version                   = "1.28"
  vpc_id                            = "vpc-12345678"
  subnet_ids                        = ["subnet-1", "subnet-2", "subnet-3"]
  cluster_security_group_id         = "sg-cluster123"
  cluster_primary_security_group_id = "sg-primary456"

  node_groups = {
    # ON_DEMAND node group for critical workloads
    general = {
      desired_size   = 6   # 2 per AZ
      min_size       = 3   # 1 per AZ
      max_size       = 15  # 5 per AZ
      instance_types = ["t3.large", "t3a.large"]
      capacity_type  = "ON_DEMAND"
      disk_size      = 50
    }

    # SPOT node group for non-critical workloads
    spot = {
      desired_size   = 3   # 1 per AZ
      min_size       = 0
      max_size       = 12  # 4 per AZ
      instance_types = ["t3.large", "t3a.large", "t3.xlarge"]
      capacity_type  = "SPOT"
      disk_size      = 50
    }

    # Compute-optimized for specific workloads
    compute = {
      desired_size   = 3
      min_size       = 0
      max_size       = 9
      instance_types = ["c5.2xlarge", "c5a.2xlarge"]
      capacity_type  = "ON_DEMAND"
      disk_size      = 100
    }
  }

  tags = {
    Environment = "production"
    Team        = "platform"
  }
}
```

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0 |
| aws | ~> 5.0 |

## Providers

| Name | Version |
|------|---------|
| aws | ~> 5.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| cluster_name | Name of the EKS cluster | `string` | n/a | yes |
| cluster_version | Kubernetes version of the cluster | `string` | n/a | yes |
| vpc_id | VPC ID where node groups will be created | `string` | n/a | yes |
| subnet_ids | List of subnet IDs (should span 3 AZs) | `list(string)` | n/a | yes |
| cluster_security_group_id | Security group ID of the EKS cluster | `string` | n/a | yes |
| cluster_primary_security_group_id | EKS-created primary security group ID | `string` | n/a | yes |
| node_groups | Map of node group configurations | `map(object)` | n/a | yes |
| tags | Tags to apply to resources | `map(string)` | `{}` | no |

### Node Groups Configuration

```hcl
node_groups = {
  "<node-group-name>" = {
    desired_size   = number        # Desired number of nodes
    min_size       = number        # Minimum number of nodes
    max_size       = number        # Maximum number of nodes
    instance_types = list(string)  # EC2 instance types
    capacity_type  = string        # "ON_DEMAND" or "SPOT"
    disk_size      = number        # EBS volume size in GB
  }
}
```

## Outputs

| Name | Description |
|------|-------------|
| node_groups | Map of node group attributes (id, arn, status, etc.) |
| node_security_group_id | Security group ID for node groups |
| node_iam_role_arn | IAM role ARN for node groups |
| node_iam_role_name | IAM role name for node groups |

## Resources Created

This module creates the following resources:

### IAM Resources
- **aws_iam_role** - IAM role for node groups (shared across all groups)
- **aws_iam_role_policy_attachment** - Policy attachments (5):
  - AmazonEKSWorkerNodePolicy
  - AmazonEKS_CNI_Policy
  - AmazonEC2ContainerRegistryReadOnly
  - AmazonSSMManagedInstanceCore
  - AmazonEBSCSIDriverPolicy

### Security Groups
- **aws_security_group** - Node security group (shared)
- **aws_security_group_rule** - Node-to-node communication
- **aws_security_group_rule** - Node-to-cluster communication
- **aws_security_group_rule** - Cluster-to-node HTTPS

### Per Node Group
- **aws_launch_template** - Custom launch template with encryption
- **aws_eks_node_group** - EKS managed node group

**Total Resources:** 10 base resources + (2 × number of node groups)

## Features in Detail

### Multi-AZ Distribution

Nodes are automatically distributed across all provided subnets (typically 3 AZs):

```hcl
subnet_ids = [
  "subnet-az1",  # Availability Zone 1
  "subnet-az2",  # Availability Zone 2
  "subnet-az3",  # Availability Zone 3
]
```

EKS will distribute nodes evenly across these AZs for high availability.

### Launch Template Features

Each node group gets a custom launch template with:

- **Encrypted EBS Volumes**
  - Volume type: gp3
  - IOPS: 3000
  - Throughput: 125 MB/s
  - Encryption: Enabled

- **IMDSv2 Enforcement**
  - HTTP endpoint: Enabled
  - HTTP tokens: Required
  - Hop limit: 2
  - Instance metadata tags: Enabled

- **Detailed Monitoring**
  - CloudWatch detailed monitoring enabled

- **Network Configuration**
  - No public IP addresses
  - Security group attached
  - Delete on termination

- **Node Labels**
  - Automatic labeling: `nodegroup=<name>`

### Security Groups

The module creates security group rules for:

1. **Node-to-Node Communication**
   - All traffic allowed between nodes
   - Protocol: All
   - Ports: All

2. **Cluster-to-Node Communication**
   - Kubelet and pod communication
   - Protocol: TCP
   - Ports: 1025-65535

3. **Node-to-Cluster Communication**
   - API server access from pods
   - Protocol: TCP
   - Port: 443

### IAM Policies

All nodes get the following AWS managed policies:

- **AmazonEKSWorkerNodePolicy** - Basic EKS node permissions
- **AmazonEKS_CNI_Policy** - VPC CNI networking
- **AmazonEC2ContainerRegistryReadOnly** - ECR image pulls
- **AmazonSSMManagedInstanceCore** - Systems Manager access
- **AmazonEBSCSIDriverPolicy** - EBS volume management

### Capacity Types

#### ON_DEMAND
- Guaranteed capacity
- Fixed pricing
- Best for: Critical workloads, databases, stateful applications

#### SPOT
- Up to 90% cost savings
- Can be interrupted with 2-minute warning
- Best for: Batch jobs, CI/CD, fault-tolerant workloads
- Supports multiple instance types for better availability

### Autoscaling

Each node group has independent autoscaling configuration:

```hcl
desired_size = 6   # Current desired capacity
min_size     = 3   # Minimum nodes (for scaling down)
max_size     = 15  # Maximum nodes (for scaling up)
```

Use with Cluster Autoscaler for automatic scaling based on pod demands.

## Best Practices

### Multi-AZ Sizing

For 3 AZs, use node counts divisible by 3:

```hcl
desired_size = 6   # 2 nodes per AZ
min_size     = 3   # 1 node per AZ
max_size     = 15  # 5 nodes per AZ
```

### SPOT Instance Configuration

Use multiple instance types for better availability:

```hcl
spot = {
  instance_types = [
    "t3.large",
    "t3a.large",
    "t3.xlarge",
    "t2.large"
  ]
  capacity_type = "SPOT"
}
```

### Disk Sizing

- **50 GB**: Sufficient for most workloads
- **100 GB**: Data-intensive applications
- **200+ GB**: Container build nodes, data processing

### Instance Types

Choose based on workload:

- **General Purpose**: t3, t3a, m5, m5a
- **Compute Optimized**: c5, c5a, c6i
- **Memory Optimized**: r5, r5a, r6i
- **GPU**: p3, p4, g4dn

## Examples

### Dedicated Node Groups

```hcl
node_groups = {
  # System components
  system = {
    desired_size   = 3
    min_size       = 3
    max_size       = 3
    instance_types = ["t3.medium"]
    capacity_type  = "ON_DEMAND"
    disk_size      = 30
  }

  # Application workloads
  application = {
    desired_size   = 6
    min_size       = 3
    max_size       = 20
    instance_types = ["m5.large", "m5a.large"]
    capacity_type  = "ON_DEMAND"
    disk_size      = 50
  }

  # Batch processing
  batch = {
    desired_size   = 0
    min_size       = 0
    max_size       = 10
    instance_types = ["c5.2xlarge", "c5a.2xlarge"]
    capacity_type  = "SPOT"
    disk_size      = 100
  }
}
```

## Post-Deployment

### Verify Node Groups

```bash
# List node groups
aws eks list-nodegroups --cluster-name my-eks-cluster

# Describe specific node group
aws eks describe-nodegroup \
  --cluster-name my-eks-cluster \
  --nodegroup-name general

# Check nodes in cluster
kubectl get nodes -o wide

# View node labels
kubectl get nodes --show-labels
```

### Configure Cluster Autoscaler

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cluster-autoscaler
  namespace: kube-system
spec:
  template:
    spec:
      containers:
      - name: cluster-autoscaler
        image: k8s.gcr.io/autoscaling/cluster-autoscaler:v1.28.0
        command:
          - ./cluster-autoscaler
          - --v=4
          - --cloud-provider=aws
          - --skip-nodes-with-local-storage=false
          - --expander=least-waste
          - --node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/my-eks-cluster
```

## Troubleshooting

### Nodes Not Joining Cluster

```bash
# Check node group status
aws eks describe-nodegroup \
  --cluster-name my-eks-cluster \
  --nodegroup-name general \
  --query 'nodegroup.health'

# Check EC2 instances
aws ec2 describe-instances \
  --filters "Name=tag:eks:cluster-name,Values=my-eks-cluster"

# View node group events
kubectl get events -A | grep node
```

### SPOT Interruptions

```bash
# Check SPOT interruption warnings
kubectl get events -A | grep -i spot

# View node conditions
kubectl describe node <node-name>
```

### Launch Template Issues

```bash
# Verify launch template
aws ec2 describe-launch-templates \
  --launch-template-names my-eks-cluster-*

# Check launch template versions
aws ec2 describe-launch-template-versions \
  --launch-template-id lt-xxxxx
```

## Cost Optimization

### Instance Sizing
- Start with smaller instances (t3.medium/large)
- Monitor CPU/memory usage
- Scale up only when needed

### SPOT Savings
- Use SPOT for non-critical workloads
- Can save up to 90% vs ON_DEMAND
- Use multiple instance types

### Autoscaling
- Set appropriate min/max values
- Use Cluster Autoscaler to scale based on demand
- Scale to zero for dev/test environments

### Example Costs (us-east-1)

**ON_DEMAND:**
- t3.large: ~$75/month per node
- t3.xlarge: ~$150/month per node
- m5.large: ~$88/month per node

**SPOT (approximate):**
- t3.large: ~$20-30/month per node
- t3.xlarge: ~$40-50/month per node
- m5.large: ~$25-35/month per node

## Known Limitations

- Launch template updates require node group recreation
- Maximum 3 node groups recommended per cluster for maintainability
- SPOT instances can be interrupted (use for fault-tolerant workloads only)
- IMDSv2 requirement may break legacy applications

## Upgrade Guide

### Kubernetes Version

```hcl
# Update cluster_version
cluster_version = "1.29"

# Apply changes
terraform plan
terraform apply
```

Nodes will be gradually replaced with new version.

### Instance Types

```hcl
# Update instance_types
instance_types = ["t3.xlarge", "t3a.xlarge"]

# Apply changes
terraform apply
```

This will trigger a rolling update of nodes.

## Contributing

Contributions welcome! Please submit issues and pull requests.

## License

MIT License - See [LICENSE](../../LICENSE) for details.

## Related Modules

- [eks-cluster](../eks-cluster) - EKS control plane
- [rds](../rds) - RDS database
- [iam-roles](../iam-roles) - IRSA roles

## Support

- **Issues**: [GitHub Issues](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/issues)
- **Discussions**: [GitHub Discussions](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/discussions)

## Version

**Current Version**: 1.0.0

**Terraform Registry**: `github.com/asarkar157/Multi-AZ-EKS-Cluster//modules/eks-node-groups?ref=eks-node-groups-v1.0.0`
