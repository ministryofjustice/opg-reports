SHELL := $(shell which bash)
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | sed 's/x86_64/amd64/')

# version of the gov uk front end to download
GOVUK_FRONT_VERSION := "v5.4.0"
# location of go builds
BUILD_FOLDER := ${ROOT_DIR}builds
# location for this arch
OS_ARCH := ${OS}_${ARCH}
BUILD_ARCH_FOLDER := ${BUILD_FOLDER}/${OS_ARCH}
# where govuk assets are downloaded into
GOVUK_DOWNLOAD_FOLDER := ${BUILD_FOLDER}/govuk-frontend
# location of go report commands
REPORTS_FOLDER := ${ROOT_DIR}cmd/report
# location of where services are based from
SERVICES_FOLDER := ${ROOT_DIR}services
# location of the api service
API_FOLDER := ${SERVICES_FOLDER}/api
# location of the api services data files
API_DATA_FOLDER := ${API_FOLDER}/data
# location of the web front end
FRONT_FOLDER := ${SERVICES_FOLDER}/front
# location to download data from the remote bucket into
BUCKET_FOLDER := ./from-bucket
# aws vault profile to use to connect to the dev bucket
AWS_PROFILE ?= shared-development-operator
# name of the dev bucket
BUCKET ?= report-data-development

# file used to track os, arch and build folders for use
# in github actions
GO_BUILD_INFO := "build.log"
# these vars are overriden per go-* target, but each
# type uses the same recipe so not repeating
GO_SOURCE_FOLDER := ${API_FOLDER}
GO_TARGET_FOLDER := ${BUILD_ARCH_FOLDER}/api
GO_BIN_NAME := api

.DEFAULT_GOAL: all
.PHONY: test tests benchmarks coverage assets assets-front assets-api

all:
	@echo "Nothing to run, choose a target."


##############################
# TESTS
##############################

# run a test based on the $name passed
# pass along github token from env and setup log levels and destinations
test:
	@go clean -testcache
	@clear && env GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)"

tests:
	@go clean -testcache
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -cover -covermode=count -v ./...

coverage:
	@rm -Rf ./code-coverage.out
	@go clean -testcache
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover -v ./...
	@go tool cover -html=code-coverage.out

benchmarks:
	@go clean -testcache
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1  -v ./... -bench=. -run=xxx -benchmem

##############################
# DOCKER
# main commands
##############################

build-production: assets
	@echo "-----"
	@echo "[Docker](production) Building..."
	@env DOCKER_BUILDKIT=0 docker compose build --no-cache
	@echo "[Docker](production) Built."
# build dev version
build: assets
	@echo "-----"
	@echo "[Docker](development) Building..."
	@env DOCKER_BUILDKIT=0 docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build
	@echo "[Docker](development) Built."

down:
	@echo "-----"
	@echo "[Docker] down"
	@docker compose down
clean: down
	@echo "-----"
	@echo "[Docker] clean"
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"
dev: clean build
	@echo "-----"
	@echo "[Docker] Starting..."
	docker compose --verbose -f docker-compose.yml -f docker/docker-compose.dev.yml up -d api front
up:
	@echo "-----"
	@echo "[Docker] Running with images from ECR..."
	docker compose --verbose -f docker-compose.yml up -d api front


##############################
# ASSETS
##############################

# get data from s3 for (production)
assets-api:
	@echo "-----"
	@echo "[Assets](api) Building..."
	@echo "	source: [${BUCKET}]"
	@echo "	target: [${API_DATA_FOLDER}]"
	@rm -Rf ${BUCKET_FOLDER}
	@rm -Rf ${API_DATA_FOLDER}
	@if test "$(AWS_SESSION_TOKEN)" = "" ; then \
		echo "	warning: AWS_SESSION_TOKEN not set, running as aws-vault profile [${AWS_PROFILE}] "; \
		aws-vault exec ${AWS_PROFILE} -- aws s3 sync --quiet s3://${BUCKET} ${BUCKET_FOLDER}; \
	else \
		echo "	AWS_SESSION_TOKEN set, running as is"; \
		aws s3 sync --quiet s3://${BUCKET} ${BUCKET_FOLDER}; \
	fi
	@echo "	note: You can use a different bucket by running the command with BUCKET=name added."
	@mv ${BUCKET_FOLDER} ${API_DATA_FOLDER}
	@echo "[Assets](api) Built."

