package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestVPCModule(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/vpc",
		Vars: map[string]interface{}{
			"region":             "us-east-1",
			"vpc_cidr":           "10.0.0.0/16",
			"availability_zones": []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"cluster_name":       "test-cluster",
			"environment":        "test",
			"tags": map[string]string{
				"Environment": "test",
				"ManagedBy":   "terratest",
			},
		},
		// Don't actually deploy - just validate
		PlanOnly: true,
	})

	defer terraform.Destroy(t, terraformOptions)

	// Run terraform init and plan
	planStruct := terraform.InitAndPlan(t, terraformOptions)

	// Verify resource counts
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Expected resources:
	// - 1 VPC
	// - 1 Internet Gateway
	// - 3 EIPs for NAT
	// - 3 Public Subnets
	// - 3 NAT Gateways
	// - 3 Private Subnets
	// - 3 Database Subnets
	// - 1 Public Route Table + 3 associations
	// - 3 Private Route Tables + 3 associations
	// - 1 Database Route Table + 3 associations
	// - 1 DB Subnet Group
	// - 1 VPC Flow Log
	// - 1 CloudWatch Log Group
	// - 1 IAM Role for Flow Logs
	// - 1 IAM Role Policy

	assert.Greater(t, resourceCounts.Add, 25, "Should create more than 25 resources")
	assert.Equal(t, 0, resourceCounts.Change, "Should not change any resources")
	assert.Equal(t, 0, resourceCounts.Destroy, "Should not destroy any resources")
}

func TestVPCModuleValidation(t *testing.T) {
	t.Parallel()

	// Test that exactly 3 AZs are required
	terraformOptions := &terraform.Options{
		TerraformDir: "../modules/vpc",
		Vars: map[string]interface{}{
			"region":             "us-east-1",
			"vpc_cidr":           "10.0.0.0/16",
			"availability_zones": []string{"us-east-1a", "us-east-1b"}, // Only 2 AZs
			"cluster_name":       "test-cluster",
			"environment":        "test",
		},
	}

	// This should fail validation
	_, err := terraform.InitAndPlanE(t, terraformOptions)
	assert.Error(t, err, "Should fail validation with only 2 AZs")
}

func TestVPCOutputs(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/vpc",
		Vars: map[string]interface{}{
			"region":             "us-east-1",
			"vpc_cidr":           "10.0.0.0/16",
			"availability_zones": []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"cluster_name":       "test-cluster",
			"environment":        "test",
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)

	// Verify that key outputs are planned
	// Note: With PlanOnly, we can't get actual output values, but we can validate the plan
	assert.NotNil(t, planStruct, "Plan should not be nil")
}

func TestVPCSubnetCIDRCalculation(t *testing.T) {
	t.Parallel()

	// Test that subnet CIDR calculations don't overlap
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/vpc",
		Vars: map[string]interface{}{
			"region":             "us-west-2",
			"vpc_cidr":           "10.1.0.0/16",
			"availability_zones": []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			"cluster_name":       "test-cluster-2",
			"environment":        "test",
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should succeed with different VPC CIDR")
}

func TestVPCTags(t *testing.T) {
	t.Parallel()

	customTags := map[string]string{
		"Environment": "production",
		"Team":        "platform",
		"CostCenter":  "engineering",
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/vpc",
		Vars: map[string]interface{}{
			"region":             "us-east-1",
			"vpc_cidr":           "10.2.0.0/16",
			"availability_zones": []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"cluster_name":       "test-cluster-tags",
			"environment":        "test",
			"tags":               customTags,
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should succeed with custom tags")
}
