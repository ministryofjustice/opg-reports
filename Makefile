SHELL := $(shell which bash)
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

VERSION_UK_GOV_FRONT := "v5.4.0"


OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
# check and set the correct goarch
ifeq (${ARCH}, 'x86_64')
	ARCH := 'amd64'
endif

BUILD_FOLDER := ${ROOT_DIR}builds
REPORTS_SRC := ${ROOT_DIR}cmd/report
API_SRC := ${ROOT_DIR}/services/api
FRONT_SRC := ${ROOT_DIR}/services/front

.PHONY: test tests benchmarks coverage assets-front assets-api

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
##############################

down:
	@docker compose down
clean: down
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"
up: clean docker-dev-build
	docker compose --verbose -f docker-compose.yml -f docker/docker-compose.dev.yml up -d api front

# docker build and dev stages
docker-build: assets-front assets-api
	@env DOCKER_BUILDKIT=0 docker compose -f docker-compose.yml build
# build docker setup - call assets first to make sure files are copied
docker-dev-build: assets-front assets-api
	@env DOCKER_BUILDKIT=0 docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build

##############################
# ASSETS
##############################
# get data from s3 (eventually)
assets-api:
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


##############################
# BUILD / RELEASES
##############################
build-reports:
	@rm -Rf ${BUILD_FOLDER}/${OS}_${ARCH}/reports/
	@mkdir -p ${BUILD_FOLDER}/${OS}_${ARCH}/reports/
	@cd ${REPORTS_SRC}/aws/cost/monthly/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/reports/aws_cost_monthly main.go
	@cd ${REPORTS_SRC}/github/standards/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/reports/github_standards main.go

build-api: assets-api
	@rm -Rf ${BUILD_FOLDER}/${OS}_${ARCH}/api/
	@mkdir -p ${BUILD_FOLDER}/${OS}_${ARCH}/api/
	@cd ${API_SRC}/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/api/api main.go

build-front: assets-front
	@rm -Rf ${BUILD_FOLDER}/${OS}_${ARCH}/front/
	@mkdir -p ${BUILD_FOLDER}/${OS}_${ARCH}/front/
	@cd ${API_SRC}/ && go mod download && env GOOS=${OS} GOARCH=${ARCH} go build -o ${BUILD_FOLDER}/${OS}_${ARCH}/front/front main.go

build: build-reports build-api build-front


##############################
# DEV
##############################

dev-api:
	@clear && cd ./services/api/ && go run main.go

dev-front:
	@clear && cd ./services/front/ && go run main.go
