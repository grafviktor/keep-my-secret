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
	@$(MAKE) test

## test: run unit tests
.PHONY: test
test:
	@echo 'Running unit tests'
	go test -race -vet=off -count=1 -coverprofile unit.txt -covermode atomic ./...

## build-server: create binary with debugging symbols in /cmd/kms folder
.PHONY: build-server
build-server:
	# https://github.com/mattn/go-sqlite3/issues/384
	# https://words.filippo.io/easy-windows-and-linux-cross-compilers-for-macos/
	# @-rm -r ./build
	@echo 'Creating debug build'
	GOOS=linux GOARCH=amd64 go build $(LD_FLAGS) -o ./build/kms-linux-amd64 ./cmd/kms/*.go
	GOOS=darwin GOARCH=amd64 go build $(LD_FLAGS) -o ./build/kms-darwin-amd64 ./cmd/kms/*.go

## build-client: create client application
.PHONY: build-client
build-client:
	@echo 'Creating client build'
	@cd website && npm install && npm run build

## build: build the whole project
.PHONY: build
build:
	@echo 'Building the whole project'
	@$(MAKE) build-server
	@$(MAKE) build-client

## http-tls-key: create self-signed certificate and store it in /tls folder
.PHONY: http-tls-key
http-tls-key:
	@echo 'Creating self-signed HTTP TLS certificate'
	@cd tls && go run $$GOPATH/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

## run: start application
.PHONY: run
run:
	@echo 'Attempting to start the app'
	DEV=true go run $(LD_FLAGS) ./cmd/kms/*.go

## lint: run linter
.PHONY: lint
lint:
	@echo 'Running linter'
	@golangci-lint run
