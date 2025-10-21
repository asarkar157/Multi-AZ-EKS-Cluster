package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestRegionalEKSModule(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/regional-eks",
		Vars: map[string]interface{}{
			"region":             "us-east-1",
			"cluster_name":       "test-regional-cluster",
			"vpc_id":             "vpc-12345678",
			"availability_zones": []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"environment":        "test",
			"organizational_units": []map[string]interface{}{
				{
					"name":        "test-ou",
					"ou_id":       "ou-test-001",
					"permissions": []string{"admin"},
				},
			},
			"kubernetes_version": "1.28",
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
			"create_rds": true,
			"rds_config": map[string]interface{}{
				"engine":                  "postgres",
				"engine_version":          "15.4",
				"instance_class":          "db.t3.medium",
				"allocated_storage":       100,
				"database_name":           "testdb",
				"master_username":         "dbadmin",
				"backup_retention_period": 7,
				"multi_az":                true,
				"storage_encrypted":       true,
			},
		},
		PlanOnly: true,
	})

	defer terraform.Destroy(t, terraformOptions)

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should create resources for:
	// - EKS cluster
	// - Node groups
	// - RDS
	// - IAM roles
	assert.Greater(t, resourceCounts.Add, 30, "Should create more than 30 resources for complete regional setup")
}

func TestRegionalEKSWithoutRDS(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/regional-eks",
		Vars: map[string]interface{}{
			"region":             "us-west-2",
			"cluster_name":       "test-cluster-no-rds",
			"vpc_id":             "vpc-87654321",
			"availability_zones": []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			"environment":        "test",
			"organizational_units": []map[string]interface{}{
				{
					"name":        "test-ou",
					"ou_id":       "ou-test-001",
					"permissions": []string{"view"},
				},
			},
			"kubernetes_version": "1.28",
			"node_groups": map[string]interface{}{
				"workers": map[string]interface{}{
					"desired_size":   3,
					"min_size":       1,
					"max_size":       9,
					"instance_types": []string{"t3.medium"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      50,
				},
			},
			"create_rds": false, // No RDS
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should create fewer resources without RDS
	assert.Greater(t, resourceCounts.Add, 20, "Should create resources for EKS without RDS")
}

func TestRegionalEKSWithReadReplica(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/regional-eks",
		Vars: map[string]interface{}{
			"region":             "eu-west-1",
			"cluster_name":       "test-cluster-replica",
			"vpc_id":             "vpc-replica123",
			"availability_zones": []string{"eu-west-1a", "eu-west-1b", "eu-west-1c"},
			"environment":        "test",
			"organizational_units": []map[string]interface{}{
				{
					"name":        "test-ou",
					"ou_id":       "ou-test-001",
					"permissions": []string{"deploy", "view"},
				},
			},
			"kubernetes_version": "1.28",
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
			"create_rds":     true,
			"rds_primary_arn": "arn:aws:rds:us-east-1:123456789012:db:primary-db",
			"rds_config": map[string]interface{}{
				"engine":                  "postgres",
				"engine_version":          "15.4",
				"instance_class":          "db.t3.medium",
				"allocated_storage":       100,
				"database_name":           "replicadb",
				"master_username":         "dbadmin",
				"backup_retention_period": 7,
				"multi_az":                true,
				"storage_encrypted":       true,
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should succeed with RDS read replica configuration")
}

func TestRegionalEKSMultipleNodeGroups(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/regional-eks",
		Vars: map[string]interface{}{
			"region":             "us-east-1",
			"cluster_name":       "test-cluster-multi-ng",
			"vpc_id":             "vpc-12345678",
			"availability_zones": []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"environment":        "test",
			"organizational_units": []map[string]interface{}{
				{
					"name":        "test-ou",
					"ou_id":       "ou-test-001",
					"permissions": []string{"admin"},
				},
			},
			"kubernetes_version": "1.28",
			"node_groups": map[string]interface{}{
				"general": map[string]interface{}{
					"desired_size":   6,
					"min_size":       3,
					"max_size":       15,
					"instance_types": []string{"t3.large"},
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
			"create_rds": false,
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should create resources for 3 node groups
	assert.Greater(t, resourceCounts.Add, 25, "Should create resources for multiple node groups")
}

func TestRegionalEKSMultipleOUs(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/regional-eks",
		Vars: map[string]interface{}{
			"region":             "us-west-2",
			"cluster_name":       "test-cluster-multi-ou",
			"vpc_id":             "vpc-12345678",
			"availability_zones": []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			"environment":        "production",
			"organizational_units": []map[string]interface{}{
				{
					"name":        "platform-ops",
					"ou_id":       "ou-ops-001",
					"permissions": []string{"admin"},
				},
				{
					"name":        "engineering",
					"ou_id":       "ou-eng-001",
					"permissions": []string{"deploy", "view"},
				},
				{
					"name":        "sre",
					"ou_id":       "ou-sre-001",
					"permissions": []string{"admin", "deploy", "view"},
				},
			},
			"kubernetes_version": "1.28",
			"node_groups": map[string]interface{}{
				"workers": map[string]interface{}{
					"desired_size":   6,
					"min_size":       3,
					"max_size":       15,
					"instance_types": []string{"m5.large"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      100,
				},
			},
			"create_rds": true,
			"rds_config": map[string]interface{}{
				"engine":                  "postgres",
				"engine_version":          "15.4",
				"instance_class":          "db.r6g.xlarge",
				"allocated_storage":       500,
				"database_name":           "proddb",
				"master_username":         "dbadmin",
				"backup_retention_period": 30,
				"multi_az":                true,
				"storage_encrypted":       true,
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should create IAM resources for 3 OUs plus RDS access
	assert.Greater(t, resourceCounts.Add, 35, "Should create resources for multiple OUs and RDS")
}
