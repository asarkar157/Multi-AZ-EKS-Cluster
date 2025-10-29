# EKS Node Groups Module v1.0.0

Production-ready Terraform module for AWS EKS managed node groups with multi-AZ distribution, custom launch templates, and support for both ON_DEMAND and SPOT capacity.

## 🚀 Quick Start

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

    spot = {
      desired_size   = 3
      min_size       = 0
      max_size       = 12
      instance_types = ["t3.large", "t3a.large", "t3.xlarge"]
      capacity_type  = "SPOT"
      disk_size      = 50
    }
  }

  tags = {
    Environment = "production"
    ManagedBy   = "terraform"
  }
}
```

## ✨ Key Features

### Multi-AZ High Availability
- ✅ **Automatic Distribution** - Nodes spread across all availability zones
- ✅ **Even Placement** - EKS distributes nodes evenly across AZs
- ✅ **3 AZ Support** - Optimized for 3 availability zones per region

### Flexible Node Groups
- ✅ **Multiple Groups** - Deploy multiple node groups with different configs
- ✅ **ON_DEMAND** - Guaranteed capacity for critical workloads
- ✅ **SPOT** - Up to 90% cost savings for fault-tolerant workloads
- ✅ **Mixed Instance Types** - Multiple instance types per group

### Custom Launch Templates
- ✅ **Encrypted EBS Volumes** - gp3 volumes with 3000 IOPS
- ✅ **IMDSv2 Enforcement** - Enhanced instance metadata security
- ✅ **Detailed Monitoring** - CloudWatch metrics enabled
- ✅ **Custom Disk Sizes** - Configurable per node group

### Security & IAM
- ✅ **5 AWS Managed Policies** - All required permissions included
- ✅ **Security Groups** - Least-privilege network rules
- ✅ **Systems Manager** - SSM access for debugging
- ✅ **ECR Integration** - Container image pulling enabled

### Autoscaling
- ✅ **Independent Scaling** - Per node group min/max/desired
- ✅ **Cluster Autoscaler Ready** - Ignores desired_size changes
- ✅ **Rolling Updates** - 33% max unavailable during updates

## 📦 What's Included

### Module Files
```
modules/eks-node-groups/
├── main.tf          # Resource definitions
├── variables.tf     # Input variables
├── outputs.tf       # Module outputs
├── user_data.sh     # Node bootstrap script
├── README.md        # Complete documentation
└── CHANGELOG.md     # Version history
```

### Resources Created

**Base Resources (10):**
- 1 IAM Role for nodes (shared)
- 5 IAM Policy Attachments:
  - AmazonEKSWorkerNodePolicy
  - AmazonEKS_CNI_Policy
  - AmazonEC2ContainerRegistryReadOnly
  - AmazonSSMManagedInstanceCore
  - AmazonEBSCSIDriverPolicy
- 1 Security Group
- 3 Security Group Rules

**Per Node Group (2 each):**
- 1 Launch Template
- 1 EKS Node Group

**Total:** 10 + (2 × number of node groups)

## 🔒 Security Features

- **Encrypted Volumes**: All EBS volumes encrypted by default
- **IMDSv2 Required**: Enhanced instance metadata security
- **Security Groups**: Least-privilege communication rules
- **IAM Roles**: Minimal required permissions
- **SSM Access**: Secure shell access via Systems Manager

## 📊 Requirements

| Component | Version |
|-----------|---------|
| Terraform | >= 1.0 |
| AWS Provider | ~> 5.0 |
| EKS Cluster | Required (must exist) |

## 🎯 Use Cases

### General Purpose Workloads
```hcl
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
```

### Cost-Optimized with SPOT
```hcl
node_groups = {
  spot_workers = {
    desired_size   = 9
    min_size       = 3
    max_size       = 30
    instance_types = ["t3.large", "t3a.large", "t2.large"]
    capacity_type  = "SPOT"
    disk_size      = 50
  }
}
```

### Dedicated Workload Groups
```hcl
node_groups = {
  system = {
    desired_size   = 3
    min_size       = 3
    max_size       = 3
    instance_types = ["t3.medium"]
    capacity_type  = "ON_DEMAND"
    disk_size      = 30
  }

  application = {
    desired_size   = 6
    min_size       = 3
    max_size       = 20
    instance_types = ["m5.large"]
    capacity_type  = "ON_DEMAND"
    disk_size      = 50
  }

  batch = {
    desired_size   = 0
    min_size       = 0
    max_size       = 10
    instance_types = ["c5.2xlarge"]
    capacity_type  = "SPOT"
    disk_size      = 100
  }
}
```

## 📝 Inputs

### Required

| Name | Description | Type |
|------|-------------|------|
| cluster_name | EKS cluster name | `string` |
| cluster_version | Kubernetes version | `string` |
| vpc_id | VPC ID | `string` |
| subnet_ids | Subnet IDs (3 AZs) | `list(string)` |
| cluster_security_group_id | Cluster SG ID | `string` |
| cluster_primary_security_group_id | EKS primary SG ID | `string` |
| node_groups | Node group configs | `map(object)` |

### Optional

| Name | Description | Type | Default |
|------|-------------|------|---------|
| tags | Resource tags | `map(string)` | `{}` |

## 📤 Outputs

| Name | Description |
|------|-------------|
| node_groups | Node group details (id, arn, status) |
| node_security_group_id | Node security group ID |
| node_iam_role_arn | IAM role ARN |
| node_iam_role_name | IAM role name |

## 🛠️ Post-Deployment

### Verify Node Groups
```bash
# List node groups
kubectl get nodes -o wide

