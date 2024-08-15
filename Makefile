.DEFAULT_GOAL: all
.PHONY: test tests benchmarks coverage clean data go-build

all:
	@echo "Nothing to run, choose a target."

##############################
AWS_VAULT_PROFILE ?= shared-development-operator
AWS_BUCKET ?= report-data-development
SERVICES ?= api front
##############################
AWS_VAULT_COMMAND = echo "using existing session token" &&
##############################
ifndef AWS_SESSION_TOKEN
AWS_VAULT_COMMAND = aws-vault exec ${AWS_VAULT_PROFILE} --
endif
# docker images
images := $(shell docker images -a | grep 'opg-reports/*' | awk '{print $$1":"$$2}')

##############################
# TESTS
##############################
# run a test based on the $name passed
# pass along github token from env and setup log levels and destinations
test:
	@go clean -testcache
	@clear
	@echo "============== test: [$(name)] =============="
	@env CGO_ENABLED=1 GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)"

tests:
	@go clean -testcache
	@clear
	@echo "============== tests =============="
	@env env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -cover -covermode=atomic -v ./...

coverage:
	@rm -Rf ./code-coverage.out
	@clear
	@echo "============== coverage =============="
	@env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover -v ./...
	@go tool cover -html=code-coverage.out

benchmarks:
	@clear
	@echo "============== benchmarks =============="
	@env LOG_LEVEL="warn" LOG_TO="stdout" go test -v ./... -bench=. -run=xxx -benchmem -benchtime=10s

benchmark:
	@clear
	@echo "============== benchmark: [$(name)] =============="
	@env LOG_LEVEL="info" LOG_TO="stdout" go test -v ./... -bench=$(name) -run=xxx -benchmem -benchtime=10s


##############################
# DATA
##############################
sqlc:
	@cd ./datastore/github_standards && sqlc generate
	@cd ./datastore/aws_costs && sqlc generate

data: vars
	@mkdir -p ./builds/api/github_standards/data
	@mkdir -p ./builds/api/aws_costs/data

	${AWS_VAULT_COMMAND} aws s3 sync --quiet s3://${AWS_BUCKET}/github_standards ./builds/api/github_standards/data/ || echo bucket_github_standards_failed; \
	${AWS_VAULT_COMMAND} aws s3 sync --quiet s3://${AWS_BUCKET}/aws_costs ./builds/api/aws_costs/data/ || echo bucket_aws_costs_failed; \

vars:
	@echo "AWS_VAULT_PROFILE: ${AWS_VAULT_PROFILE}"
	@echo "AWS_BUCKET: ${AWS_BUCKET}"
	@echo "AWS_VAULT_COMMAND: ${AWS_VAULT_COMMAND}"
	@echo "SERVICES: ${SERVICES}"

##############################
# DOCKER BUILD
##############################
down:
	@docker compose down

stop:
	@docker compose stop ${SERVICES}

start:
	@docker compose start ${SERVICES}

clean: down
	@rm -f ./servers/api/*.db
	@rm -f ./servers/api/*.csv
# @rm -Rf ./servers/front/govuk
	@rm -Rf ./builds
	@mkdir -p ./builds
	@docker image rm $(image) || echo "ok"
	@docker compose rm api front
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"

build: data
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker/docker-compose.dev.yml \
		build ${SERVICES} \
		--parallel

up: build
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker/docker-compose.dev.yml \
		up \
		-d ${SERVICES}

# production versions
build-production: data
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		build ${SERVICES} \
		--parallel

up-production:
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		up \
		-d ${SERVICES}

##############################
go-build: clean
	@go build -o ./builds/front/front_server ./servers/front/main.go
	@go build -o ./builds/api/api_server ./servers/api/main.go
	@go build -o ./builds/api/seed_cmd ./commands/seed/main.go
	@go build -o ./builds/commands/github_standards ./commands/github_standards/main.go
	@go build -o ./builds/commands/aws_costs ./commands/aws_costs/main.go


##############################
# close approx of the dockerfile for setup without docker
##############################
mirror-api: clean data go-build
	@mkdir -p ./builds/api/github_standards/data
	@mkdir -p ./builds/api/aws_costs/data
#	github_standards
	@cp ./datastore/github_standards/github_standards*.sql ./builds/api/github_standards/
	@echo "seeding github_standards"
	@./builds/api/seed_cmd \
		-table github_standards \
		-db ./builds/api/github_standards.db \
		-schema ./builds/api/github_standards/github_standards.sql \
		-data "./builds/api/github_standards/data/*.json"
#	aws_costs
	@cp ./datastore/aws_costs/aws_costs*.sql ./builds/api/aws_costs/
	@echo "seeding aws_costs"
	@./builds/api/seed_cmd \
		-table aws_costs \
		-db ./builds/api/aws_costs.db \
		-schema ./builds/api/aws_costs/aws_costs.sql \
		-data "./builds/api/aws_costs/data/*.json"
