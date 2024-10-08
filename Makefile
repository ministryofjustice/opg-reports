#================================
AWS_VAULT_PROFILE ?= shared-development-operator
AWS_BUCKET ?= report-data-development
SERVICES ?= api front
AWS_VAULT_COMMAND = echo "using existing session token" &&
#================================
ifndef AWS_SESSION_TOKEN
AWS_VAULT_COMMAND = aws-vault exec ${AWS_VAULT_PROFILE} --
endif
#================================
images := $(shell docker images -a | grep 'opg-reports/*' | awk '{print $$1":"$$2}')

.DEFAULT_GOAL = help
#================================
# TESTS
#================================
.PHONY: tests tests/all tests/named tests/benchmarks tests/benchmark tests/coverage
## run all the go suite tests
tests: tests/all

## run all tests
tests/all:
	@go clean -testcache
	@clear
	@echo "============== tests =============="
	@env env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -cover -covermode=atomic -v ./...

## run go suite tests that match the file pattern
## usage:
## `make tests/named name=<pattern>`
tests/named test:
	@go clean -testcache
	@clear
	@echo "============== test: [$(name)] =============="
	@env CGO_ENABLED=1 GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)"

## run go benchmarking
tests/benchmarks:
	@clear
	@echo "============== benchmarks =============="
	@env LOG_LEVEL="warn" LOG_TO="stdout" go test -v ./... -bench=. -run=xxx -benchmem -benchtime=10s

## run the named benchmark
## usage:
## `make benchmark name="<pattern>"`
tests/benchmark:
	@clear
	@echo "============== benchmark: [$(name)] =============="
	@env LOG_LEVEL="info" LOG_TO="stdout" go test -v ./... -bench=$(name) -run=xxx -benchmem -benchtime=10s

## Run the go code coverage tool
tests/coverage:
	@rm -Rf ./code-coverage.out
	@clear
	@echo "============== coverage =============="
	@env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover -v ./...
	@go tool cover -html=code-coverage.out


#================================
# DOCKER
#================================
.PHONY: docker docker/up docker/down docker/stop docker/start docker/build docker/production/build docker/production/up
## short form alias for the docker/up command which runs docker compose build and up commands
docker: docker/up

## run docker compose down, turning off all docker containers
docker/down:
	@docker compose down

## run docker compose stop
## usage:
## `make docker/stop`
## `make docker/stop SERVICES="<A> <B>"`
docker/stop:
	@docker compose stop ${SERVICES}

## run docker compose start
## usage:
## `make docker/start`
## `make docker/start SERVICES="<A> <B>"`
docker/start:
	@docker compose start ${SERVICES}

## run docker compose build with the `./docker/docker-compose.dev.yml` file
## calls `data` target as well to sync content down before build
## usage:
## `make docker/build`
## `make docker/build SERVICES="<A> <B>"`
docker/build: data
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker/docker-compose.dev.yml \
		build ${SERVICES} \
		--parallel

## run docker compose build and then compose up with the `./docker/docker-compose.dev.yml` file
## usage:
## `make docker/up`
## `make docker/up SERVICES="<A> <B>"`
docker/up: docker/build
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker/docker-compose.dev.yml \
		up \
		-d ${SERVICES}

## run docker compose build
## calls `data` target as well to sync content down before build
## usage:
## `make docker/production/build`
## `make docker/production/build SERVICES="<A> <B>"`
docker/production/build: data
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		build ${SERVICES} \
		--parallel

## Run docker compose up *without* build - so will also pull from registry
## usage:
## `make docker/production/up`
## `make docker/production/up SERVICES="<A> <B>"`
docker/production/up:
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		up \
		-d ${SERVICES}

#================================
# DATA
#================================
.PHONY: data data/sync data/sync/vars
## short form alias for data/sync, which fetches content from buckets
data: data/sync

## run sqlc generate for all known items
data/sqlc:
	@cd ./datastore/github_standards && sqlc generate
#--fork-remove-start
	@cd ./datastore/aws_costs && sqlc generate
	@cd ./datastore/aws_uptime && sqlc generate
#--fork-remove-end

## download all data from bucket
## Can overwrite bucket name to download using:
## `make data/sync AWS_BUCKET="<bucket-name>"
data/sync: data/sync/vars
#	github_standards
	@mkdir -p ./builds/api/github_standards/data
	@echo "getting github_standards" && ${AWS_VAULT_COMMAND} aws s3 sync --quiet s3://${AWS_BUCKET}/github_standards ./builds/api/github_standards/data/ && echo bucket_github_standards_done || echo bucket_github_standards_failed;
#--fork-remove-start
#	aws_costs
	@mkdir -p ./builds/api/aws_costs/data
	@echo "getting aws_costs" && ${AWS_VAULT_COMMAND} aws s3 sync --quiet s3://${AWS_BUCKET}/aws_costs ./builds/api/aws_costs/data/ && echo bucket_aws_costs_done || echo bucket_aws_costs_failed;
