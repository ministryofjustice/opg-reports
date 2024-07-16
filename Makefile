SHELL := $(shell which bash)
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
# check and set the correct goarch
ifeq (${ARCH}, 'x86_64')
	ARCH := 'amd64'
endif

BUILD_FOLDER := ${ROOT_DIR}builds
REPORTS_SRC := ${ROOT_DIR}cmd/report
API_SRC := ${ROOT_DIR}services/api
FRONT_SRC := ${ROOT_DIR}services/front

DEV_AWS_VAULT_PROFILE := "shared-development"
DEV_BUCKET := "report-data-development"
# PROD_AWS_VAULT_PROFILE := "shared-production"
# PROD_BUCKET := "report-data-production"
PROD_AWS_VAULT_PROFILE := "shared-development"
PROD_BUCKET := "report-data-development"
LOCAL_FOLDER := "./bucket-sync"
API_DATA_FOLDER := ${API_SRC}/data

VERSION_UK_GOV_FRONT := "v5.4.0"

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

build: assets
	@env DOCKER_BUILDKIT=0 docker compose build --no-cache
# build dev version
build-dev: assets-dev
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
# get data from s3 for development
assets-api-dev:
	@clear
	@rm -Rf ${LOCAL_FOLDER}
	@rm -Rf ${API_DATA_FOLDER}
	@if test "$(AWS_SESSION_TOKEN)" = "" ; then \
		echo "warning: AWS_SESSION_TOKEN not set, running as aws-vault profile [${DEV_AWS_VAULT_PROFILE}] "; \
		aws-vault exec ${DEV_AWS_VAULT_PROFILE} -- aws s3 sync s3://${DEV_BUCKET} ${LOCAL_FOLDER}; \
	else \
		echo "AWS_SESSION_TOKEN set, running as is"; \
		aws s3 sync s3://${DEV_BUCKET} ${LOCAL_FOLDER}; \
	fi
	@mv ${LOCAL_FOLDER} ${API_DATA_FOLDER}

# get data from s3 for (production)
assets-api:
	@clear
	@rm -Rf ${LOCAL_FOLDER}
	@rm -Rf ${API_DATA_FOLDER}
	@if test "$(AWS_SESSION_TOKEN)" = "" ; then \
		echo "warning: AWS_SESSION_TOKEN not set, running as aws-vault profile [${PROD_AWS_VAULT_PROFILE}] "; \
		aws-vault exec ${PROD_AWS_VAULT_PROFILE} -- aws s3 sync s3://${PROD_BUCKET} ${LOCAL_FOLDER}; \
	else \
		echo "AWS_SESSION_TOKEN set, running as is"; \
		aws s3 sync s3://${PROD_BUCKET} ${LOCAL_FOLDER}; \
	fi
	@mv ${LOCAL_FOLDER} ${API_DATA_FOLDER}

# get the gov uk front end assets and move them into local folders
assets-front:
	@rm -Rf ./builds/govuk-frontend
	@rm -Rf ./services/front/assets/css/
	@rm -Rf ./services/front/assets/fonts/
	@rm -Rf ./services/front/assets/images/
	@rm -Rf ./services/front/assets/manifest.json
	@mkdir -p ./builds/govuk-frontend
	@cd ./builds/govuk-frontend && gh release download ${VERSION_UK_GOV_FRONT} -R alphagov/govuk-frontend
	@cd ./builds/govuk-frontend && unzip -qq release-${VERSION_UK_GOV_FRONT}.zip
	@cd ./builds/govuk-frontend && mkdir -p ./assets/css/ && mv govuk-frontend-*.css* ./assets/css/
	@mv ./builds/govuk-frontend/assets/css/ ./services/front/assets/
	@mv ./builds/govuk-frontend/assets/fonts/ ./services/front/assets/
	@mv ./builds/govuk-frontend/assets/images/ ./services/front/assets/
	@mv ./builds/govuk-frontend/assets/manifest.json ./services/front/assets/
	@rm -Rf ./builds/govuk-frontend
	@echo "Downloaded alphagov/govuk-frontend@${VERSION_UK_GOV_FRONT} to ./services/front/assets/"

assets-dev: assets-api-dev assets-front
assets: assets-api assets-front
##############################
# RELEASE ARTIFACTS
# Can also be used for local binary builds
##############################

# builds the reports into a build folder - typically for a release artifact rather than docker setup
artifact-reports:
	@rm -Rf ${BUILD_FOLDER}/${OS}_${ARCH}/reports/
	@mkdir -p ${BUILD_FOLDER}/${OS}_${ARCH}/reports/
	@cd ${REPORTS_SRC}/aws/cost/monthly/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/reports/aws_cost_monthly main.go
	@cd ${REPORTS_SRC}/github/standards/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/reports/github_standards main.go

# builds the api into a build folder - typically for a release artifact rather than docker setup
artifact-api: assets-api
	@rm -Rf ${BUILD_FOLDER}/${OS}_${ARCH}/api/
	@mkdir -p ${BUILD_FOLDER}/${OS}_${ARCH}/api/
	@cd ${API_SRC}/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/api/api main.go

# builds the api into a build folder - typically for a release artifact rather than docker setup
artifact-front: assets-front
	@rm -Rf ${BUILD_FOLDER}/${OS}_${ARCH}/front/
	@mkdir -p ${BUILD_FOLDER}/${OS}_${ARCH}/front/
	@cd ${API_SRC}/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/front/front main.go

artifacts: artifact-reports artifact-api artifact-front

##############################
# DEV
##############################

dev-run-api:
	@clear && cd ./services/api/ && go run main.go

dev-run-front:
	@clear && cd ./services/front/ && go run main.go
