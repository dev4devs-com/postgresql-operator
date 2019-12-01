#!/usr/bin/env bash
set -e

source ./scripts/common.sh

operator_logs() {
    kubectl describe pods -n postgresql-operator
    kubectl logs deployment/postgresql-operator -n postgresql-operator
}

remove_operator() {
   header_text "Removing all"
   make uninstall
}

apply_ci_image() {
  header_text "Replace image"
  sed -i".bak" -e "s/master/ci/g" deploy/operator.yaml; rm -f deploy/operator.yaml.bak
}


test_operator() {
  header_text "Installing Operator"
  make install

# todo: make timeout work

#  header_text "Checking Database"
#  if ! timeout 2m bash -c -- 'until kubectl describe Database -n postgresql-operator | grep OK; do sleep 1; done';
#  then
#      error_text "Error to deploy Database"
#      operator_logs
#      exit 1
#  fi

  header_text "Installing Backup Service"
  make install-backup

#  header_text "Checking Backup Service"
#  if ! timeout 2m bash -c -- 'until kubectl get cronjob.batch/backup -n postgresql-operator | grep backup; do sleep 1; done';
#  then
#      error_text "Error to install Backup Service"
#      operator_logs
#      exit 1
#  fi
}

apply_ci_image
trap_add 'remove_operator' EXIT
test_operator
remove_operator
