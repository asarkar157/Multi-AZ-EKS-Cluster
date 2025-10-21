# Multi-Region EKS Cluster with Multi-AZ and Shared RDS

This Terraform module creates a production-ready, highly available multi-region EKS cluster infrastructure with the following features:

- **Multi-Region Deployment**: Primary and secondary regions for disaster recovery
- **Multi-AZ Architecture**: Each region spans exactly 3 availability zones for high availability
- **Production OU Support**: Multiple organizational units with role-based access control
- **Shared RDS Access**: RDS instance with multi-AZ deployment accessible from EKS clusters
- **Security Best Practices**: Encryption at rest, IRSA, VPC Flow Logs, and secure networking

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Primary Region                            │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   AZ-1       │  │   AZ-2       │  │   AZ-3       │          │
│  ├──────────────┤  ├──────────────┤  ├──────────────┤          │
│  │ EKS Nodes    │  │ EKS Nodes    │  │ EKS Nodes    │          │
│  │ Private Sub  │  │ Private Sub  │  │ Private Sub  │          │
│  │ Public Sub   │  │ Public Sub   │  │ Public Sub   │          │
│  │ DB Sub       │  │ DB Sub       │  │ DB Sub       │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│         │                 │                 │                   │
│         └─────────────────┴─────────────────┘                   │
│                           │                                     │
│                    ┌──────▼──────┐                              │
│                    │ RDS Primary │                              │
│                    │  Multi-AZ   │                              │
│                    └─────────────┘                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                    VPC Peering Connection
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                       Secondary Region                           │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   AZ-1       │  │   AZ-2       │  │   AZ-3       │          │
│  ├──────────────┤  ├──────────────┤  ├──────────────┤          │
│  │ EKS Nodes    │  │ EKS Nodes    │  │ EKS Nodes    │          │
│  │ Private Sub  │  │ Private Sub  │  │ Private Sub  │          │
│  │ Public Sub   │  │ Public Sub   │  │ Public Sub   │          │
│  │ DB Sub       │  │ DB Sub       │  │ DB Sub       │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│         │                 │                 │                   │
│         └─────────────────┴─────────────────┘                   │
│                           │                                     │
│                    ┌──────▼──────┐                              │
│                    │ RDS Replica │                              │
│                    │  Multi-AZ   │                              │
│                    └─────────────┘                              │
└─────────────────────────────────────────────────────────────────┘
```

## Features

### High Availability
- **3 AZs per Region**: Node groups and subnets distributed across 3 availability zones
- **Multi-Region**: Primary and secondary regions with VPC peering
- **RDS Multi-AZ**: Database instances deployed in multi-AZ mode
- **Multiple NAT Gateways**: One NAT gateway per AZ for redundancy

### Security
- **Encryption at Rest**: KMS encryption for EKS secrets and RDS storage
- **Encryption in Transit**: TLS for all communications
- **IRSA (IAM Roles for Service Accounts)**: Fine-grained IAM permissions for pods
- **Security Groups**: Least-privilege network access controls
- **Secrets Management**: RDS credentials stored in AWS Secrets Manager
- **VPC Flow Logs**: Network traffic monitoring

### Access Control
- **Organizational Units**: Multiple production OUs with different permission levels
  - Admin: Full cluster access
  - Deploy: Deployment and edit permissions
  - View: Read-only access
- **RBAC Integration**: EKS access entries for OU-based access

### Observability
- **CloudWatch Logs**: EKS control plane logs
- **Performance Insights**: RDS performance monitoring
- **Enhanced Monitoring**: RDS instance metrics
- **VPC Flow Logs**: Network traffic analysis

## Prerequisites

- Terraform >= 1.0
- AWS CLI configured with appropriate credentials
- Existing VPCs in both regions with the following subnets tagged:
  - Private subnets: `Type=private` (for EKS nodes)
  - Database subnets: `Type=database` (for RDS)
- AWS Organizations setup with organizational units (if using OU-based access)

## Required VPC Subnet Tags

Your existing VPCs must have subnets tagged appropriately:

```hcl
# Private subnets (for EKS nodes)
tags = {
  Type = "private"
}

# Database subnets (for RDS)
tags = {
  Type = "database"
}
```

## Usage

### Basic Example

```hcl
module "multi_region_eks" {
  source = "."

  # Region Configuration
  primary_region   = "us-east-1"
  secondary_region = "us-west-2"

  # Existing VPC IDs
  primary_vpc_id   = "vpc-xxxxx"
  secondary_vpc_id = "vpc-yyyyy"

  # Availability Zones (must specify exactly 3 per region)
  primary_availability_zones   = ["us-east-1a", "us-east-1b", "us-east-1c"]
  secondary_availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

