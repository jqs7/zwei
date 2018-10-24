.DEFAULT_GOAL := help

APP_NAME=zwei
APP_CMD_DIR=cmd/$(APP_NAME)
APP_BINARY=bin/$(APP_NAME)
APP_BINARY_UNIX=bin/$(APP_NAME)_unix_amd64

all: postgres build migration zwei

.PHONY: test
test: ## test
	go test -v ./...

.PHONY: zwei
zwei: ## run zwei in docker 
	docker-compose build zwei
	docker-compose up -d zwei

.PHONY: postgres
postgres: ## run postgres in docker 
	docker-compose up -d postgres

.PHONY: migration
migration: ## migration
	cd cmd/migrate && ./migrate

.PHONY: build
build: ## build
	go build -o $(APP_BINARY) -v $(APP_CMD_DIR)/main.go
	go build -o bin/migrate -v cmd/migrate/*.go

.PHONY: install
install: ## install
	go install $(APP_CMD_DIR)/main.go

.PHONY: clean
clean: ## clean 
	go clean
	rm -f $(APP_BINARY)
	rm -f $(APP_BINARY_UNIX)

.PHONY: run
run: build ## run
	./$(APP_BINARY)

.PHONY: build-linux
build-linux: ## build linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(APP_BINARY_UNIX) -v $(APP_CMD_DIR)/main.go


.PHONY: help
help: 
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
