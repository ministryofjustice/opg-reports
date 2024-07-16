SHELL := $(shell which bash)
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
# check and set the correct goarch
ifeq (${ARCH}, 'x86_64')
	ARCH := 'amd64'
endif

# version of the gov uk front end to download
GOVUK_FRONT_VERSION := "v5.4.0"
# location of go builds
BUILD_FOLDER := ${ROOT_DIR}builds
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
AWS_PROFILE ?= "shared-development"
# name of the dev bucket
BUCKET ?= "report-data-development"


.PHONY: test tests benchmarks coverage assets assets-front assets-api

##############################
# TESTS
##############################

# run a test based on the $name passed
# pass along github token from env and setup log levels and destinations
test:
	@go clean -testcache
	@clear && env GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)"

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
	@env DOCKER_BUILDKIT=0 docker compose build --no-cache
# build dev version
build: assets
	@env DOCKER_BUILDKIT=0 docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build
down:
	@docker compose down
clean: down
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"
up: clean build-dev
	docker compose --verbose -f docker-compose.yml -f docker/docker-compose.dev.yml up -d api front


##############################
# ASSETS
##############################

# get data from s3 for (production)
assets-api:
	@rm -Rf ${BUCKET_FOLDER}
	@rm -Rf ${API_DATA_FOLDER}
	@if test "$(AWS_SESSION_TOKEN)" = "" ; then \
		echo "warning: AWS_SESSION_TOKEN not set, running as aws-vault profile [${AWS_PROFILE}] "; \
		aws-vault exec ${AWS_PROFILE} -- aws s3 sync s3://${BUCKET} ${BUCKET_FOLDER}; \
	else \
		echo "AWS_SESSION_TOKEN set, running as is"; \
		aws s3 sync s3://${BUCKET} ${BUCKET_FOLDER}; \
	fi
	@mv ${BUCKET_FOLDER} ${API_DATA_FOLDER}

# get the gov uk front end assets and move them into local folders
assets-front:
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
	@echo "Downloaded alphagov/govuk-frontend@${GOVUK_FRONT_VERSION} to ${SERVICES_FOLDER}/front/assets/"

assets: assets-api assets-front
##############################
# RELEASE ARTIFACTS
# Can also be used for local binary builds
##############################

# builds the reports into a build folder - typically for a release artifact rather than docker setup
artifact-reports:
	@rm -Rf ${BUILD_FOLDER}/${OS}_${ARCH}/reports/
	@mkdir -p ${BUILD_FOLDER}/${OS}_${ARCH}/reports/
	@cd ${REPORTS_FOLDER}/aws/cost/monthly/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/reports/aws_cost_monthly main.go
	@cd ${REPORTS_FOLDER}/github/standards/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/reports/github_standards main.go

# builds the api into a build folder - typically for a release artifact rather than docker setup
artifact-api: assets-api
	@rm -Rf ${BUILD_FOLDER}/${OS}_${ARCH}/api/
	@mkdir -p ${BUILD_FOLDER}/${OS}_${ARCH}/api/
	@cd ${API_FOLDER}/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/api/api main.go

# builds the api into a build folder - typically for a release artifact rather than docker setup
artifact-front: assets-front
	@rm -Rf ${BUILD_FOLDER}/${OS}_${ARCH}/front/
	@mkdir -p ${BUILD_FOLDER}/${OS}_${ARCH}/front/
	@cd ${API_FOLDER}/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/front/front main.go

artifacts: artifact-reports artifact-api artifact-front

##############################
# DEV
##############################

dev-run-api:
	@clear && cd ${API_FOLDER} && go run main.go

dev-run-front:
	@clear && cd ${FRONT_FOLDER} && go run main.go
