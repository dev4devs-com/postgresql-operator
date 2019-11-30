#!/usr/bin/env bash
set -e

source ./scripts/common.sh

fetch_go_linter

header_text "Running golangci-lint"
golangci-lint run --disable-all \
    --deadline 5m \
    --enable=nakedret \
    --enable=maligned \
    --enable=ineffassign \
    --enable=goconst \
    --enable=errcheck \
    --enable=dupl \
    --enable=golint \
    --enable=gocyclo \
    --enable=misspell \
    --enable=interfacer \
    --enable=misspell \
    --enable=varcheck \
    --enable=structcheck \
    --enable=unparam \

# todo: enable
# --enable=goimports
# --enable=lll \
# --enable=gosec \