# View node labels
kubectl get nodes --show-labels

# Check node group status
aws eks describe-nodegroup \
  --cluster-name my-eks-cluster \
  --nodegroup-name general
```

### Install Cluster Autoscaler
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/autoscaler/master/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml

kubectl -n kube-system annotate serviceaccount cluster-autoscaler \
  eks.amazonaws.com/role-arn=<autoscaler-role-arn>
```

## 💰 Cost Optimization

### Instance Pricing (us-east-1)

**ON_DEMAND (per month):**
- t3.medium: ~$30
- t3.large: ~$75
- t3.xlarge: ~$150
- m5.large: ~$88

**SPOT (per month, approximate):**
- t3.large: ~$20-30 (70-90% savings)
- t3.xlarge: ~$40-50
- m5.large: ~$25-35

### Savings Tips
1. **Use SPOT** for non-critical workloads
2. **Right-size instances** - Start small, scale up
3. **Use Cluster Autoscaler** - Scale down during low usage
4. **Multiple instance types** - Better SPOT availability
5. **Schedule scaling** - Scale down dev/test at night

## 🏗️ Architecture

### Multi-AZ Distribution
```
┌────────────────────────────────────────┐
│          EKS Cluster                   │
├────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────┐  ┌─────────┐│
│  │  AZ-1   │  │  AZ-2   │  │  AZ-3   ││
│  ├─────────┤  ├─────────┤  ├─────────┤│
│  │ Node 1  │  │ Node 3  │  │ Node 5  ││
│  │ Node 2  │  │ Node 4  │  │ Node 6  ││
│  └─────────┘  └─────────┘  └─────────┘│
└────────────────────────────────────────┘
```

### Node Group Types
```
┌──────────────────┐
│   ON_DEMAND      │ - Critical workloads
│   Node Group     │ - Guaranteed capacity
└──────────────────┘

┌──────────────────┐
│   SPOT           │ - Cost-optimized
│   Node Group     │ - Fault-tolerant
└──────────────────┘

┌──────────────────┐
│   Specialized    │ - GPU/compute
│   Node Group     │ - Dedicated
└──────────────────┘
```

## 🔄 Upgrade Guide

### Kubernetes Version
```hcl
# Update version
cluster_version = "1.29"

# Apply - triggers rolling update
terraform apply
```

### Instance Types
```hcl
# Update instance types
instance_types = ["t3.xlarge", "t3a.xlarge"]

# Apply - replaces nodes
terraform apply
```

## 🐛 Troubleshooting

### Nodes Not Ready
```bash
# Check node status
kubectl get nodes

# Describe problematic node
kubectl describe node <node-name>

# Check node group health
aws eks describe-nodegroup \
  --cluster-name my-cluster \
  --nodegroup-name general \
  --query 'nodegroup.health'
```

### SPOT Interruptions
```bash
# View interruption events
kubectl get events -A | grep -i spot

# Check node conditions
kubectl get nodes -o json | jq '.items[].status.conditions'
```

### Launch Template Issues
```bash
# List launch templates
aws ec2 describe-launch-templates \
  --filters "Name=tag:eks:cluster-name,Values=my-cluster"

# Check latest version
aws ec2 describe-launch-template-versions \
  --launch-template-id lt-xxxxx \
  --versions '$Latest'
```

## 📚 Documentation

- [README.md](README.md) - Complete module documentation
- [CHANGELOG.md](CHANGELOG.md) - Version history
- [Main Repo](../../) - Full project documentation

## 🎯 Best Practices

1. **Multi-AZ**: Always use 3 availability zones
2. **Node Counts**: Use numbers divisible by 3
3. **SPOT**: Use multiple instance types for availability
4. **Monitoring**: Enable CloudWatch Container Insights
5. **Tagging**: Use consistent tags for cost allocation
6. **Updates**: Use rolling updates, not in-place
7. **Autoscaling**: Set realistic min/max values

## 🔗 Related Modules

- [eks-cluster](../eks-cluster) - EKS control plane
- [rds](../rds) - Database
- [iam-roles](../iam-roles) - Service account roles

## 🤝 Contributing

Contributions welcome! Please submit issues and pull requests.

## 📄 License

MIT License - See [LICENSE](../../LICENSE) for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/issues)
- **Discussions**: [GitHub Discussions](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/discussions)

---

**Module Version**: 1.0.0
**Release Date**: October 29, 2025
**Terraform Registry**: `github.com/asarkar157/Multi-AZ-EKS-Cluster//modules/eks-node-groups?ref=eks-node-groups-v1.0.0`

**Status**: ✅ Production Ready
