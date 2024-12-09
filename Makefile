SHELL := /bin/bash
.DEFAULT_GOAL := build

BIN         = $(CURDIR)/bin
BUILD_DIR   = $(CURDIR)/build

GOPATH      = $(HOME)/go
GOBIN       = $(GOPATH)/bin
GO          ?= GOGC=off $(shell which go)
PKGS        = $(or $(PKG),$(shell env $(GO) list ./...))
VERSION     ?= $(shell git describe --tags --always --match=v*)
SHORT_COMMIT ?= $(shell git rev-parse --short HEAD)
ENVIRONMENT ?= local

PATH := $(BIN):$(GOBIN):$(PATH)

LDFLAGS = -w -s -X "github.com/algoboyz/garm/init.version=$(VERSION)" -X "github.com/algoboyz/garm/init.commit=$(SHORT_COMMIT)"

# Targets
NAME = $(shell basename $(CURDIR))
TARGET = $(BUILD_DIR)/$(NAME).o
OUT = $(BIN)/$(NAME)

# Printing
V ?= 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

$(BIN):
	@mkdir -p $@
$(BUILD_DIR):
	@mkdir -p $@

# Tools
$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(@F)…)
	$Q GOBIN=$(BIN) $(GO) install -v $(shell $(GO) list -tags=tools -e -f '{{ join .Imports "\n" }}' tools/tools.go | grep $(@F))

GOLANGCI_LINT = $(BIN)/golangci-lint
TOOLCHAIN = $(GOLANGCI_LINT)

# Targets
.PHONY: lint # Run project linters
lint: | $(GOLANGCI_LINT)
	$(info $(M) running linter…)
	$Q $(GOLANGCI_LINT) run --max-issues-per-linter 10 --timeout 5m

.PHONY: test # Run all tests and generate coverage report
test:
	$(info $(M) running tests locally…)
	$Q $(GO) test $(PKGS) -race -count 1 -timeout 5m -covermode=atomic -coverprofile cover.out fmt

.PHONY: build # Build service
build: $(BUILD_DIR)
	$(info $(M) $(GOOS) $(GOARCH) building executable…)
	$Q CGO_ENABLED=0 GOOS=$(or $(GOOS),linux) GOARCH=$(or $(GOARCH),amd64) $(GO) build -ldflags '$(LDFLAGS)' -a -o $(BUILD_DIR)/simd ./simd.go
	$Q chmod +x $(BUILD_DIR)/simd
	@true

.PHONY: build-mac # Build mac executable
mac: $(BIN) $(BUILD_DIR)
	$(info $(M) building mac executable…)
	$Q as $(NAME).s -o $(TARGET)
	$Q ld $(TARGET) -o $(OUT) -l System -syslibroot `xcrun -sdk macosx --show-sdk-path` -e _start -arch arm64
	@true

.PHONY: linux # Build linux executable
linux: $(BIN) $(BUILD_DIR)
	$(info $(M) building linux executable…)
	$Q set -o nounset
	$Q set -o errexit
	$Q gcc -o $(TARGET) $(NAME).s
	@true

.PHONY: fmt # Run gofmt on go source files
fmt:
	$(info $(M) running fmt…)
	$Q

.PHONY: generate # Run go generate on go source files
generate: | $(TOOLCHAIN)
	$(info $(M) running go generate…)
	$Q $(GO) generate $(PKGS)
	$Q $(MAKE) fmt

.PHONY: clean # Cleanup project root
clean:
	$(info $(M) cleaning…)
	@rm -rf $(BIN)
	@rm -rf $(BUILD_DIR)

.PHONY: help # Display help
help:
	@grep  -E '^.PHONY' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "#|: "}; {printf "\033[36m%-20s\033[0m %s\n",$$2,$$3}'