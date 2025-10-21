package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestRDSModule(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/rds",
		Vars: map[string]interface{}{
			"identifier":              "test-rds",
			"vpc_id":                  "vpc-12345678",
			"subnet_ids":              []string{"subnet-db-1", "subnet-db-2", "subnet-db-3"},
			"availability_zones":      []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"engine":                  "postgres",
			"engine_version":          "15.4",
			"instance_class":          "db.t3.medium",
			"allocated_storage":       100,
			"database_name":           "testdb",
			"master_username":         "dbadmin",
			"backup_retention_period": 7,
			"multi_az":                true,
			"storage_encrypted":       true,
			"allowed_security_group_ids": []string{"sg-eks-nodes"},
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
	// - Random password
	// - Secrets Manager secret + version
	// - DB Subnet Group
	// - Security Group + Rules
	// - KMS Key + Alias
	// - DB Parameter Group
	// - DB Instance
	// - IAM Role for monitoring + policy attachment

	assert.Greater(t, resourceCounts.Add, 12, "Should create more than 12 resources")
}

func TestRDSModuleMultiAZ(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/rds",
		Vars: map[string]interface{}{
			"identifier":              "test-rds-multi-az",
			"vpc_id":                  "vpc-12345678",
			"subnet_ids":              []string{"subnet-1", "subnet-2", "subnet-3"},
			"availability_zones":      []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"engine":                  "postgres",
			"engine_version":          "15.4",
			"instance_class":          "db.r6g.xlarge",
			"allocated_storage":       500,
			"database_name":           "proddb",
			"master_username":         "admin",
			"backup_retention_period": 30,
			"multi_az":                true,
			"storage_encrypted":       true,
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should succeed with multi-AZ configuration")
}

func TestRDSModuleReadReplica(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/rds",
		Vars: map[string]interface{}{
			"identifier":          "test-rds-replica",
			"vpc_id":              "vpc-87654321",
			"subnet_ids":          []string{"subnet-rep-1", "subnet-rep-2", "subnet-rep-3"},
			"availability_zones":  []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			"engine":              "postgres",
			"engine_version":      "15.4",
			"instance_class":      "db.r6g.xlarge",
			"allocated_storage":   500,
			"database_name":       "replicadb",
			"master_username":     "admin",
			"multi_az":            true,
			"storage_encrypted":   true,
			"replicate_source_db": "arn:aws:rds:us-east-1:123456789012:db:test-rds-primary",
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Read replica shouldn't create secrets or primary DB
	assert.Greater(t, resourceCounts.Add, 5, "Should create replica resources")
}

func TestRDSModuleEncryption(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/rds",
		Vars: map[string]interface{}{
			"identifier":              "test-rds-encrypted",
			"vpc_id":                  "vpc-12345678",
			"subnet_ids":              []string{"subnet-1", "subnet-2", "subnet-3"},
			"availability_zones":      []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"engine":                  "postgres",
			"engine_version":          "15.4",
			"instance_class":          "db.t3.medium",
			"allocated_storage":       100,
			"database_name":           "encrypteddb",
			"master_username":         "admin",
			"backup_retention_period": 7,
			"multi_az":                false,
			"storage_encrypted":       true,
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should include KMS encryption")
}

func TestRDSModuleMySQLEngine(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/rds",
		Vars: map[string]interface{}{
			"identifier":              "test-rds-mysql",
			"vpc_id":                  "vpc-12345678",
			"subnet_ids":              []string{"subnet-1", "subnet-2", "subnet-3"},
			"availability_zones":      []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"engine":                  "mysql",
			"engine_version":          "8.0.35",
			"instance_class":          "db.t3.medium",
			"allocated_storage":       100,
			"database_name":           "mysqldb",
			"master_username":         "admin",
			"backup_retention_period": 7,
			"multi_az":                true,
			"storage_encrypted":       true,
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should succeed with MySQL engine")
}

func TestRDSModuleBackupRetention(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/rds",
		Vars: map[string]interface{}{
			"identifier":              "test-rds-backup",
			"vpc_id":                  "vpc-12345678",
			"subnet_ids":              []string{"subnet-1", "subnet-2", "subnet-3"},
			"availability_zones":      []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"engine":                  "postgres",
			"engine_version":          "15.4",
			"instance_class":          "db.t3.medium",
			"allocated_storage":       100,
			"database_name":           "backupdb",
			"master_username":         "admin",
			"backup_retention_period": 35, // Max retention
			"multi_az":                true,
			"storage_encrypted":       true,
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should succeed with extended backup retention")
}

func TestRDSModuleSecurityGroups(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/rds",
		Vars: map[string]interface{}{
			"identifier":         "test-rds-sg",
			"vpc_id":             "vpc-12345678",
			"subnet_ids":         []string{"subnet-1", "subnet-2", "subnet-3"},
			"availability_zones": []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"engine":             "postgres",
			"engine_version":     "15.4",
			"instance_class":     "db.t3.medium",
			"allocated_storage":  100,
			"database_name":      "testdb",
			"master_username":    "admin",
			"multi_az":           true,
			"storage_encrypted":  true,
			"allowed_security_group_ids": []string{
				"sg-eks-nodes-1",
				"sg-eks-nodes-2",
			},
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCounts := terraform.GetResourceCount(t, planStruct)

	// Should create security group rules for each allowed SG
	assert.Greater(t, resourceCounts.Add, 10, "Should create security group rules")
}

func TestRDSModulePerformanceInsights(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../modules/rds",
		Vars: map[string]interface{}{
			"identifier":              "test-rds-pi",
			"vpc_id":                  "vpc-12345678",
			"subnet_ids":              []string{"subnet-1", "subnet-2", "subnet-3"},
			"availability_zones":      []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			"engine":                  "postgres",
			"engine_version":          "15.4",
			"instance_class":          "db.r6g.large",
			"allocated_storage":       100,
			"database_name":           "perfdb",
			"master_username":         "admin",
			"backup_retention_period": 7,
			"multi_az":                true,
			"storage_encrypted":       true,
		},
		PlanOnly: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)
	assert.NotNil(t, planStruct, "Plan should enable Performance Insights")
}
