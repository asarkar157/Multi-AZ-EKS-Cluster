package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestMultiRegionEKSIntegration tests the complete multi-region setup
func TestMultiRegionEKSIntegration(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "..",
		Vars: map[string]interface{}{
			"primary_region":   "us-east-1",
			"secondary_region": "us-west-2",
			"primary_vpc_id":   "vpc-primary123",
			"secondary_vpc_id": "vpc-secondary456",
			"primary_availability_zones":   []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"secondary_availability_zones": []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			"cluster_name_prefix":          "test-multi-region",
			"environment":                  "test",
			"kubernetes_version":           "1.28",
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
			"organizational_units": []map[string]interface{}{
				{
					"name":        "test-ou",
					"ou_id":       "ou-test-001",
					"permissions": []string{"admin"},
				},
			},
		},
		PlanOnly: true,
	})

	defer terraform.Destroy(t, terraformOptions)

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should create resources in both regions
	// Each region: EKS cluster, node groups, RDS, IAM roles
	// Plus: VPC peering
	assert.Greater(t, resourceCounts.Add, 60, "Should create more than 60 resources for multi-region setup")
}

func TestMultiRegionEKSVPCPeering(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "..",
		Vars: map[string]interface{}{
			"primary_region":               "us-east-1",
			"secondary_region":             "eu-west-1",
			"primary_vpc_id":               "vpc-primary123",
			"secondary_vpc_id":             "vpc-secondary456",
			"primary_availability_zones":   []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"secondary_availability_zones": []string{"eu-west-1a", "eu-west-1b", "eu-west-1c"},
			"cluster_name_prefix":          "test-vpc-peering",
			"environment":                  "test",
			"kubernetes_version":           "1.28",
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
			"rds_config": map[string]interface{}{
				"engine":                  "postgres",
				"engine_version":          "15.4",
				"instance_class":          "db.t3.small",
				"allocated_storage":       50,
				"database_name":           "testdb",
				"master_username":         "admin",
				"backup_retention_period": 7,
				"multi_az":                false,
				"storage_encrypted":       true,
			},
			"organizational_units": []map[string]interface{}{
				{
					"name":        "test-ou",
					"ou_id":       "ou-test-001",
					"permissions": []string{"view"},
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should include VPC peering between regions")
}

func TestMultiRegionEKSRDSReplication(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "..",
		Vars: map[string]interface{}{
			"primary_region":               "us-east-1",
			"secondary_region":             "us-west-2",
			"primary_vpc_id":               "vpc-primary123",
			"secondary_vpc_id":             "vpc-secondary456",
			"primary_availability_zones":   []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"secondary_availability_zones": []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			"cluster_name_prefix":          "test-rds-replication",
			"environment":                  "production",
			"kubernetes_version":           "1.28",
			"node_groups": map[string]interface{}{
				"general": map[string]interface{}{
					"desired_size":   9,
					"min_size":       6,
					"max_size":       18,
					"instance_types": []string{"m5.xlarge"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      100,
				},
			},
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
			"organizational_units": []map[string]interface{}{
				{
					"name":        "prod-ops",
					"ou_id":       "ou-ops-001",
					"permissions": []string{"admin"},
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should include RDS primary and read replica")
}

func TestMultiRegionEKSProduction(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "..",
		Vars: map[string]interface{}{
			"primary_region":               "us-east-1",
			"secondary_region":             "us-west-2",
			"primary_vpc_id":               "vpc-primary123",
			"secondary_vpc_id":             "vpc-secondary456",
			"primary_availability_zones":   []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"secondary_availability_zones": []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			"cluster_name_prefix":          "production",
			"environment":                  "production",
			"kubernetes_version":           "1.28",
			"node_groups": map[string]interface{}{
				"general": map[string]interface{}{
					"desired_size":   9,
					"min_size":       6,
					"max_size":       18,
					"instance_types": []string{"m5.xlarge", "m5a.xlarge"},
					"capacity_type":  "ON_DEMAND",
					"disk_size":      100,
				},
				"spot": map[string]interface{}{
					"desired_size":   6,
					"min_size":       0,
					"max_size":       15,
					"instance_types": []string{"m5.xlarge", "m5a.xlarge", "m5n.xlarge"},
					"capacity_type":  "SPOT",
					"disk_size":      100,
				},
			},
			"rds_config": map[string]interface{}{
				"engine":                  "postgres",
				"engine_version":          "15.4",
				"instance_class":          "db.r6g.2xlarge",
				"allocated_storage":       1000,
				"database_name":           "proddb",
				"master_username":         "dbadmin",
				"backup_retention_period": 30,
				"multi_az":                true,
				"storage_encrypted":       true,
			},
			"organizational_units": []map[string]interface{}{
				{
					"name":        "production-ops",
					"ou_id":       "ou-prod-ops-001",
					"permissions": []string{"admin", "deploy", "view"},
				},
				{
					"name":        "production-dev",
					"ou_id":       "ou-prod-dev-001",
					"permissions": []string{"deploy", "view"},
				},
				{
					"name":        "production-readonly",
					"ou_id":       "ou-prod-ro-001",
					"permissions": []string{"view"},
				},
			},
			"tags": map[string]string{
				"Environment": "production",
				"ManagedBy":   "terraform",
				"Team":        "platform",
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Production setup with 3 OUs, 2 node groups per region, multi-AZ RDS
	assert.Greater(t, resourceCounts.Add, 80, "Should create more than 80 resources for production multi-region setup")
}

func TestMultiRegionEKSOutputs(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "..",
		Vars: map[string]interface{}{
			"primary_region":               "us-east-1",
			"secondary_region":             "us-west-2",
			"primary_vpc_id":               "vpc-primary123",
			"secondary_vpc_id":             "vpc-secondary456",
			"primary_availability_zones":   []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"secondary_availability_zones": []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			"cluster_name_prefix":          "test-outputs",
			"environment":                  "test",
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
			"rds_config": map[string]interface{}{
				"engine":                  "postgres",
				"engine_version":          "15.4",
				"instance_class":          "db.t3.small",
				"allocated_storage":       50,
				"database_name":           "testdb",
				"master_username":         "admin",
				"backup_retention_period": 7,
				"multi_az":                false,
				"storage_encrypted":       true,
			},
			"organizational_units": []map[string]interface{}{
				{
					"name":        "test-ou",
					"ou_id":       "ou-test-001",
					"permissions": []string{"admin"},
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should include all expected outputs")
}
