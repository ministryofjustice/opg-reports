SHELL := $(shell which bash)
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION_UK_GOV_FRONT := "v5.4.0"
# BUILD_FOLDER := ${ROOT_DIR}builds
# REPORTS_DIR := ${ROOT_DIR}cmd/report

# check and set the correct goarch
ifeq (${ARCH}, 'x86_64')
	BUILD_ARCH := 'amd64'
else
	BUILD_ARCH := ${ARCH}
endif

.PHONY: test tests benchmarks coverage govuk-frontend

govuk-frontend:
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

test:
	@go clean -testcache
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)"

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
