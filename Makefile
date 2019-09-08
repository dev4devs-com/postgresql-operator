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

##############################
# INSTALL/UNINSTALL          #
##############################

.PHONY: install
install:
	@echo ....... Creating namespace ....... 
	- kubectl create namespace ${NAMESPACE}
	@echo ....... Applying CRDS and Operator .......
	- kubectl apply -f deploy/crds/postgresql_v1alpha1_database_crd.yaml -n ${NAMESPACE}
	- kubectl apply -f deploy/crds/postgresql_v1alpha1_backup_crd.yaml -n ${NAMESPACE}
	@echo ....... Applying Rules and Service Account .......
	- kubectl apply -f deploy/role.yaml -n ${NAMESPACE}
	- kubectl apply -f deploy/role_binding.yaml  -n ${NAMESPACE}
	- kubectl apply -f deploy/service_account.yaml  -n ${NAMESPACE}
	@echo ....... Applying Database Operator .......
	- kubectl apply -f deploy/operator.yaml -n ${NAMESPACE}
	@echo ....... Creating the Database .......
	- kubectl apply -f deploy/crds/postgresql_v1alpha1_database_cr.yaml -n ${NAMESPACE}

.PHONY: uninstall
uninstall:
	@echo ....... Uninstalling .......
	@echo ....... Deleting CRDs.......
	- kubectl delete -f deploy/crds/postgresql_v1alpha1_backup_crd.yaml -n ${NAMESPACE}
	- kubectl delete -f deploy/crds/postgresql_v1alpha1_database_crd.yaml -n ${NAMESPACE}
	@echo ....... Deleting Rules and Service Account .......
	- kubectl delete -f deploy/role.yaml -n ${NAMESPACE}
	- kubectl delete -f deploy/role_binding.yaml -n ${NAMESPACE}
	- kubectl delete -f deploy/service_account.yaml -n ${NAMESPACE}
	@echo ....... Deleting Operator .......
	- kubectl delete -f deploy/operator.yaml -n ${NAMESPACE}
	@echo ....... Deleting namespace ${NAMESPACE}.......
	- kubectl delete namespace ${NAMESPACE}

.PHONY: backup/install
backup/install:
	@echo Installing backup service in ${NAMESPACE} :
	- kubectl apply -f deploy/crds/postgresql_v1alpha1_backup_cr.yaml -n ${NAMESPACE}

.PHONY: backup/uninstall
backup/uninstall:
	@echo Uninstalling backup service from ${NAMESPACE} :
	- kubectl delete -f deploy/crds/postgresql_v1alpha1_backup_cr.yaml -n ${NAMESPACE}

##############################
# CI                         #
##############################

.PHONY: code/build/linux
code/build/linux:
	env GOOS=linux GOARCH=amd64 go build $(APP_FILE)

.PHONY: image/build/master
image/build/master:
	@echo Building operator with the tag: $(IMAGE_MASTER_TAG)
	operator-sdk build $(IMAGE_MASTER_TAG)

.PHONY: image/build/release
image/build/release:
	@echo Building operator with the tag: $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_RELEASE_TAG)
	operator-sdk build $(IMAGE_LATEST_TAG)

.PHONY: image/push/master
image/push/master:
	@echo Pushing operator with tag $(IMAGE_MASTER_TAG) to $(IMAGE_REGISTRY)
	@docker login --username $(QUAY_USERNAME) --password $(QUAY_PASSWORD) quay.io
	docker push $(IMAGE_MASTER_TAG)

.PHONY: image/push/release
image/push/release:
	@echo Pushing operator with tag $(IMAGE_RELEASE_TAG) to $(IMAGE_REGISTRY)
	@docker login --username $(QUAY_USERNAME) --password $(QUAY_PASSWORD) quay.io
	docker push $(IMAGE_RELEASE_TAG)
	@echo Pushing operator with tag $(IMAGE_LATEST_TAG) to $(IMAGE_REGISTRY)
	docker push $(IMAGE_LATEST_TAG)


##############################
# Local Development          #
##############################

.PHONY: setup/debug
setup/debug:
	@echo Exporting env vars to run operator locally:
	- . ./scripts/export_local_envvars.sh
	@echo Installing ...
	- make install

.PHONY: setup
setup:
	go mod tidy

.PHONY: code/run/local
code/run/local:
	@echo Exporting env vars to run operator locally:
	- . ./scripts/export_local_envvars.sh
	@echo  ....... Installing ...
	- make install
	@echo Starting ...
	- operator-sdk up local

.PHONY: code/vet
code/vet:
	@echo go vet
	go vet $$(go list ./... )

.PHONY: code/fmt
code/fmt:
	@echo go fmt
	go fmt $$(go list ./... )

.PHONY: code/dev
code/dev:
	- make code/fmt
	- make code/vet
	- make code/gen

.PHONY: code/gen
code/gen:
	operator-sdk generate k8s
	operator-sdk generate openapi

##############################
# Tests                      #
##############################

.PHONY: test/run
test/run:
	@echo Running tests:
	go test -cover $(TEST_PKGS)

.PHONY: test/integration-cover
test/integration-cover:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -failfast -tags=integration -coverprofile=coverage.out -covermode=count $(addprefix $(PKG)/,$(pkg)) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: test/compile/e2e
test/compile/e2e:
	 @GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -c -o=$(TEST_COMPILE_OUTPUT) ./test/e2e ...

.PHONY: test/e2e
test/e2e:
	  - kubectl create namespace ${NAMESPACE}
	  operator-sdk test local ./test/e2e --up-local --namespace=${NAMESPACE}