  # Cluster Configuration
  cluster_name_prefix = "my-company"
  environment         = "production"

  # Organizational Units
  organizational_units = [
    {
      name        = "production-ops"
      ou_id       = "ou-prod-ops-001"
      permissions = ["admin", "deploy", "view"]
    },
    {
      name        = "production-dev"
      ou_id       = "ou-prod-dev-001"
      permissions = ["deploy", "view"]
    },
    {
      name        = "production-readonly"
      ou_id       = "ou-prod-ro-001"
      permissions = ["view"]
    }
  ]

  tags = {
    Environment = "production"
    ManagedBy   = "terraform"
    Team        = "platform"
  }
}
```

### Advanced Configuration

```hcl
module "multi_region_eks" {
  source = "."

  primary_region   = "us-east-1"
  secondary_region = "eu-west-1"

  primary_vpc_id   = "vpc-xxxxx"
  secondary_vpc_id = "vpc-yyyyy"

  primary_availability_zones   = ["us-east-1a", "us-east-1b", "us-east-1d"]
  secondary_availability_zones = ["eu-west-1a", "eu-west-1b", "eu-west-1c"]

  # Custom EKS Configuration
  kubernetes_version = "1.28"

  node_groups = {
    general = {
      desired_size   = 9   # 3 per AZ
      min_size       = 6   # 2 per AZ
      max_size       = 18  # 6 per AZ
      instance_types = ["m5.xlarge", "m5a.xlarge"]
      capacity_type  = "ON_DEMAND"
      disk_size      = 100
    }
    spot = {
      desired_size   = 6   # 2 per AZ
      min_size       = 3   # 1 per AZ
      max_size       = 15  # 5 per AZ
      instance_types = ["m5.xlarge", "m5a.xlarge", "m5n.xlarge"]
      capacity_type  = "SPOT"
      disk_size      = 100
    }
  }

  # Custom RDS Configuration
  rds_config = {
    engine                  = "postgres"
    engine_version          = "15.4"
    instance_class          = "db.r6g.2xlarge"
    allocated_storage       = 500
    database_name           = "myapp"
    master_username         = "dbadmin"
    backup_retention_period = 30
    multi_az                = true
    storage_encrypted       = true
  }

