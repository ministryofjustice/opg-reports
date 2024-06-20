SHELL := $(shell which bash)
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_FOLDER := ${ROOT_DIR}builds
REPORTS_DIR := ${ROOT_DIR}cmd/report

# check and set the correct goarch
ifeq (${ARCH}, 'x86_64')
	BUILD_ARCH := 'amd64'
else
	BUILD_ARCH := ${ARCH}
endif

.PHONY: test tests benchmarks coverage
# clean out any build files
clean:
	@rm -Rf ${BUILD_FOLDER}
	@mkdir -p ${BUILD_FOLDER} ${BUILD_FOLDER}/report

# Builds all the commands in the report folder into the builds folder
build_reports: ${REPORTS_DIR}/*
	@for f in $^; do  \
        echo "building: " $${f##*/} && \
		env GOOS=${OS} GOARCH=${BUILD_ARCH} go build -o ${BUILD_FOLDER}/report/$${f##*/} $${f}/main.go ; \
    done

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
