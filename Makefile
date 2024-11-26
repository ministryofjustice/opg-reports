SHELL = bash
#============ PARAMS ==============
BUCKET_PROFILE ?= shared-development-operator
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
	@echo "=== tests"
	@env env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -cover -covermode=atomic ./... && echo "" && echo "passed ${tick}"
.PHONY: tests

## run go suite tests that match the file pattern
## usage:
## `make tests name=<pattern>`
test:
	@go clean -testcache
	@clear
	@echo "=== test: [$(name)]"
	@env CGO_ENABLED=1 GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)" && echo "" && echo "passed ${tick}"
.PHONY: test


## Run the go code coverage tool
coverage:
	@rm -Rf ./code-coverage.out
	@clear
	@echo "=== coverage"
	@env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover ./...
	@go tool cover -html=code-coverage.out
.PHONY: coverage

## Output the open api spec
openapi:
	@env CGO_ENABLED=1 go run ./servers/api/main.go openapi > openapi.yaml
.PHONY: openapi

## Output build info
buildinfo:
	@echo "=== PARAMS"
	@echo "BUCKET_PROFILE: ${BUCKET_PROFILE}"
	@echo "=== BUILD INFO"
	@echo "API_VERSION:    ${API_VERSION}"
	@echo "COMMIT:         ${COMMIT}"
	@echo "MODE:           ${MODE}"
	@echo "ORGANISATION:   ${ORGANISATION}"
	@echo "SEMVER:         ${SEMVER}"
	@echo "TIMESTAMP:      ${TIMESTAMP}"
	@echo "==="
.PHONY: buildinfo

## Removes an existing build artifacts
clean:
	@rm -f ./code-coverage.out
	@rm -f ./openapi.yaml
	@rm -Rf ./builds
	@rm -Rf ./databases
	@rm -Rf ./servers/api/databases
	@rm -Rf ./servers/front/assets
	@rm -Rf ./collectors/awscosts/data
	@rm -Rf ./collectors/githubstandards/data
	@rm -Rf ./collectors/githubreleases/data
	@rm -Rf ./collectors/awsuptime/data
.PHONY: clean
#========= RUN =========

## Run the api from dev source
api:
	@cd ./servers/api && go run main.go
.PHONY: api

## Run the front from dev source
front:
	@cd ./servers/front && go run main.go
.PHONY: front

#========= IMPORT DATA =========
# Import all old data - order is important due to data gaps
import: build
	@cd ./builds/ && aws-vault exec ${BUCKET_PROFILE} -- ./bin/convertor --download=false
	@cd ./builds && ./bin/importer -type=github-standards -file=./converted-data/github_standards.json
	@cd ./builds && ./bin/importer -type=aws-uptime -file=./converted-data/aws_uptime.json
	@cd ./builds && ./bin/importer -type=aws-costs -file=./converted-data/aws_costs.json
.PHONY: import
#========= BUILD GO BINARIES =========
# Build all binaries
build: buildinfo build/collectors build/convertor build/importer build/servers
.PHONY: build

build/servers: build/servers/api build/servers/front
.PHONY: build/servers

## Build the api into build directory
build/servers/api:
	@echo -n "[building] servers/api .................. "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/api ./servers/api/main.go && echo "${tick}"
.PHONY: build/servers/api

## Build the api into build directory
build/servers/front:
	@echo -n "[building] servers/front ................ "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/front ./servers/front/main.go && echo "${tick}"
.PHONY: build/servers/front

## build the convertor tool
build/convertor:
	@echo -n "[building] convertor .................... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/convertor ./convertor/main.go && echo "${tick}"
.PHONY: build/convertor


## build the importer tool
build/importer:
	@echo -n "[building] importer ..................... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/importer ./importer/main.go && echo "${tick}"
.PHONY: build/importer

## build all collectors
build/collectors: build/collectors/awscosts build/collectors/awsuptime build/collectors/githubstandards build/collectors/githubreleases
.PHONY: build/collectors

## build the aws costs collector
build/collectors/awscosts:
	@echo -n "[building] collectors/awscosts .......... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/awscosts ./collectors/awscosts/main.go && echo "${tick}"
.PHONY: build/collectors/awscosts

## build the aws uptime collector
build/collectors/awsuptime:
	@echo -n "[building] collectors/awsuptime ......... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/awsuptime ./collectors/awsuptime/main.go && echo "${tick}"
.PHONY: build/collectors/awsuptime

## build the github standards collector
build/collectors/githubstandards:
	@echo -n "[building] collectors/githubstandards ... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/githubstandards ./collectors/githubstandards/main.go && echo "${tick}"
.PHONY: build/collectors/githubstandards

## build the github releases collector
build/collectors/githubreleases:
	@echo -n "[building] collectors/githubreleases .... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/bin/githubreleases ./collectors/githubreleases/main.go && echo "${tick}"
.PHONY: build/collectors/githubreleases

