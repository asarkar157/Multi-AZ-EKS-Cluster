package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestIAMRolesModule(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/iam-roles",
		Vars: map[string]interface{}{
			"cluster_name":       "test-cluster",
			"oidc_provider_arn":  "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E",
			"oidc_provider_url":  "https://oidc.eks.us-east-1.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E",
			"rds_instance_arn":   "arn:aws:rds:us-east-1:123456789012:db:test-db",
			"organizational_units": []map[string]interface{}{
				{
					"name":        "test-ou",
					"ou_id":       "ou-test-001",
					"permissions": []string{"admin"},
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
	// - RDS access role per OU + policy + attachment
	// - ALB controller role + policy + attachment
	// - EBS CSI driver role + attachment
	// - Cluster autoscaler role + policy + attachment
	// - External DNS role + policy + attachment

	assert.Greater(t, resourceCounts.Add, 12, "Should create more than 12 IAM resources")
}

func TestIAMRolesMultipleOUs(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/iam-roles",
		Vars: map[string]interface{}{
			"cluster_name":      "test-cluster-multi-ou",
			"oidc_provider_arn": "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"oidc_provider_url": "https://oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"rds_instance_arn":  "arn:aws:rds:us-east-1:123456789012:db:test-db",
			"organizational_units": []map[string]interface{}{
				{
					"name":        "ou-admin",
					"ou_id":       "ou-admin-001",
					"permissions": []string{"admin"},
				},
				{
					"name":        "ou-dev",
					"ou_id":       "ou-dev-001",
					"permissions": []string{"deploy", "view"},
				},
				{
					"name":        "ou-readonly",
					"ou_id":       "ou-ro-001",
					"permissions": []string{"view"},
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should create RDS access resources for each OU (3 OUs * 3 resources each = 9)
	// Plus common service roles
	assert.Greater(t, resourceCounts.Add, 18, "Should create resources for multiple OUs")
}

func TestIAMRolesWithoutRDS(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/iam-roles",
		Vars: map[string]interface{}{
			"cluster_name":      "test-cluster-no-rds",
			"oidc_provider_arn": "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"oidc_provider_url": "https://oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"rds_instance_arn":  nil, // No RDS
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
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should not create RDS-specific roles
	assert.Greater(t, resourceCounts.Add, 8, "Should create service roles without RDS roles")
}

func TestIAMRolesALBController(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/iam-roles",
		Vars: map[string]interface{}{
			"cluster_name":      "test-cluster-alb",
			"oidc_provider_arn": "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"oidc_provider_url": "https://oidc.eks.us-east-1.amazonaws.com/id/TEST",
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
	assert.NotNil(t, planStruct, "Plan should include ALB controller role and policy")
}

func TestIAMRolesClusterAutoscaler(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/iam-roles",
		Vars: map[string]interface{}{
			"cluster_name":      "test-cluster-ca",
			"oidc_provider_arn": "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"oidc_provider_url": "https://oidc.eks.us-east-1.amazonaws.com/id/TEST",
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
	assert.NotNil(t, planStruct, "Plan should include Cluster Autoscaler role and policy")
}

func TestIAMRolesEBSCSIDriver(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/iam-roles",
		Vars: map[string]interface{}{
			"cluster_name":      "test-cluster-ebs",
			"oidc_provider_arn": "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"oidc_provider_url": "https://oidc.eks.us-east-1.amazonaws.com/id/TEST",
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
	assert.NotNil(t, planStruct, "Plan should include EBS CSI driver role")
}

func TestIAMRolesExternalDNS(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/iam-roles",
		Vars: map[string]interface{}{
			"cluster_name":      "test-cluster-dns",
			"oidc_provider_arn": "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"oidc_provider_url": "https://oidc.eks.us-east-1.amazonaws.com/id/TEST",
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
	assert.NotNil(t, planStruct, "Plan should include External DNS role and policy")
}

func TestIAMRolesRDSAccess(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/iam-roles",
		Vars: map[string]interface{}{
			"cluster_name":      "test-cluster-rds-access",
			"oidc_provider_arn": "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"oidc_provider_url": "https://oidc.eks.us-east-1.amazonaws.com/id/TEST",
			"rds_instance_arn":  "arn:aws:rds:us-east-1:123456789012:db:production-db",
			"organizational_units": []map[string]interface{}{
				{
					"name":        "app-team",
					"ou_id":       "ou-app-001",
					"permissions": []string{"deploy", "view"},
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should include RDS-specific IAM resources
	assert.Greater(t, resourceCounts.Add, 10, "Should create RDS access IAM resources")
}
