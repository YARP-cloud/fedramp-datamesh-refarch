# FedRAMP High Event-Driven Data Mesh
# Main Makefile

# Variables
SHELL := /bin/bash
GO := go
GOFMT := gofmt
TERRAFORM := terraform
KUBECTL := kubectl
DOCKER := docker

# Go CLI Configuration
CLI_DIR := ./cli
CLI_BIN := dmesh
CLI_BUILD_DIR := $(CLI_DIR)/bin
CLI_MAIN := $(CLI_DIR)/main.go

# Terraform Configuration
TF_DIR := ./platform/infrastructure/terraform
TF_ENV ?= dev

# Kubernetes Configuration
K8S_DIR := ./platform/infrastructure/kubernetes
K8S_NAMESPACE := fedramp-data-mesh

# Default target
.PHONY: all
all: help

# Help menu
.PHONY: help
help:
	@echo "FedRAMP High Event-Driven Data Mesh"
	@echo ""
	@echo "Usage:"
	@echo "  make build           Compile the CLI tool and check dependencies"
	@echo "  make lint            Run linters on Go code"
	@echo "  make test            Run unit tests"
	@echo "  make integration-test Run integration tests"
	@echo "  make tf-init         Initialize Terraform for selected environment"
	@echo "  make tf-plan         Generate and show Terraform execution plan"
	@echo "  make tf-apply        Apply Terraform changes"
	@echo "  make tf-destroy      Destroy Terraform-managed infrastructure"
	@echo "  make k8s-deploy      Deploy Kubernetes components"
	@echo "  make k8s-delete      Delete Kubernetes components"
	@echo "  make docker-build    Build Docker images"
	@echo "  make clean           Remove build artifacts"
	@echo ""
	@echo "Environment:"
	@echo "  TF_ENV              Terraform environment (default: dev)"
	@echo "                      Supported values: dev, test, prod"
	@echo ""

# Build the CLI tool
.PHONY: build
build:
	@echo "Building CLI tool..."
	mkdir -p $(CLI_BUILD_DIR)
	cd $(CLI_DIR) && $(GO) build -o bin/$(CLI_BIN) main.go
	@echo "Build complete: $(CLI_BUILD_DIR)/$(CLI_BIN)"

# Run linters
.PHONY: lint
lint:
	@echo "Running linters..."
	cd $(CLI_DIR) && $(GO) vet ./...
	cd $(CLI_DIR) && $(GOFMT) -s -w .
	@echo "Linting complete"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	cd $(CLI_DIR) && $(GO) test -v ./...
	@echo "Tests complete"

# Run integration tests
.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	cd $(CLI_DIR) && $(GO) test -v -tags=integration ./...
	@echo "Integration tests complete"

# Terraform init
.PHONY: tf-init
tf-init:
	@echo "Initializing Terraform for $(TF_ENV) environment..."
	cd $(TF_DIR) && $(TERRAFORM) init -reconfigure -backend-config=environments/$(TF_ENV)/backend.tfvars
	@echo "Terraform initialization complete"

# Terraform plan
.PHONY: tf-plan
tf-plan:
	@echo "Planning Terraform changes for $(TF_ENV) environment..."
	cd $(TF_DIR) && $(TERRAFORM) plan -var-file=environments/$(TF_ENV)/terraform.tfvars
	@echo "Terraform plan complete"

# Terraform apply
.PHONY: tf-apply
tf-apply:
	@echo "Applying Terraform changes for $(TF_ENV) environment..."
	cd $(TF_DIR) && $(TERRAFORM) apply -var-file=environments/$(TF_ENV)/terraform.tfvars
	@echo "Terraform apply complete"

# Terraform destroy
.PHONY: tf-destroy
tf-destroy:
	@echo "Destroying Terraform-managed infrastructure for $(TF_ENV) environment..."
	cd $(TF_DIR) && $(TERRAFORM) destroy -var-file=environments/$(TF_ENV)/terraform.tfvars
	@echo "Terraform destroy complete"

# Deploy Kubernetes components
.PHONY: k8s-deploy
k8s-deploy:
	@echo "Deploying Kubernetes components..."
	$(KUBECTL) apply -f $(K8S_DIR)/namespace.yaml
	$(KUBECTL) apply -k $(K8S_DIR)/schema-registry
	$(KUBECTL) apply -k $(K8S_DIR)/kafka-connect
	$(KUBECTL) apply -k $(K8S_DIR)/monitoring
	@echo "Kubernetes deployment complete"

# Delete Kubernetes components
.PHONY: k8s-delete
k8s-delete:
	@echo "Deleting Kubernetes components..."
	$(KUBECTL) delete -k $(K8S_DIR)/monitoring
	$(KUBECTL) delete -k $(K8S_DIR)/kafka-connect
	$(KUBECTL) delete -k $(K8S_DIR)/schema-registry
	$(KUBECTL) delete -f $(K8S_DIR)/namespace.yaml
	@echo "Kubernetes deletion complete"

# Build Docker images
.PHONY: docker-build
docker-build:
	@echo "Building Docker images..."
	$(DOCKER) build -t fedramp-data-mesh/kafka-connect:latest -f domains/project-management/producers/project-state/Dockerfile domains/project-management/producers/project-state
	@echo "Docker build complete"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(CLI_BUILD_DIR)
	@echo "Clean complete"
