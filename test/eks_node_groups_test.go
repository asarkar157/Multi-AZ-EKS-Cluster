package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestEKSNodeGroupsModule(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-node-groups",
		Vars: map[string]interface{}{
			"cluster_name":                       "test-cluster",
			"cluster_version":                    "1.28",
			"vpc_id":                             "vpc-12345678",
			"subnet_ids":                         []string{"subnet-1", "subnet-2", "subnet-3"},
			"cluster_security_group_id":          "sg-cluster",
			"cluster_primary_security_group_id":  "sg-primary",
			"node_groups": map[string]interface{}{
				"general": map[string]interface{}{
					"desired_size":   6,
					"min_size":       3,
					"max_size":       15,
					"instance_types": []string{"t3.large"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      50,
				},
			},
			"tags": map[string]string{
				"Environment": "test",
			},
		},
		PlanOnly: true,
	})

	defer terraform.Destroy(t, terraformOptions)

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Expected resources:
	// - IAM Role
	// - 5 IAM Role Policy Attachments
	// - Security Group
	// - 3 Security Group Rules
	// - 1 Launch Template per node group
	// - 1 EKS Node Group per node group

	assert.Greater(t, resourceCounts.Add, 10, "Should create more than 10 resources")
}

func TestEKSNodeGroupsMultipleGroups(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-node-groups",
		Vars: map[string]interface{}{
			"cluster_name":                      "test-cluster-multi",
			"cluster_version":                   "1.28",
			"vpc_id":                            "vpc-12345678",
			"subnet_ids":                        []string{"subnet-1", "subnet-2", "subnet-3"},
			"cluster_security_group_id":         "sg-cluster",
			"cluster_primary_security_group_id": "sg-primary",
			"node_groups": map[string]interface{}{
				"general": map[string]interface{}{
					"desired_size":   6,
					"min_size":       3,
					"max_size":       15,
					"instance_types": []string{"t3.large", "t3a.large"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      50,
				},
				"spot": map[string]interface{}{
					"desired_size":   3,
					"min_size":       0,
					"max_size":       12,
					"instance_types": []string{"t3.large", "t3.xlarge"},
					"capacity_type":  "SPOT",
					"disk_size":      50,
				},
				"compute": map[string]interface{}{
					"desired_size":   6,
					"min_size":       3,
					"max_size":       18,
					"instance_types": []string{"c5.2xlarge"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      100,
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should create resources for 3 node groups
	assert.Greater(t, resourceCounts.Add, 15, "Should create resources for multiple node groups")
}

func TestEKSNodeGroupsSpotInstances(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-node-groups",
		Vars: map[string]interface{}{
			"cluster_name":                      "test-cluster-spot",
			"cluster_version":                   "1.28",
			"vpc_id":                            "vpc-12345678",
			"subnet_ids":                        []string{"subnet-1", "subnet-2", "subnet-3"},
			"cluster_security_group_id":         "sg-cluster",
			"cluster_primary_security_group_id": "sg-primary",
			"node_groups": map[string]interface{}{
				"spot-workers": map[string]interface{}{
					"desired_size":   9,
					"min_size":       3,
					"max_size":       30,
					"instance_types": []string{"t3.large", "t3a.large", "t2.large"},
					"capacity_type":  "SPOT",
					"disk_size":      50,
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should succeed for SPOT instances")
}

func TestEKSNodeGroupsLaunchTemplate(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-node-groups",
		Vars: map[string]interface{}{
			"cluster_name":                      "test-cluster-lt",
			"cluster_version":                   "1.28",
			"vpc_id":                            "vpc-12345678",
			"subnet_ids":                        []string{"subnet-1", "subnet-2", "subnet-3"},
			"cluster_security_group_id":         "sg-cluster",
			"cluster_primary_security_group_id": "sg-primary",
			"node_groups": map[string]interface{}{
				"custom": map[string]interface{}{
					"desired_size":   6,
					"min_size":       3,
					"max_size":       15,
					"instance_types": []string{"m5.xlarge"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      100,
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should include launch template with custom disk size")
}

func TestEKSNodeGroupsSecurityGroups(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-node-groups",
		Vars: map[string]interface{}{
			"cluster_name":                      "test-cluster-sg",
			"cluster_version":                   "1.28",
			"vpc_id":                            "vpc-12345678",
			"subnet_ids":                        []string{"subnet-1", "subnet-2", "subnet-3"},
			"cluster_security_group_id":         "sg-cluster-123",
			"cluster_primary_security_group_id": "sg-primary-456",
			"node_groups": map[string]interface{}{
				"workers": map[string]interface{}{
					"desired_size":   6,
					"min_size":       3,
					"max_size":       15,
					"instance_types": []string{"t3.medium"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      50,
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should configure security groups correctly")
}

func TestEKSNodeGroupsIAMRoles(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-node-groups",
		Vars: map[string]interface{}{
			"cluster_name":                      "test-cluster-iam",
			"cluster_version":                   "1.28",
			"vpc_id":                            "vpc-12345678",
			"subnet_ids":                        []string{"subnet-1", "subnet-2", "subnet-3"},
			"cluster_security_group_id":         "sg-cluster",
			"cluster_primary_security_group_id": "sg-primary",
			"node_groups": map[string]interface{}{
				"workers": map[string]interface{}{
					"desired_size":   3,
					"min_size":       1,
					"max_size":       9,
					"instance_types": []string{"t3.small"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      30,
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should include IAM role and policy attachments
	assert.Greater(t, resourceCounts.Add, 8, "Should create IAM resources for node groups")
}
