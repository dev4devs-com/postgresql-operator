APP_NAME = postgresql-operator
ORG_NAME = dev4devs-com
PKG = github.com/$(ORG_NAME)/$(APP_NAME)
TOP_SRC_DIRS = pkg
PACKAGES ?= $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*_test.go \
              -exec dirname {} \\; | sort | uniq")
TEST_PKGS = $(addprefix $(PKG)/,$(PACKAGES))
APP_FILE=./cmd/manager/main.go
BIN_DIR := $(GOPATH)/bin
BINARY ?= postgresql-operator
IMAGE_REGISTRY=quay.io
IMAGE_LATEST_TAG=$(IMAGE_REGISTRY)/$(ORG_NAME)/$(APP_NAME):latest
IMAGE_MASTER_TAG=$(IMAGE_REGISTRY)/$(ORG_NAME)/$(APP_NAME):master
IMAGE_RELEASE_TAG=$(IMAGE_REGISTRY)/$(ORG_NAME)/$(APP_NAME):$(CIRCLE_TAG)
NAMESPACE=postgresql-operator
TEST_COMPILE_OUTPUT ?= build/_output/bin/$(APP_NAME)-test

# This follows the output format for goreleaser
BINARY_LINUX_64 = ./dist/linux_amd64/$(BINARY)

LDFLAGS=-ldflags "-w -s -X main.Version=${TAG}"
.DEFAULT_GOAL:=help

export GOPROXY?=https://proxy.golang.org/

##############################
# INSTALL-UNINSTALL          #
##############################

##@ Application

.PHONY: install
install:  ## Install all resources (CR-CRD's, RBCA and Operator)
	@echo ....... Creating namespace ....... 
	- kubectl create namespace ${NAMESPACE}
	@echo ....... Applying CRDS and Operator .......
	- kubectl apply -f deploy/crds/postgresql.dev4devs.com_databases_crd.yaml -n ${NAMESPACE}
	- kubectl apply -f deploy/crds/postgresql.dev4devs.com_backups_crd.yaml -n ${NAMESPACE}
	@echo ....... Applying Rules and Service Account .......
	- kubectl apply -f deploy/role.yaml -n ${NAMESPACE}
	- kubectl apply -f deploy/role_binding.yaml  -n ${NAMESPACE}
	- kubectl apply -f deploy/service_account.yaml  -n ${NAMESPACE}
	@echo ....... Applying Database Operator .......
	- kubectl apply -f deploy/operator.yaml -n ${NAMESPACE}
	@echo ....... Creating the Database .......
	- kubectl apply -f deploy/crds/postgresql.dev4devs.com_v1alpha1_database_cr.yaml -n ${NAMESPACE}

.PHONY: uninstall
uninstall:  ## Uninstall all that all performed in the $ make install
	@echo ....... Uninstalling .......
	@echo ....... Deleting CRDs.......
	- kubectl delete -f deploy/crds/postgresql.dev4devs.com_backups_crd.yaml -n ${NAMESPACE}
	- kubectl delete -f deploy/crds/postgresql.dev4devs.com_databases_crd.yaml -n ${NAMESPACE}
	@echo ....... Deleting Rules and Service Account .......
	- kubectl delete -f deploy/role.yaml -n ${NAMESPACE}
	- kubectl delete -f deploy/role_binding.yaml -n ${NAMESPACE}
	- kubectl delete -f deploy/service_account.yaml -n ${NAMESPACE}
	@echo ....... Deleting Operator .......
	- kubectl delete -f deploy/operator.yaml -n ${NAMESPACE}
	@echo ....... Deleting namespace ${NAMESPACE}.......
	- kubectl delete namespace ${NAMESPACE}

.PHONY: install-backup
install-backup: ## Install backup feature ( Backup CR )
	@echo Installing backup service in ${NAMESPACE} :
	- kubectl apply -f deploy/crds/postgresql.dev4devs.com_v1alpha1_backup_cr.yaml -n ${NAMESPACE}

.PHONY: uninstall-backup
uninstall-backup: ## Uninstall backup feature ( Backup CR )
	@echo Uninstalling backup service from ${NAMESPACE} :
	- kubectl delete -f deploy/crds/postgresql.dev4devs.com_v1alpha1_backup_cr.yaml -n ${NAMESPACE}

##############################
# CI                         #
##############################

##@ CI

.PHONY: code-build-linux
code-build-linux:  ## Build binary for Linux SO (amd64)
	env GOOS=linux GOARCH=amd64 go build $(APP_FILE)

.PHONY: image-build-master
image-build-master:  ## Build master branch image
	@echo Building operator with the tag: $(IMAGE_MASTER_TAG)
	operator-sdk build $(IMAGE_MASTER_TAG)

.PHONY: image-build-release
image-build-release: ## Build release and latest tag image
	@echo Building operator with the tag: $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_LATEST_TAG)

.PHONY: image-push-master
image-push-master: ## Push master branch image
	@echo Pushing operator with tag $(IMAGE_MASTER_TAG) to $(IMAGE_REGISTRY)
	@docker login --username $(QUAY_USERNAME) --password $(QUAY_PASSWORD) quay.io
	docker push $(IMAGE_MASTER_TAG)

.PHONY: image-push-release
image-push-release: ## Push release and latest tag image
	@echo Pushing operator with tag $(IMAGE_RELEASE_TAG) to $(IMAGE_REGISTRY)
	@docker login --username $(QUAY_USERNAME) --password $(QUAY_PASSWORD) quay.io
	docker push $(IMAGE_RELEASE_TAG)
	@echo Pushing operator with tag $(IMAGE_LATEST_TAG) to $(IMAGE_REGISTRY)
	docker push $(IMAGE_LATEST_TAG)


##############################
# Local Development          #
##############################

##@ Development

.PHONY: setup-debug
setup-debug:  ## Setup local env to debug. It will export env vars and install the project in the cluster
	@echo Exporting env vars to run operator locally:
	- . .-scripts-export_local_envvars.sh
	@echo Installing ...
	- make install

.PHONY: setup
setup:
	go mod tidy

.PHONY: run-local
run-local:  ## Run project locally for debbug purposes.
	@echo Exporting env vars to run operator locally:
	- . .-scripts-export_local_envvars.sh
	@echo  ....... Installing ...
	- make install
	@echo Starting ...
	- operator-sdk up local

.PHONY: vet
vet:  ## Run go vet for the project
	@echo go vet
	go vet $$(go list ./... )

.PHONY: fmt
fmt: ## Run go fmt for the project
	@echo go fmt
	go fmt $$(go list ./... )

.PHONY: dev
dev: ## Run all required dev commands. (It should be used always before send a PR)
	- make fmt
	- make vet
	- make gen

.PHONY: gen
gen:  ## Run SDK commands to generated-upddate the project
	operator-sdk generate k8s
	operator-sdk generate openapi

##############################
# Tests                      #
##############################

##@ Tests

.PHONY: test
test:  ## Run unit test
	@echo Running tests:
	go test -cover $(TEST_PKGS)

.PHONY: integration-cover
integration-cover:  ## Run coveralls
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -failfast -tags=integration -coverprofile=coverage.out -covermode=count $(addprefix $(PKG)/,$(pkg)) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: compile-e2e
compile-e2e:  ## Compile binary to run integration tests
	 @GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -c -o=$(TEST_COMPILE_OUTPUT) .-test-e2e ...

.PHONY: test-e2e
test-e2e:  ## Run integration tests locally
	  - kubectl create namespace ${NAMESPACE}
	  operator-sdk test local .-test-e2e --up-local --namespace=${NAMESPACE}

##@ General

.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)