#	aws_uptime
	@mkdir -p ./builds/api/aws_uptime/data
	@echo "getting aws_uptime" && ${AWS_VAULT_COMMAND} aws s3 sync --quiet s3://${AWS_BUCKET}/aws_uptime ./builds/api/aws_uptime/data/ && echo bucket_aws_uptime_done || echo bucket_aws_uptime_failed;
#--fork-remove-end

## output current values used by data/sync
data/sync/vars:
	@echo "AWS_VAULT_PROFILE: ${AWS_VAULT_PROFILE}"
	@echo "AWS_BUCKET: ${AWS_BUCKET}"
	@echo "AWS_VAULT_COMMAND: ${AWS_VAULT_COMMAND}"
	@echo "SERVICES: ${SERVICES}"

#================================
# CLEAN
#================================
.PHONY: clean
## removes all generated files and docker images to ensure clean build and run
clean: docker/down
	@rm -f ./servers/api/*.db
	@rm -f ./servers/api/*.csv
	@rm -Rf ./servers/front/govuk
	@rm -Rf ./builds
	@mkdir -p ./builds
	@docker image rm $(images) || echo "ok"
	@docker compose rm api front
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"

#================================
# DEV
#================================
.PHONY: dev dev/build dev/run dev/api dev/front dev/mirror dev/mirror/api seed seed/api
## short form alias for dev/mirror which will build and prep the local env for usage
dev: dev/mirror

## run the local api server
dev/api api:
	@cd ./servers/api && go run main.go

## run the local fron server
dev/front front:
	@cd ./servers/front && go run main.go

## builds all local binaries, calls clean first
dev/build:
	@go build -o ./builds/front/front_server ./servers/front/main.go
	@go build -o ./builds/api/api_server ./servers/api/main.go
	@go build -o ./builds/api/seed_cmd ./commands/seed/main.go
	@go build -o ./builds/commands/github_standards ./commands/github_standards/main.go
#--fork-remove-start
	@go build -o ./builds/commands/aws_costs ./commands/aws_costs/main.go
	@go build -o ./builds/commands/aws_uptime ./commands/aws_uptime/main.go
#--fork-remove-end

#================================
# DATA SEEDING
#================================
## short form alias
seed: seed/api
## seed local databases with data
## note: run the build process first as it uses the build file locations
seed/api:
#	github_standards
	@mkdir -p ./builds/api/github_standards/data
	@cp ./datastore/github_standards/github_standards*.sql ./builds/api/github_standards/
	@echo "seeding github_standards"
	@./builds/api/seed_cmd \
		-table github_standards \
		-db ./builds/api/github_standards.db \
		-schema ./builds/api/github_standards/github_standards.sql \
		-data "./builds/api/github_standards/data/*.json"
	@echo "done github_standards"
#--fork-remove-start
#	aws_costs
	@mkdir -p ./builds/api/aws_costs/data
	@cp ./datastore/aws_costs/aws_costs*.sql ./builds/api/aws_costs/
	@echo "seeding aws_costs"
	@./builds/api/seed_cmd \
		-table aws_costs \
		-db ./builds/api/aws_costs.db \
		-schema ./builds/api/aws_costs/aws_costs.sql \
		-data "./builds/api/aws_costs/data/*.json"
	@echo "done aws_costs"
#	aws_uptime
	@mkdir -p ./builds/api/aws_uptime/data
	@cp ./datastore/aws_uptime/aws_uptime*.sql ./builds/api/aws_uptime/
	@echo "seeding aws_uptime"
	@./builds/api/seed_cmd \
		-table aws_uptime \
		-db ./builds/api/aws_uptime.db \
		-schema ./builds/api/aws_uptime/aws_uptime.sql \
		-data "./builds/api/aws_uptime/data/*.json"
	@echo "done aws_uptime"
#--fork-remove-end

#================================
# MIRROR DOCKER SETUP
#================================
## mirrors build setup
mirror: clean data/sqlc dev/build dev/build data/sync seed mirror/api
## mirror position of data files for the api server
mirror/api:
	@mv ./builds/api/github_standards.db ./servers/api/github_standards.db
#--fork-remove-start
	@mv ./builds/api/aws_costs.db ./servers/api/aws_costs.db
	@mv ./builds/api/aws_uptime.db ./servers/api/aws_uptime.db
#--fork-remove-end


#================================
## Running in local dev
## To setup local folders with data and faster dev process than docker rebuilds, run
## `make mirror`
## This will clean out everything, create builds, sync and seed data and then move
## the databases into the api folder ready
## You can then run `make dev/front` or `make dev/api` to run either and they will
## include generated / fetched data
usage:
	@echo ""
#================================
# HELP
#================================
help:
	@echo "============================"
	@FILE=Makefile ./scripts/help.mk
	@echo "============================"
.PHONY: help
