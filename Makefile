SHELL = bash
#============ PARAMS ==============
BUCKET_PROFILE ?= shared-development-operator
IMPORT_DOWNLOAD ?= true
GITHUB_ACCESS_TOKEN ?= ${GITHUB_TOKEN}
#============ BUILD INFO ==============
BUILD_DIR = ./builds/
API_VERSION = v1

COMMIT = $(shell git rev-parse HEAD)
TIMESTAMP = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SEMVER ?= v0.0.1

ORGANISATION = OPG
DATASET ?= real
FIXTURES ?= full
BUCKET_NAME ?= report-data-development
#======================================
pkg=github.com/ministryofjustice/opg-reports/info
LDFLAGS:="-X '${pkg}.Commit=${COMMIT}' -X '${pkg}.Timestamp=${TIMESTAMP}' -X '${pkg}.Semver=${SEMVER}' -X '${pkg}.Organisation=${ORGANISATION}' -X '${pkg}.Dataset=${DATASET}' -X '${pkg}.Fixtures=${FIXTURES}' -X '${pkg}.BucketName=${BUCKET_NAME}'"
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
	@echo "=== BUILD INFO"
	@echo "COMMIT:           ${COMMIT}"
	@echo "TIMESTAMP:        ${TIMESTAMP}"
	@echo "SEMVER:           ${SEMVER}"
	@echo "=== CONFIG INFO"
	@echo "ORGANISATION:     ${ORGANISATION}"
	@echo "DATASET:          ${DATASET}"
	@echo "FIXTURES:         ${FIXTURES}"
	@echo "BUCKET_NAME:      ${BUCKET_NAME}"
	@echo "=== PARAMS"
	@echo "BUCKET_PROFILE:   ${BUCKET_PROFILE}"
	@echo "IMPORT_DOWNLOAD:  ${IMPORT_DOWNLOAD}"
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

## Run the api from dev source - will copy existing db to location
api:
	@mkdir -p ./servers/api/databases
	@cp ./builds/databases/api.db ./servers/api/databases && echo "copied db ${tick}" || echo "no db to copy"
	@cd ./servers/api && go run main.go
.PHONY: api

## Run the front from dev source
front:
	@cd ./servers/front && go run main.go
.PHONY: front

#========= DOWNLOAD DATA =========
data/download: build
	@mkdir -p ./builds/databases
	@aws-vault exec ${BUCKET_PROFILE} -- aws s3 sync --quiet s3://${BUCKET_NAME}/databases/api.db ./builds/databases/api.db
.PHONY: data/download

data/upload: build
	@mkdir -p ./builds/databases
	@aws-vault exec ${BUCKET_PROFILE} -- aws s3 cp --sse AES256 --recursive ./builds/databases/ s3://${BUCKET_NAME}/databases/
.PHONY: data/upload
#========= IMPORT DATA =========
# Import all old data - order is important due to data gaps
import: build
	@cd ./builds/ && aws-vault exec ${BUCKET_PROFILE} -- env GITHUB_ACCESS_TOKEN=${GITHUB_ACCESS_TOKEN} ./convertor --download=${IMPORT_DOWNLOAD}
	@cd ./builds && ./importer -type=github-standards -file=./converted-data/github_standards.json
	@cd ./builds && ./importer -type=github-releases -file=./converted-data/github_releases.json
	@cd ./builds && ./importer -type=aws-uptime -file=./converted-data/aws_uptime.json
	@cd ./builds && ./importer -type=aws-costs -file=./converted-data/aws_costs.json
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
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/api ./servers/api/main.go && echo "${tick}"
.PHONY: build/servers/api

## Build the api into build directory
build/servers/front:
	@echo -n "[building] servers/front ................ "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/front ./servers/front/main.go && echo "${tick}"
.PHONY: build/servers/front

## build the convertor tool
build/convertor:
	@echo -n "[building] convertor .................... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/convertor ./convertor/main.go && echo "${tick}"
.PHONY: build/convertor


## build the importer tool
build/importer:
	@echo -n "[building] importer ..................... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/importer ./importer/main.go && echo "${tick}"
.PHONY: build/importer

## build all collectors
build/collectors: build/collectors/awscosts build/collectors/awsuptime build/collectors/githubstandards build/collectors/githubreleases
.PHONY: build/collectors

## build the aws costs collector
build/collectors/awscosts:
	@echo -n "[building] collectors/awscosts .......... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/awscosts ./collectors/awscosts/main.go && echo "${tick}"
.PHONY: build/collectors/awscosts

## build the aws uptime collector
build/collectors/awsuptime:
	@echo -n "[building] collectors/awsuptime ......... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/awsuptime ./collectors/awsuptime/main.go && echo "${tick}"
.PHONY: build/collectors/awsuptime

## build the github standards collector
build/collectors/githubstandards:
	@echo -n "[building] collectors/githubstandards ... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/githubstandards ./collectors/githubstandards/main.go && echo "${tick}"
.PHONY: build/collectors/githubstandards

## build the github releases collector
build/collectors/githubreleases:
	@echo -n "[building] collectors/githubreleases .... "
	@env CGO_ENABLED=1 go build -ldflags=${LDFLAGS} -o ${BUILD_DIR}/githubreleases ./collectors/githubreleases/main.go && echo "${tick}"
.PHONY: build/collectors/githubreleases

