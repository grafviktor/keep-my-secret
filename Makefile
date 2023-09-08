LD_FLAGS = -ldflags="-X main.buildVersion=v0.9.9 -X main.buildDate=$(shell date +%Y-%m-%d) -X main.buildCommit=$(shell git rev-parse --short=8 HEAD)"
LD_FLAGS_WIN_LIN = -ldflags="-X main.buildVersion=v0.9.9 -X main.buildDate=$(shell date +%Y-%m-%d) -X main.buildCommit=$(shell git rev-parse --short=8 HEAD) -linkmode external -extldflags -static"

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
	# https://www.andrewheiss.com/blog/2020/01/10/makefile-subdirectory-zips/
	# https://github.com/mattn/go-sqlite3/issues/384
	# @-rm -r ./build
	@echo 'Creating debug build'
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build $(LD_FLAGS) -o ./build/kms-darwin-amd64 ./cmd/kms/*.go
	# https://words.filippo.io/easy-windows-and-linux-cross-compilers-for-macos/
	# brew install FiloSottile/musl-cross/musl-cross
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build $(LD_FLAGS_WIN_LIN) -o ./build/kms-linux-amd64 ./cmd/kms/*.go
	# https://words.filippo.io/easy-windows-and-linux-cross-compilers-for-macos/
	# 1. brew install FiloSottile/musl-cross/musl-cross + 2. brew install mingw-w64
	CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOARCH=amd64 GOOS=windows CGO_ENABLED=1 go build $(LD_FLAGS_WIN_LIN) -o ./build/kms-windows-amd64.exe ./cmd/kms/*.go

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

## swagger: runs swag utility which generates API documentation
# 'swag' should be installed: go install github.com/swaggo/swag/cmd/swag@latest
# you should only run it from a folder where "main.go" file is located
.PHONY: swagger
swagger:
	@echo 'Generating API docs in cmd/kms/swagger folder'
	@cd cmd/kms/ && swag init --output ./swagger/