# get the gov uk front end assets and move them into local folders
assets-front:
	@echo "-----"
	@echo "[Assets](front) Building..."
	@echo "	source: [alphagov/govuk-frontend@${GOVUK_FRONT_VERSION}]"
	@echo "	target: [${SERVICES_FOLDER}/front/assets/]"
	@rm -Rf ${GOVUK_DOWNLOAD_FOLDER}
	@rm -Rf ${SERVICES_FOLDER}/front/assets/css/
	@rm -Rf ${SERVICES_FOLDER}/front/assets/fonts/
	@rm -Rf ${SERVICES_FOLDER}/front/assets/images/
	@rm -Rf ${SERVICES_FOLDER}/front/assets/manifest.json
	@mkdir -p ${GOVUK_DOWNLOAD_FOLDER}
	@cd ${GOVUK_DOWNLOAD_FOLDER} && gh release download ${GOVUK_FRONT_VERSION} -R alphagov/govuk-frontend
	@cd ${GOVUK_DOWNLOAD_FOLDER} && unzip -qq release-${GOVUK_FRONT_VERSION}.zip
	@cd ${GOVUK_DOWNLOAD_FOLDER} && mkdir -p ./assets/css/ && mv govuk-frontend-*.css* ./assets/css/
	@mv ${GOVUK_DOWNLOAD_FOLDER}/assets/css/ ${SERVICES_FOLDER}/front/assets/
	@mv ${GOVUK_DOWNLOAD_FOLDER}/assets/fonts/ ${SERVICES_FOLDER}/front/assets/
	@mv ${GOVUK_DOWNLOAD_FOLDER}/assets/images/ ${SERVICES_FOLDER}/front/assets/
	@mv ${GOVUK_DOWNLOAD_FOLDER}/assets/manifest.json ${SERVICES_FOLDER}/front/assets/
	@rm -Rf ${GOVUK_DOWNLOAD_FOLDER}
	@echo "[Assets](front) Built."

assets: assets-api assets-front
	@echo "[Assets] Built all"

##############################
# RELEASE ARTIFACTS
# Will build the go code. Uses target specific variables to
# set folders and binary names
##############################

gbuildinfo:
	@rm -f ${GO_BUILD_INFO}
	@echo "arch ${OS_ARCH}" >> ${GO_BUILD_INFO}
	@echo "build_folder ${BUILD_FOLDER}" >> ${GO_BUILD_INFO}
	@echo "build_arch_folder ${BUILD_ARCH_FOLDER}" >> ${GO_BUILD_INFO}

# set variables for the api binary
go-api: GO_SOURCE_FOLDER=${SERVICES_FOLDER}/api
go-api: GO_TARGET_FOLDER=${BUILD_ARCH_FOLDER}/api
go-api: GO_BIN_NAME=api
# set variables for the the front binary
go-front: GO_SOURCE_FOLDER=${SERVICES_FOLDER}/front
go-front: GO_TARGET_FOLDER=${BUILD_ARCH_FOLDER}/front
go-front: GO_BIN_NAME=front
# set variables for the github standards report
go-report-gh-standards: GO_SOURCE_FOLDER=${REPORTS_FOLDER}/github/standards
go-report-gh-standards: GO_TARGET_FOLDER=${BUILD_ARCH_FOLDER}/reports
go-report-gh-standards: GO_BIN_NAME=github_standards
# set variables for the aws monthly costs report
go-report-aws-monthly-costs: GO_SOURCE_FOLDER=${REPORTS_FOLDER}/aws/cost/monthly
go-report-aws-monthly-costs: GO_TARGET_FOLDER=${BUILD_ARCH_FOLDER}/reports
go-report-aws-monthly-costs: GO_BIN_NAME=aws_cost_monthly
#
go-report-aws-daily-uptime: GO_SOURCE_FOLDER=${REPORTS_FOLDER}/aws/uptime/daily
go-report-aws-daily-uptime: GO_TARGET_FOLDER=${BUILD_ARCH_FOLDER}/reports
go-report-aws-daily-uptime: GO_BIN_NAME=aws_uptime_daily

# these all share the same recipe, but each target changes the $GO_ variables
# to build the correct binary
go-api go-front go-report-gh-standards go-report-aws-monthly-costs go-report-aws-daily-uptime: gbuildinfo
	@echo "-----"
	@echo "[Go](${GO_BIN_NAME}) Building..."
	@echo "	source: [${GO_SOURCE_FOLDER}]"
	@echo "	target: [${GO_TARGET_FOLDER}]"
	@echo "	binary: [${GO_BIN_NAME}]"
	@mkdir -p ${GO_TARGET_FOLDER}
	@rm -Rf ${GO_TARGET_FOLDER}/${GO_BIN_NAME}
	@cd ${GO_SOURCE_FOLDER} && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${GO_TARGET_FOLDER}/${GO_BIN_NAME} main.go
	@echo "[Go](${GO_BIN_NAME}) Built."
	@echo "${GO_BIN_NAME}_target_folder ${GO_TARGET_FOLDER}" >> ${GO_BUILD_INFO}

go-reports: go-report-gh-standards go-report-aws-monthly-costs go-report-aws-daily-uptime
	@echo "[Go] Built reports."

go go-all: go-api go-front go-reports
	@echo "[Go] Built all"

godocs:
	@echo "-----"
	@echo "[Go](docs) building docs "
	@./scripts/go-docs.sh
	@echo "[Go](docs) done"
##############################
# DEV
##############################

dev-run-api:
	@echo "Running api..."
	@clear && cd ${API_FOLDER} && go run main.go

dev-run-front:
	@echo "Running front..."
	@clear && cd ${FRONT_FOLDER} && env CONFIG_FILE="./config.json" go run main.go
