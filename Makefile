APP_NAME = own-redis
APP_VERSION = 0.1.0
APP_DESCRIPTION = a simple redis with golang

SOURCE_DIR = $(shell pwd)/own-redis
BUILD_DIR = $(shell pwd)/bin

LISTEN ?= 0.0.0.0
PORT ?= 6379

pre_commands:
	@printf "=============================================\n"
	@printf "App Name: $(APP_NAME)\n"
	@printf "App Version: $(APP_VERSION)\n"
	@printf "App Description: $(APP_DESCRIPTION)\n"
	@printf "=============================================\n"
	@printf "\n"

deps: pre_commands
	@printf "===> Installing dependencies...\n"
	@go mod tidy
	@printf "===> Dependencies installed.\n"
	@printf "\n"

build: pre_commands deps
	@printf "===> Building...\n"
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SOURCE_DIR)
	@printf "===> Build complete.\n"
	@printf "\n"

test: pre_commands deps
	@echo "Running tests..."
	@sleep 3
	@clear
	@go test -v $(SOURCE_DIR)/...

dev: pre_commands deps
	@printf "===> Running app in development mode...\n"
	@sleep 3
	@clear
	@air -c .air.toml

start: pre_commands build
	@printf "===> Running app...\n"
	@sleep 3
	@clear
	@$(BUILD_DIR)/$(APP_NAME) -listen $(LISTEN):$(PORT)

clean: pre_commands
	@printf "===> Cleaning...\n"
	@rm -rf $(BUILD_DIR)
	@printf "===> Clean complete.\n"
	@printf "\n"

fmt: pre_commands
	@printf "===> Formatting...\n"
	@go fmt $(SOURCE_DIR)/...
	@printf "===> Format complete.\n"
	@printf "\n"

lint: pre_commands deps fmt
	@printf "===> Linting...\n"
	@golangci-lint run $(SOURCE_DIR)/...
	@printf "===> Lint complete.\n"
	@printf "\n"

.PHONY: build start test deps clean fmt lint