  organizational_units = [
    {
      name        = "platform-team"
      ou_id       = "ou-platform-001"
      permissions = ["admin"]
    },
    {
      name        = "engineering"
      ou_id       = "ou-eng-001"
      permissions = ["deploy", "view"]
    },
    {
      name        = "sre"
      ou_id       = "ou-sre-001"
      permissions = ["admin", "deploy", "view"]
    }
  ]
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| primary_region | Primary AWS region | `string` | `"us-east-1"` | no |
| secondary_region | Secondary AWS region | `string` | `"us-west-2"` | no |
| primary_vpc_id | ID of existing VPC in primary region | `string` | n/a | yes |
| secondary_vpc_id | ID of existing VPC in secondary region | `string` | n/a | yes |
| primary_availability_zones | Availability zones for primary region (must be exactly 3) | `list(string)` | `["us-east-1a", "us-east-1b", "us-east-1c"]` | no |
| secondary_availability_zones | Availability zones for secondary region (must be exactly 3) | `list(string)` | `["us-west-2a", "us-west-2b", "us-west-2c"]` | no |
| cluster_name_prefix | Prefix for EKS cluster names | `string` | `"multi-region-eks"` | no |
| environment | Environment name | `string` | `"production"` | no |
| kubernetes_version | Kubernetes version | `string` | `"1.28"` | no |
| node_groups | Configuration for EKS node groups | `map(object)` | See variables.tf | no |
| rds_config | RDS instance configuration | `object` | See variables.tf | no |
| organizational_units | List of organizational units for access control | `list(object)` | See variables.tf | no |
| tags | Common tags to apply to all resources | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| primary_cluster_endpoint | Endpoint for primary EKS cluster |
| primary_cluster_name | Name of the primary EKS cluster |
| secondary_cluster_endpoint | Endpoint for secondary EKS cluster |
| secondary_cluster_name | Name of the secondary EKS cluster |
| primary_rds_endpoint | Connection endpoint for primary RDS instance |
| secondary_rds_endpoint | Connection endpoint for secondary RDS instance |
| vpc_peering_connection_id | ID of the VPC peering connection |

## Post-Deployment Steps

### 1. Configure kubectl Access

```bash
# Primary cluster
aws eks update-kubeconfig --region us-east-1 --name <primary-cluster-name>

# Secondary cluster
aws eks update-kubeconfig --region us-west-2 --name <secondary-cluster-name>
```

### 2. Retrieve RDS Credentials

```bash
# Get the secret ARN from Terraform outputs
aws secretsmanager get-secret-value --secret-id <secret-arn> --region us-east-1
```

### 3. Install Essential Add-ons

```bash
# AWS Load Balancer Controller
helm repo add eks https://aws.github.io/eks-charts
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=<cluster-name> \
  --set serviceAccount.annotations."eks\.amazonaws\.com/role-arn"=<alb-controller-role-arn>

# Cluster Autoscaler
kubectl apply -f https://raw.githubusercontent.com/kubernetes/autoscaler/master/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml
kubectl -n kube-system annotate serviceaccount cluster-autoscaler \
  eks.amazonaws.com/role-arn=<cluster-autoscaler-role-arn>

# External DNS (optional)
kubectl create serviceaccount external-dns -n kube-system
kubectl annotate serviceaccount external-dns -n kube-system \
  eks.amazonaws.com/role-arn=<external-dns-role-arn>
```

### 4. Configure Database Access

Create a Kubernetes secret for your application to access RDS:

```bash
kubectl create secret generic rds-credentials \
  --from-literal=host=<rds-endpoint> \
  --from-literal=username=<master-username> \
  --from-literal=password=<password> \
  --from-literal=database=<database-name>
```

## Module Structure

```
.
├── main.tf                      # Root module configuration
├── variables.tf                 # Root module variables
├── outputs.tf                   # Root module outputs
├── README.md                    # This file
├── examples/                    # Example configurations
│   ├── basic/
│   └── advanced/
└── modules/
    ├── regional-eks/            # Regional EKS cluster module
    │   ├── main.tf
    │   ├── variables.tf
    │   └── outputs.tf
    ├── eks-cluster/             # EKS cluster configuration
    │   ├── main.tf
    │   ├── variables.tf
    │   └── outputs.tf
    ├── eks-node-groups/         # EKS node groups
    │   ├── main.tf
    │   ├── variables.tf
    │   ├── outputs.tf
    │   └── user_data.sh
    ├── rds/                     # RDS instance module
    │   ├── main.tf
    │   ├── variables.tf
    │   └── outputs.tf
    ├── iam-roles/               # IAM roles for IRSA
    │   ├── main.tf
    │   ├── variables.tf
    │   ├── outputs.tf
    │   └── policies/
    │       └── alb-controller-policy.json
    └── vpc/                     # VPC module (reference only)
        ├── main.tf
        ├── variables.tf
        └── outputs.tf
```

## Cost Considerations

This infrastructure creates the following billable resources:

**Per Region:**
- EKS Cluster: ~$73/month
- NAT Gateways (3): ~$97/month each = ~$291/month
- EKS Nodes: Varies by instance type and count
- RDS Instance: Varies by instance class
- Data Transfer: Varies by usage

**Cross-Region:**
- VPC Peering: Data transfer charges apply
- RDS Replication: Data transfer charges apply

**Estimated Monthly Cost (minimum):**
- Primary Region: ~$500-1000
- Secondary Region: ~$500-1000
- **Total: ~$1000-2000/month** (excluding compute and data transfer)

## Best Practices

1. **Use Spot Instances**: Configure spot node groups for non-critical workloads
2. **Enable Cluster Autoscaler**: Automatically scale nodes based on demand
3. **Monitor Costs**: Use AWS Cost Explorer and set up billing alerts
4. **Regular Backups**: RDS automated backups are enabled by default
5. **Security Scanning**: Scan container images before deployment
6. **Network Policies**: Implement Kubernetes network policies
7. **Resource Quotas**: Set resource quotas and limits for namespaces
8. **Pod Security**: Use Pod Security Standards

## Troubleshooting

### Nodes Not Joining Cluster

Check the node group status:
```bash
aws eks describe-nodegroup --cluster-name <cluster-name> --nodegroup-name <nodegroup-name>
```

### RDS Connection Issues

1. Verify security group rules allow EKS node security group
2. Check RDS endpoint is correct
3. Verify credentials from Secrets Manager

### VPC Peering Not Working

1. Verify peering connection is active
2. Check route tables have routes to peer VPC
3. Verify security groups allow cross-VPC traffic

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

MIT License - See LICENSE file for details

## Support

For issues and questions:
- Open a GitHub issue
- Contact your AWS support team
- Check AWS EKS documentation: https://docs.aws.amazon.com/eks/

## Changelog

### Version 1.0.0
- Initial release
- Multi-region EKS cluster support
- Multi-AZ architecture
- RDS integration with multi-AZ
- OU-based access control
- IRSA for service accounts
