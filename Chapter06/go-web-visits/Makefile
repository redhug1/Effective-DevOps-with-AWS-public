# on AWS
# GOOS: linux
# GOARCH: amd64

MAIN=visits-counter
BUILD=build
GOOS=linux
GOARCH=amd64
BUILD_ARCH=$(BUILD)
BIN_DIR?=.

export GOOS?=$(GOOS)
export GOARCH?=$(GOARCH)

.PHONY: build-AWS
build-aws:
		@mkdir -p $(BUILD_ARCH)/$(BIN_DIR)
		go build -o $(BUILD_ARCH)/$(BIN_DIR)/$(MAIN) $(MAIN).go
