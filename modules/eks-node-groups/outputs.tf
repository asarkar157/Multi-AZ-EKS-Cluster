output "node_groups" {
  description = "Map of node group attributes"
  value = {
    for k, ng in aws_eks_node_group.main : k => {
      id             = ng.id
      arn            = ng.arn
      status         = ng.status
      capacity_type  = ng.capacity_type
      instance_types = ng.instance_types
    }
  }
}

output "node_security_group_id" {
  description = "Security group ID for node groups"
  value       = aws_security_group.node_group.id
}

output "node_iam_role_arn" {
  description = "IAM role ARN for node groups"
  value       = aws_iam_role.node_group.arn
}

output "node_iam_role_name" {
  description = "IAM role name for node groups"
  value       = aws_iam_role.node_group.name
}
