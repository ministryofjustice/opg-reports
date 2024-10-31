SHELL = bash
#============ BUILD INFO ==============
BUILD_DIR = ./builds/
API_VERSION = v1
COMMIT = $(shell git rev-parse HEAD)
ORGANISATION = OPG
SEMVER ?= v0.0.1
MODE ?= simple
TIMESTAMP = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
#======================================
pkg=github.com/ministryofjustice/opg-reports/pkg
LDFLAGS:=" -X '${pkg}/bi.ApiVersion=${API_VERSION}' -X '${pkg}/bi.Commit=${COMMIT}' -X '${pkg}/bi.Organisation=${ORGANISATION}' -X '${pkg}/bi.Semver=${SEMVER}' -X '${pkg}/bi.Timestamp=${TIMESTAMP}' -X '${pkg}/bi.Mode=${MODE}' "
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
	@env CGO_ENABLED=1 go run ./servers/sapi/main.go openapi > openapi.yaml
.PHONY: openapi

## Output build info
buildinfo:
	@echo "============ BUILD INFO =============="
	@echo "API_VERSION:  ${API_VERSION}"
	@echo "COMMIT:       ${COMMIT}"
	@echo "MODE:         ${MODE}"
	@echo "ORGANISATION: ${ORGANISATION}"
	@echo "SEMVER:       ${SEMVER}"
	@echo "TIMESTAMP:    ${TIMESTAMP}"
	@echo "======================================"
.PHONY: buildinfo

## Removes an existing build artifacts
clean:
	@rm -Rf ./builds
	@rm -Rf ./servers/sapi/databases
	@rm -Rf ./servers/sfront/assets
	@rm -Rf ./databases
	@rm -Rf ./collectors/cawscosts/data
.PHONY: clean
#========= RUN =========

## Run the api from dev source
api:
	@cd ./servers/sapi && go run main.go
.PHONY: api

## Run the front from dev source
front:
	@cd ./servers/sfront && go run main.go
.PHONY: front

#========= BUILD GO BINARIES =========
# Build all binaries
build: buildinfo build/collectors build/importers build/servers
.PHONY: build

build/servers: build/servers/api build/servers/front
.PHONY: build/servers

## Build the api into build directory
build/servers/api:
	@echo -n "[building] servers/sapi .................. "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/sapi ./servers/sapi/main.go && echo "${tick}"
.PHONY: build/servers/api

## Build the api into build directory
build/servers/front:
	@echo -n "[building] servers/sfront ................ "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/sfront ./servers/sfront/main.go && echo "${tick}"
.PHONY: build/servers/front


## build all importers
build/importers: build/importers/isqlite
.PHONY: build/importers

## build the sqlite importer tool
build/importers/isqlite:
	@echo -n "[building] importers/isqlite ............. "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/isqlite ./importers/isqlite/main.go && echo "${tick}"
.PHONY: build/importers/sqlite

## build all collectors
build/collectors: build/collectors/cawscosts
.PHONY: build/collectors

## build the aws costs collector
build/collectors/cawscosts:
	@echo -n "[building] collectors/cawscosts .......... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/cawscosts ./collectors/cawscosts/main.go && echo "${tick}"
.PHONY: build/collectors/awscosts
