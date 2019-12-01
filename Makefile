APP_FILE=./cmd/manager/main.go
BIN_DIR := $(GOPATH)/bin
BINARY ?= postgresql-operator
IMAGE_LATEST_TAG=dev4devscom/postgresql-operator:latest
IMAGE_MASTER_TAG=dev4devscom/postgresql-operator:master
IMAGE_CI_TAG=dev4devscom/postgresql-operator:ci
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
install:  ## Creates the `{namespace}` namespace, application CRDS, cluster role and service account. Installs the operator and DB
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
uninstall:  ## Uninstalls the operator and DB. Deletes the `{namespace}`` namespace, application CRDS, cluster role and service account. i.e. all configuration applied by `make install`
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
install-backup: ## Installs the backup Service in the operator's namespace
	@echo Installing backup service in ${NAMESPACE} :
	- kubectl apply -f deploy/crds/postgresql.dev4devs.com_v1alpha1_backup_cr.yaml -n ${NAMESPACE}

.PHONY: uninstall-backup
uninstall-backup: ## Uninstalls the backup Service from the operator's namespace.
	@echo Uninstalling backup service from ${NAMESPACE} :
	- kubectl delete -f deploy/crds/postgresql.dev4devs.com_v1alpha1_backup_cr.yaml -n ${NAMESPACE}

##############################
# CI                         #
##############################

##@ CI

.PHONY: code-build-linux
code-build-linux:  ## Build binary for Linux SO (amd64)
	env GOOS=linux GOARCH=amd64 go build $(APP_FILE)

.PHONY: image-build-ci
image-build-ci:  ## Used by CI to build operator image from pr source code branch and add `:ci` tag.
	@echo Building operator with the tag: $(IMAGE_CI_TAG)
	operator-sdk build $(IMAGE_CI_TAG)

.PHONY: image-build-master
image-build-master:  ## Used by CI to build operator image from `master` branch and add `:master` tag.
	@echo Building operator with the tag: $(IMAGE_MASTER_TAG)
	operator-sdk build $(IMAGE_MASTER_TAG)

.PHONY: image-build-release
image-build-release: ## Used by CI to build operator image for relase tags
	@echo Building operator with the tag: $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_LATEST_TAG)

.PHONY: image-push-master
image-push-master: ## Used by CI to push the `master` image to https://hub.docker.com/r/dev4devscom/postgresql-operator
	@echo Pushing operator with tag $(IMAGE_MASTER_TAG)
	docker push $(IMAGE_MASTER_TAG)

.PHONY: image-push-ci
image-push-ci: ## Used by CI to push the `ci` image to https://hub.docker.com/r/dev4devscom/postgresql-operator
	@echo Pushing operator with tag $(IMAGE_CI_TAG)
	docker push $(IMAGE_CI_TAG)

.PHONY: image-push-release
image-push-release: ## Used by CI to push the `release` and `latest` image to https://hub.docker.com/r/dev4devscom/postgresql-operator
	@echo Pushing operator with tag $(IMAGE_RELEASE_TAG)
	docker push $(IMAGE_RELEASE_TAG)
	@echo Pushing operator with tag $(IMAGE_LATEST_TAG)
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
run-local:  ## Runs the operator locally for development purposes.
	@echo Exporting env vars to run operator locally:
	- . .-scripts-export_local_envvars.sh
	@echo  ....... Installing ...
	- make install
	@echo Starting ...
	- operator-sdk up local

.PHONY: vet
vet:  ## Examines source code and reports suspicious constructs using https://golang.org/cmd/vet/[vet].
	@echo go vet
	go vet $$(go list ./... )

.PHONY: fmt
fmt: ## Formats code using https://golang.org/cmd/gofmt/[gofmt].
	@echo go fmt
	go fmt $$(go list ./... )

.PHONY: dev
dev: fmt vet gen lint ##  It will tun the dev commands to check, fix and generated/update the files. (It should be used always before send a PR)

.PHONY: gen
gen:  ## It will automatically generated/update the files by using the operator-sdk based on the CR status and spec definitions.
	operator-sdk generate k8s
	operator-sdk generate openapi

.PHONY: lint
lint: ## Run golangci-lint for the project
	./scripts/check-lint.sh

##############################
# Tests                      #
##############################

##@ Tests
.PHONY: test
test:  ## Run unit test
	@echo Running tests:
	go test -coverprofile=coverage.out -covermode=count -count=1 -short ./cmd/... ./pkg/...

.PHONY: compile-e2e
compile-e2e:  ## Compile binary to run integration tests
	 @GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -c -o=$(TEST_COMPILE_OUTPUT) .-test-e2e ...

.PHONY: test-e2e
test-e2e:  ## Run integration tests locally
	  - kubectl create namespace ${NAMESPACE}
	  operator-sdk test local ./test/e2e --up-local --namespace=${NAMESPACE}

##@ General

.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)