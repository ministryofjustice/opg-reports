.DEFAULT_GOAL: all
.PHONY: test tests benchmarks coverage go-build

all:
	@echo "Nothing to run, choose a target."


##############################
AWS_VAULT_PROFILE ?= shared-development-operator
AWS_BUCKET ?= report-data-development
RUN_DOWNLOAD ?= yes
##############################
AWS_VAULT_COMMAND = echo "using session token" &&
GO_RELEASER = $(shell which goreleaser)
##############################
ifndef AWS_SESSION_TOKEN
AWS_VAULT_COMMAND = aws-vault exec ${AWS_VAULT_PROFILE} --
endif

##############################
# TESTS
##############################

# run a test based on the $name passed
# pass along github token from env and setup log levels and destinations
test:
	@clear
	@echo "============== test: [$(name)] =============="
	@env CGO_ENABLED=1 GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)"

tests:
	@clear
	@echo "============== tests =============="
	@env env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -cover -covermode=atomic -v ./...

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

clean: docker-clean
	@rm -Rf ./builds


##############################
# DATA
##############################
csv: vars
# 	download github_standards data
	@mkdir -p ./builds/api/github_standards
	@if [[ "${RUN_DOWNLOAD}" == "yes" ]]; then \
		${AWS_VAULT_COMMAND} aws s3 sync s3://${AWS_BUCKET}/github_standards ./builds/api/github_standards/ ; \
	fi


vars:
	@echo "RUN_DOWNLOAD: ${RUN_DOWNLOAD}"
	@echo "AWS_VAULT_PROFILE: ${AWS_VAULT_PROFILE}"
	@echo "AWS_BUCKET: ${AWS_BUCKET}"
	@echo "AWS_VAULT_COMMAND: ${AWS_VAULT_COMMAND}"
	@echo "GO_RELEASER: ${GO_RELEASER}"

##############################
# DOCKER BUILD
##############################
docker-clean: docker-down
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"

docker-build: csv
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker/docker-compose.dev.yml \
		build

docker-build-production: csv
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		build

docker-up: docker-clean csv
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker/docker-compose.dev.yml \
		up \
		-d api

docker-down:
	@docker compose down
