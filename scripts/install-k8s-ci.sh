#!/usr/bin/env bash
set -e

source ./scripts/common.sh

header_text "Installing kind"

KIND_IMAGE="docker.io/kindest/node:${K8S_VERSION}"

# Download the latest version of kind, which supports all versions of
# Kubernetes v1.11+.
curl -Lo kind https://github.com/kubernetes-sigs/kind/releases/latest/download/kind-$(uname)-amd64
chmod +x kind
sudo mv kind /usr/local/bin/

header_text "Create a cluster of version ${K8S_VERSION}."
kind create cluster --image="$KIND_IMAGE"

header_text "Exporting kubeconfig"
kind export kubeconfig

header_text "Installing kubectl"
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl
chmod +x kubectl
sudo mv kubectl /usr/local/bin/

header_text "Setting namespace as default"
kubectl config set-context --current --namespace=postgresql-operator
kubectl cluster-info