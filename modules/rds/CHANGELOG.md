# Changelog - RDS Module

All notable changes to the RDS module will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.1] - 2025-10-29

### Added
- Validation rule for `subnet_ids` variable to ensure at least 2 subnets are provided for high availability

### Changed
- Improved variable validation to catch configuration errors early

## [0.0.0] - Initial Development

### Added

#### Core Features
- Multi-AZ RDS instance support
- Read replica support for cross-region deployments
- PostgreSQL and MySQL engine support
- KMS encryption for data at rest
- AWS Secrets Manager integration for credentials
- Enhanced monitoring with CloudWatch
- Performance Insights enabled

#### Security
- Encrypted storage using customer-managed KMS keys
- Security group with least-privilege access
- Master password stored in AWS Secrets Manager
- Automatic key rotation enabled

#### High Availability
- Multi-AZ deployment support
- Automated backups with configurable retention
- Point-in-time recovery
- Read replica support

#### Monitoring
- Enhanced monitoring (60-second intervals)
- Performance Insights (7-day retention)
- CloudWatch logs export
- Custom DB parameter groups

#### Configuration
- Configurable instance class
- Configurable allocated storage
- Custom backup windows
- Custom maintenance windows
- Automatic minor version upgrades

### Resources Created

#### Base Resources (10-12)
1. DB subnet group
2. Security group for RDS
3. Security group ingress rules (per allowed SG)
4. Security group egress rule
5. KMS key (if encryption enabled)
6. KMS alias (if encryption enabled)
7. DB parameter group
8. RDS instance (primary or replica)
9. IAM role for monitoring
10. IAM role policy attachment
11. Secrets Manager secret (if primary)
12. Secrets Manager secret version (if primary)

**Total:** 10-12 resources depending on configuration

### Features in Detail

#### Multi-AZ Support
- Automatic failover
- Synchronous replication
- No data loss on failure

#### Read Replicas
- Cross-region replication
- Eventual consistency
- Scale read workloads

#### Encryption
- KMS-encrypted storage
- Encrypted backups
- Encrypted snapshots
- Key rotation enabled

#### Backup & Recovery
- Automated daily backups
- Configurable retention (default 7 days)
- Point-in-time recovery
- Final snapshot on deletion

#### Performance
- gp3 storage type
- Performance Insights
- Enhanced monitoring
- Custom DB parameters

### Configuration

#### Required Inputs
- `identifier` - RDS instance identifier
- `vpc_id` - VPC ID
- `subnet_ids` - Subnet IDs (at least 2 required)
- `availability_zones` - AZ list
- `engine_version` - Database version
- `instance_class` - Instance type
- `allocated_storage` - Storage in GB
- `database_name` - Database name
- `master_username` - Master username

#### Optional Inputs
- `engine` - Database engine (default: postgres)
- `backup_retention_period` - Days (default: 7)
- `multi_az` - Enable multi-AZ (default: true)
- `storage_encrypted` - Enable encryption (default: true)
- `replicate_source_db` - Source ARN for replica
- `allowed_security_group_ids` - Security groups
- `tags` - Resource tags

#### Outputs
- `db_instance_endpoint` - Connection endpoint
- `db_instance_arn` - Instance ARN
- `db_instance_id` - Instance identifier
- `db_security_group_id` - Security group ID
- `secret_arn` - Secrets Manager ARN

### Best Practices

#### Instance Sizing
- Start with db.t3.medium for development
- Use db.r5 or db.r6 for production
- Monitor CPU and memory usage

#### Storage
- gp3 storage for cost-effective IOPS
- Start with 100GB minimum
- Plan for growth (storage cannot shrink)

#### Backups
- Enable automated backups
- Use 7+ days retention for production
- Test restore procedures regularly

#### Security
- Always enable encryption
- Use Secrets Manager for passwords
- Limit security group access
- Enable deletion protection

### Known Limitations
- Cannot modify storage type after creation
- Multi-AZ change requires downtime
- Read replicas must use same or newer version
- Deletion protection prevents `terraform destroy`

### Dependencies

#### Required Providers
- `hashicorp/aws` ~> 5.0
- `hashicorp/random` ~> 3.0

#### Terraform Version
- Terraform >= 1.0

#### Module Dependencies
- VPC with subnets
- Security groups from EKS

---

## Support

For issues, questions, or contributions:
- GitHub Issues: [Report bugs or request features](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/issues)
- Pull Requests: [Contribute improvements](https://github.com/asarkar157/Multi-AZ-EKS-Cluster/pulls)

[0.0.1]: https://github.com/asarkar157/Multi-AZ-EKS-Cluster/releases/tag/rds-v0.0.1
