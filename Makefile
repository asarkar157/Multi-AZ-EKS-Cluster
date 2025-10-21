# Makefile for Multi-Region EKS Cluster Terraform Module

.PHONY: help init plan apply destroy test test-unit test-integration fmt validate clean

# Default target
.DEFAULT_GOAL := help

# Colors for terminal output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

## Help target
help: ## Show this help message
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

## Terraform Operations
init: ## Initialize Terraform
	@echo "${GREEN}Initializing Terraform...${RESET}"
	terraform init

plan: ## Run Terraform plan
	@echo "${GREEN}Running Terraform plan...${RESET}"
	terraform plan

apply: ## Apply Terraform changes
	@echo "${GREEN}Applying Terraform changes...${RESET}"
	terraform apply

destroy: ## Destroy Terraform resources
	@echo "${YELLOW}Destroying Terraform resources...${RESET}"
	terraform destroy

## Testing
test: test-setup test-unit test-integration ## Run all tests
	@echo "${GREEN}All tests completed!${RESET}"

test-setup: ## Set up test dependencies
	@echo "${GREEN}Setting up test dependencies...${RESET}"
	cd test && go mod download

test-unit: ## Run unit tests
	@echo "${GREEN}Running unit tests...${RESET}"
	cd test && go test -v -timeout 30m -parallel 5 \
		-run 'TestVPC|TestEKSCluster|TestEKSNodeGroups|TestRDS|TestIAMRoles'

test-integration: ## Run integration tests
	@echo "${GREEN}Running integration tests...${RESET}"
	cd test && go test -v -timeout 30m -parallel 3 \
		-run 'TestRegionalEKS|TestMultiRegionEKS'

test-vpc: ## Run VPC module tests only
	@echo "${GREEN}Running VPC tests...${RESET}"
	cd test && go test -v -timeout 10m -run TestVPC

test-eks: ## Run EKS module tests only
	@echo "${GREEN}Running EKS tests...${RESET}"
	cd test && go test -v -timeout 10m -run TestEKS

test-rds: ## Run RDS module tests only
	@echo "${GREEN}Running RDS tests...${RESET}"
	cd test && go test -v -timeout 10m -run TestRDS

test-iam: ## Run IAM roles module tests only
	@echo "${GREEN}Running IAM tests...${RESET}"
	cd test && go test -v -timeout 10m -run TestIAMRoles

test-coverage: ## Run tests with coverage report
	@echo "${GREEN}Running tests with coverage...${RESET}"
	cd test && go test -v -timeout 30m -cover -coverprofile=coverage.out
	cd test && go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}Coverage report generated: test/coverage.html${RESET}"

## Code Quality
fmt: ## Format Terraform code
	@echo "${GREEN}Formatting Terraform code...${RESET}"
	terraform fmt -recursive .

fmt-check: ## Check if Terraform code is formatted
	@echo "${GREEN}Checking Terraform formatting...${RESET}"
	terraform fmt -check -recursive .

validate: ## Validate Terraform configuration
	@echo "${GREEN}Validating Terraform configuration...${RESET}"
	terraform validate

lint: ## Run tflint
	@echo "${GREEN}Running tflint...${RESET}"
	tflint --init
	tflint --recursive

security: ## Run security scans
	@echo "${GREEN}Running security scans...${RESET}"
	@echo "${CYAN}Running tfsec...${RESET}"
	tfsec .
	@echo "${CYAN}Running checkov...${RESET}"
	checkov -d . --framework terraform

docs: ## Generate module documentation
	@echo "${GREEN}Generating documentation...${RESET}"
	terraform-docs markdown table --output-file README.md --output-mode inject .
	@for module in modules/*; do \
		if [ -d "$$module" ]; then \
			echo "Generating docs for $$module..."; \
			terraform-docs markdown table --output-file README.md --output-mode inject $$module; \
		fi \
	done

## Cleanup
clean: ## Clean up temporary files
	@echo "${YELLOW}Cleaning up temporary files...${RESET}"
	find . -type d -name ".terraform" -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "terraform.tfstate*" -delete
	find . -type f -name ".terraform.lock.hcl" -delete
	cd test && rm -f coverage.out coverage.html

clean-all: clean ## Clean up all generated files including test cache
	@echo "${YELLOW}Cleaning up all generated files...${RESET}"
	cd test && go clean -testcache
	find . -type f -name "*.backup" -delete

## Installation
install-tools: ## Install required tools
	@echo "${GREEN}Installing required tools...${RESET}"
	@echo "${CYAN}Installing Terraform...${RESET}"
	@which terraform > /dev/null || (echo "Please install Terraform manually" && exit 1)
	@echo "${CYAN}Installing Go...${RESET}"
	@which go > /dev/null || (echo "Please install Go manually" && exit 1)
	@echo "${CYAN}Installing tflint...${RESET}"
	@which tflint > /dev/null || brew install tflint || echo "Please install tflint manually"
	@echo "${CYAN}Installing tfsec...${RESET}"
	@which tfsec > /dev/null || brew install tfsec || echo "Please install tfsec manually"
	@echo "${CYAN}Installing terraform-docs...${RESET}"
	@which terraform-docs > /dev/null || brew install terraform-docs || echo "Please install terraform-docs manually"
	@echo "${GREEN}Tool installation complete!${RESET}"

## CI/CD
ci: fmt-check validate test ## Run CI pipeline locally
	@echo "${GREEN}CI pipeline completed successfully!${RESET}"

pre-commit: fmt validate lint security ## Run pre-commit checks
	@echo "${GREEN}Pre-commit checks passed!${RESET}"

## Development
dev-setup: install-tools test-setup ## Set up development environment
	@echo "${GREEN}Development environment ready!${RESET}"

example-basic: ## Show basic usage example
	@echo "${CYAN}Basic usage example:${RESET}"
	@cat terraform.tfvars.example

example-advanced: ## Show advanced usage example
	@echo "${CYAN}Advanced usage example - see README.md${RESET}"

## Module Specific
module-vpc: ## Validate VPC module
	@echo "${GREEN}Validating VPC module...${RESET}"
	cd modules/vpc && terraform init -backend=false && terraform validate

module-eks: ## Validate EKS module
	@echo "${GREEN}Validating EKS module...${RESET}"
	cd modules/eks-cluster && terraform init -backend=false && terraform validate

module-rds: ## Validate RDS module
	@echo "${GREEN}Validating RDS module...${RESET}"
	cd modules/rds && terraform init -backend=false && terraform validate

module-all: module-vpc module-eks module-rds ## Validate all modules
	@echo "${GREEN}All modules validated!${RESET}"

## Release
version: ## Show current version
	@echo "${CYAN}Current version: 1.0.0${RESET}"

changelog: ## Generate changelog
	@echo "${GREEN}Generating changelog...${RESET}"
	@git log --pretty=format:"- %s (%h)" --reverse > CHANGELOG.md
	@echo "${GREEN}Changelog generated: CHANGELOG.md${RESET}"
