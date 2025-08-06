# Makefile for managing buf commands

# Default shell configuration
SHELL := /bin/bash
GO_BIN := $(shell pwd)/.bin
PATH := $(GO_BIN):$(PATH)
GOCI_LINT_VERSION ?= v1.64.5

# Phony targets
.PHONY: generate-proto lint-breaking lint-proto install-tools

# Generate Go code for proto files
generate-proto:
	buf generate --template buf.gen.yaml

# Check for breaking changes
lint-breaking:
	buf breaking --against 'https://github.com/sanjevscet/news-grpc.git#branch=master'

# Lint proto files
lint-proto:
	buf lint --config buf.yaml

# Install required tools
install-tools:
	@echo "Installing tools..."
	@mkdir -p $(GO_BIN)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_BIN) $(GOCI_LINT_VERSION)
# 	@go install tool
	@echo "Tools installed successfully."