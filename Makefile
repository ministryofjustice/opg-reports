SHELL = bash
#============ BUILD INFO ==============
BUILD_DIR = ./builds/
API_VERSION = v1
COMMIT = $(shell git rev-parse HEAD)
ORGANISATION = OPG
SEMVER ?= v0.0.1
TIMESTAMP = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
#======================================
pkg=github.com/ministryofjustice/opg-reports/pkg
LDFLAGS:=" -X '${pkg}/bi.ApiVersion=${API_VERSION}' -X '${pkg}/bi.Commit=${COMMIT}' -X '${pkg}/bi.Organisation=${ORGANISATION}' -X '${pkg}/bi.Semver=${SEMVER}' -X '${pkg}/bi.Timestamp=${TIMESTAMP}'"
#======================================
tick="✅"

## run all tests
tests:
	@go clean -testcache
	@clear
	@echo "============== tests =============="
	@env env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -cover -covermode=atomic ./... && echo "" && echo "passed ${tick}"
.PHONY: tests

## run go suite tests that match the file pattern
## usage:
## `make tests name=<pattern>`
test:
	@go clean -testcache
	@clear
	@echo "============== test: [$(name)] =============="
	@env CGO_ENABLED=1 GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)"
.PHONY: test


## Run the go code coverage tool
coverage:
	@rm -Rf ./code-coverage.out
	@clear
	@echo "============== coverage =============="
	@env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover ./...
	@go tool cover -html=code-coverage.out
.PHONY: coverage

## Output the open api spec
openapi:
	@env CGO_ENABLED=1 go run ./api/main.go openapi > openapi.yaml
.PHONY: openapi

## Output build info
buildinfo:
	@echo "============ BUILD INFO =============="
	@echo "API_VERSION:  ${API_VERSION}"
	@echo "COMMIT:       ${COMMIT}"
	@echo "ORGANISATION: ${ORGANISATION}"
	@echo "SEMVER:       ${SEMVER}"
	@echo "TIMESTAMP:    ${TIMESTAMP}"
	@echo "======================================"
.PHONY: buildinfo

#========= RUN =========

## Run the api from dev source
api:
	@cd ./servers/api && go run main.go
.PHONY: api

#========= GO BUILD STEPS =========
# Build all binaries
build: build/servers build/collectors
.PHONY: build

build/servers: build/servers/api
.PHONY: build/servers
## Build the api into build directory
build/servers/api: buildinfo
	@echo -n "[building] servers/api .................. "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/servers/api ./servers/api/main.go && echo "${tick}"
.PHONY: build/servers/api

## build all collectors
build/collectors: buildinfo build/collectors/awscosts
.PHONY: build/collectors

build/collectors/awscosts:
	@echo -n "[building] collectors/awscosts .......... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/collectors/api ./collectors/awscosts/main.go && echo "${tick}"

.PHONY: build/collectors/awscosts
