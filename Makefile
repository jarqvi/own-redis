APP_NAME = own-redis
APP_VERSION = $(shell echo $$APP_VERSION)
APP_DESCRIPTION = a simple redis with golang

SOURCE_DIR = $(shell pwd)/src
BUILD_DIR = $(shell pwd)/bin

all: build

build:
	@printf "App Name: $(APP_NAME)\n"
	@printf "App Version: $(APP_VERSION)\n"
	@printf "App Description: $(APP_DESCRIPTION)\n"
	@printf "\n"

	@printf "===> Building...\n"
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SOURCE_DIR)
	@printf "===> Build complete\n"
	@printf "\n"

run: build
	@printf "===> Running...\n"
	@$(BUILD_DIR)/$(APP_NAME)

test: deps
	@echo "Running tests"
	@go test -v $(SOURCE_DIR)/...

deps:
	@printf "===> Installing dependencies...\n"
	@go mod tidy
	@printf "===> Dependencies installed\n"
	@printf "\n"

clean:
	@printf "===> Cleaning...\n"
	@rm -rf $(BUILD_DIR)
	@printf "===> Clean complete\n"
	@printf "\n"

fmt:
	@printf "===> Formatting...\n"
	@go fmt $(SOURCE_DIR)/...
	@printf "===> Format complete\n"
	@printf "\n"

lint: fmt deps
	@printf "===> Linting...\n"
	@golangci-lint run $(SOURCE_DIR)/...
	@printf "===> Lint complete\n"
	@printf "\n"

.PHONY: all build run test deps clean fmt lint