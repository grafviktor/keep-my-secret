LD_FLAGS = -ldflags="-X main.buildVersion=v0.9.9 -X main.buildDate=$(shell date +%Y-%m-%d) -X main.buildCommit=$(shell git rev-parse --short=8 HEAD)"

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	gofumpt -l -w ./..
	goimports -w -local github.com/grafviktor/keep-my-secret .
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Linting code...'
	golangci-lint run
	@#$(MAKE) mylint
	@$(MAKE) test

## test: run unit tests
.PHONY: test
test:
	@echo 'Running unit tests'
	go test -race -vet=off -count=1 -coverprofile unit.txt -covermode atomic ./...

## build: creates binary with debugging symbols in /cmd/kms folder
.PHONY: build
build:
	@echo 'Creating debug build'
	go build $(LD_FLAGS) -o ./cmd/kms/kms ./cmd/kms/*.go

## build-client: creates client application
build-client:
	@echo 'Creating client build'
	@cd website && npm install && npm run build

## run: start application
.PHONY: run
run:
	@echo 'Attempting to start the app'
	go run $(LD_FLAGS) ./cmd/kms/*.go

## lint: run linter
.PHONY: lint
lint:
	@echo 'Running linter'
	@golangci-lint run
