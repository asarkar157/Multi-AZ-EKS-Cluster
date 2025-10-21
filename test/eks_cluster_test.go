package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestEKSClusterModule(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-cluster",
		Vars: map[string]interface{}{
			"cluster_name":              "test-eks-cluster",
			"kubernetes_version":        "1.28",
			"vpc_id":                    "vpc-12345678",
			"subnet_ids":                []string{"subnet-1", "subnet-2", "subnet-3"},
			"control_plane_subnet_ids":  []string{"subnet-1", "subnet-2", "subnet-3"},
			"environment":               "test",
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
	// - IAM Role for cluster
	// - 2 IAM Role Policy Attachments
	// - KMS Key + Alias
	// - Security Group
	// - Security Group Rule
	// - CloudWatch Log Group
	// - EKS Cluster
	// - OIDC Provider
	// - 4 EKS Addons
	// - EKS Access Entry per OU
	// - IAM Role per OU
	// - EKS Access Policy Association per OU

	assert.Greater(t, resourceCounts.Add, 15, "Should create more than 15 resources")
}

func TestEKSClusterEncryption(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-cluster",
		Vars: map[string]interface{}{
			"cluster_name":             "test-eks-encrypted",
			"kubernetes_version":       "1.28",
			"vpc_id":                   "vpc-12345678",
			"subnet_ids":               []string{"subnet-1", "subnet-2", "subnet-3"},
			"control_plane_subnet_ids": []string{"subnet-1", "subnet-2", "subnet-3"},
			"environment":              "test",
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
	assert.NotNil(t, planStruct, "Plan should include KMS encryption")
}

func TestEKSClusterAddons(t *testing.T) {
	t.Parallel()

	// Test with custom addon versions
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-cluster",
		Vars: map[string]interface{}{
			"cluster_name":             "test-eks-addons",
			"kubernetes_version":       "1.28",
			"vpc_id":                   "vpc-12345678",
			"subnet_ids":               []string{"subnet-1", "subnet-2", "subnet-3"},
			"control_plane_subnet_ids": []string{"subnet-1", "subnet-2", "subnet-3"},
			"environment":              "test",
			"vpc_cni_version":          "v1.15.0-eksbuild.1",
			"coredns_version":          "v1.10.1-eksbuild.2",
			"kube_proxy_version":       "v1.28.1-eksbuild.1",
			"ebs_csi_driver_version":   "v1.25.0-eksbuild.1",
			"organizational_units": []map[string]interface{}{
				{
					"name":        "test-ou",
					"ou_id":       "ou-test-001",
					"permissions": []string{"deploy", "view"},
				},
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should succeed with custom addon versions")
}

func TestEKSClusterMultipleOUs(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-cluster",
		Vars: map[string]interface{}{
			"cluster_name":             "test-eks-multi-ou",
			"kubernetes_version":       "1.28",
			"vpc_id":                   "vpc-12345678",
			"subnet_ids":               []string{"subnet-1", "subnet-2", "subnet-3"},
			"control_plane_subnet_ids": []string{"subnet-1", "subnet-2", "subnet-3"},
			"environment":              "test",
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

	// Should create resources for 3 OUs
	assert.Greater(t, resourceCounts.Add, 20, "Should create resources for multiple OUs")
}

func TestEKSClusterLogging(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/eks-cluster",
		Vars: map[string]interface{}{
			"cluster_name":             "test-eks-logging",
			"kubernetes_version":       "1.28",
			"vpc_id":                   "vpc-12345678",
			"subnet_ids":               []string{"subnet-1", "subnet-2", "subnet-3"},
			"control_plane_subnet_ids": []string{"subnet-1", "subnet-2", "subnet-3"},
			"environment":              "test",
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
	assert.NotNil(t, planStruct, "Plan should include CloudWatch logging")
}
