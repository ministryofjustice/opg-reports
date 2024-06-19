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
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -v ./... -run="$(name)"

tests:
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -cover -covermode=count -v ./...

coverage:
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -covermode=count -coverprofile=coverage.out -cover -v ./...
	@go tool cover -html=coverage.out

benchmarks:
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -v ./... -bench=. -run=xxx -benchmem
