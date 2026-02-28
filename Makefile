project_name = cmyk
image_name = $(project_name)

socks5_proxy = socks5://127.0.0.1:8123

.PHONY: openapi

help: ## This help dialog.
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

clean-packages: ## Remove local packages
	@echo
	@echo "Removing local packages"
	go clean -modcache

clean: ## Deletes intermediate files
	@echo
	@echo "Deleting intermediate files..."
	rm -rf ./tmp

reinit: clean ## Reinitialize the app
	@echo
	@echo "Reinitializing the app..."
	rm -rf ./openapi
	rm -rf ./tmp
	rm -rf ./vendor
	rm -f go.sum
	rm -f go.mod
	go mod init $(project_name)
	go mod tidy
	go mod vendor

vendor: ## Vendor dependencies
	@echo
	@echo "Vendoring dependencies..."
	go mod vendor

deps: ## Tidy and pin dependencies
	@echo
	@echo "Tidying and pinning dependencies..."
	go mod tidy

check: ## Run static analysis checks
	@echo
	@echo "Running static analysis checks..."
	go fmt -mod=vendor ./cmd/... ./internal/...
	go vet -mod=vendor ./cmd/... ./internal/...
	revive -config revive.toml -formatter friendly ./cmd/... ./internal/...

test: ## Run unit tests
	@echo
	@echo "Running unit tests..."
	go test -timeout 5000ms -cover -mod=vendor ./cmd/... ./internal/...

openapi: ## Build API specification
	@echo
	@echo "Gnerating OpenAPI Specification..."
	rm -rf ./openapi
	swag init -o ./openapi -g main.go -d ./cmd/api,./internal/handlers,./internal/models

build: ## Build local app
	@echo
	@echo "Building local app"
	rm -rf ./tmp
	go build -ldflags "-X main.BuildDate=BUILD_DATE_NOT_SET -X main.BuildNumber=BUILD_NUMBER_NOT_SET -X main.CommitId=COMMIT_ID_NOT_SET" -mod=vendor -o ./tmp/$(project_name) ./cmd/api

run: ## Run the app locally
	@echo
	@echo "Running the app locally"
	go run cmd/api/main.go

runs: ## Run the app locally using a SOCKS5 proxy
	@echo
	@echo "Running the app locally"
	SOCKS5_PROXY=$(socks5_proxy) go run cmd/api/main.go

run_: ## Run the app locally in mock mode
	@echo
	@echo "Running the app locally"
	MOCK_MODE=1 go run cmd/api/main.go
