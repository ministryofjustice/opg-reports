.DEFAULT_GOAL: all
.PHONY: test tests benchmarks coverage go-build

all:
	@echo "Nothing to run, choose a target."



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
# DOCKER BUILD
##############################
docker-clean: docker-down
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"
docker-build:
	@env DOCKER_BUILDKIT=0 docker compose --verbose -f docker-compose.yml -f docker/docker-compose.dev.yml build
docker-up:
	@env DOCKER_BUILDKIT=0 docker compose --verbose -f docker-compose.yml -f docker/docker-compose.dev.yml up -d api front
docker-down:
	@docker compose down

# production versions
docker-build-production:
	@env DOCKER_BUILDKIT=0 docker compose --verbose -f docker-compose.yml build
docker-up-production:
	@env DOCKER_BUILDKIT=0 docker compose --verbose -f docker-compose.yml up -d api front

##############################
# GO BUILD
# - build all go binaries at once and push to ./builds/go/
#   using goreleaser
##############################

AWS_VAULT_PROFILE ?= shared-development-operator
AWS_BUCKET ?= report-data-development
GO_RELEASER = $(shell which goreleaser)

go-build:
	@echo "goreleaser: ${GO_RELEASER}"
	@env AWS_VAULT_PROFILE=${AWS_VAULT_PROFILE} AWS_BUCKET=${AWS_BUCKET} ${GO_RELEASER} build --clean --single-target --skip=validate
	@rm -f ./builds/binaries/*.json
	@rm -f ./builds/binaries/*.yml


go-run-api: go-build
	@cd ./builds/api/ && ./api_server

go-run-front: go-build
	@cd ./builds/front/ && ./front_server

