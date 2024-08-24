##############################
AWS_VAULT_PROFILE ?= shared-development-operator
AWS_BUCKET ?= report-data-development
SERVICES ?= api front
##############################
AWS_VAULT_COMMAND = echo "using existing session token" &&
ifndef AWS_SESSION_TOKEN
AWS_VAULT_COMMAND = aws-vault exec ${AWS_VAULT_PROFILE} --
endif
##############################
images := $(shell docker images -a | grep 'opg-reports/*' | awk '{print $$1":"$$2}')
##############################
.DEFAULT_GOAL: help
#############################
.PHONY: help
help: ## Show this help
	@echo "\nSpecify a command. The choices are:\n"
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;36m%-12s\033[m %s\n", $$1, $$2}'
	@echo ""

##############################
# TESTS
##############################
.PHONY: tests
tests: tests/all ## Run tests/all

.PHONY: tests/named
tests/named: ## Run only tests that match pattern (usage: make test name="<pattern>")
	@go clean -testcache
	@clear
	@echo "============== test: [$(name)] =============="
	@env CGO_ENABLED=1 GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)"

.PHONY: tests/all
tests/all: ## Run all tests
	@go clean -testcache
	@clear
	@echo "============== tests =============="
	@env env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -cover -covermode=atomic -v ./...

.PHONY: tests/coverage
tests/coverage: ## Run code coverage
	@rm -Rf ./code-coverage.out
	@clear
	@echo "============== coverage =============="
	@env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover -v ./...
	@go tool cover -html=code-coverage.out

.PHONY: tests/benchmarks
tests/benchmarks: ## Run all code benchmarks
	@clear
	@echo "============== benchmarks =============="
	@env LOG_LEVEL="warn" LOG_TO="stdout" go test -v ./... -bench=. -run=xxx -benchmem -benchtime=10s

.PHONY: tests/benchmark
tests/benchmark: ## Run the named benchmark (usage: make benchmark name="<pattern>")
	@clear
	@echo "============== benchmark: [$(name)] =============="
	@env LOG_LEVEL="info" LOG_TO="stdout" go test -v ./... -bench=$(name) -run=xxx -benchmem -benchtime=10s


##############################
# DOCKER
##############################
.PHONY: docker
docker: docker/up ## Run docker/up

.PHONY: docker/down
docker/down: ## Run docker compose down, turning off all docker containers
	@docker compose down

.PHONY: docker/stop
docker/stop: ## Run docker compose stop  (usage: make docker/stop || make docker/stop SERVICES="<A> <B>")
	@docker compose stop ${SERVICES}

.PHONY: docker/start
docker/start: ## Run docker compose start (usage: make docker/start || make docker/start SERVICES="<A> <B>")
	@docker compose start ${SERVICES}

.PHONY: docker/build
docker/build: data ## Run docker compose build with dev file (usage: make docker/build || make docker/build SERVICES="<A> <B>")
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker/docker-compose.dev.yml \
		build ${SERVICES} \
		--parallel

.PHONY: docker/up
docker/up: docker/build ## Run docker compose up with dev file - calls docker/build first (usage: make docker/up || make docker/up SERVICES="<A> <B>")
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker/docker-compose.dev.yml \
		up \
		-d ${SERVICES}

.PHONY: docker/production/build
docker/production/build: data ## Run docker compose build (usage: make docker/production/build || make docker/production/build SERVICES="<A> <B>")
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		build ${SERVICES} \
		--parallel

.PHONY: docker/production/up
docker/production/up: ## Run docker compose up without build (usage: make docker/production/up || make docker/productionup SERVICES="<A> <B>")
	@env DOCKER_BUILDKIT=0 docker compose \
		--verbose \
		-f docker-compose.yml \
		up \
		-d ${SERVICES}

##############################
# DATA
##############################
.PHONY: data
data: data/sync ## Run data sync

.PHONY: data/sqlc
data/sqlc: ## Run sqlc generate for all known items
	@cd ./datastore/github_standards && sqlc generate
#--fork-remove-start
	@cd ./datastore/aws_costs && sqlc generate
#--fork-remove-end

.PHONY: data/sync
data/sync: data/sync/vars ## Download all data from bucket
#	github_standards
	@mkdir -p ./builds/api/github_standards/data
	@echo "getting github_standards" && ${AWS_VAULT_COMMAND} aws s3 sync --quiet s3://${AWS_BUCKET}/github_standards ./builds/api/github_standards/data/ && echo bucket_github_standards_done || echo bucket_github_standards_failed;
#--fork-remove-start
#	aws_costs
	@mkdir -p ./builds/api/aws_costs/data
	@echo "getting aws_costs" && ${AWS_VAULT_COMMAND} aws s3 sync --quiet s3://${AWS_BUCKET}/aws_costs ./builds/api/aws_costs/data/ && echo bucket_aws_costs_done || echo bucket_aws_costs_failed;
#--fork-remove-end

.PHONY: data/sync/vars
data/sync/vars: ## Output current values used by data/sync
	@echo "AWS_VAULT_PROFILE: ${AWS_VAULT_PROFILE}"
	@echo "AWS_BUCKET: ${AWS_BUCKET}"
	@echo "AWS_VAULT_COMMAND: ${AWS_VAULT_COMMAND}"
	@echo "SERVICES: ${SERVICES}"

##############################
# CLEANER
##############################
.PHONY: clean
clean: docker/down ## Removes all generated files and docker images
	@rm -f ./servers/api/*.db
	@rm -f ./servers/api/*.csv
	@rm -Rf ./servers/front/govuk
	@rm -Rf ./builds
	@mkdir -p ./builds
	@docker image rm $(images) || echo "ok"
	@docker compose rm api front
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"
##############################
# DEV
##############################
.PHONY: dev
dev: dev/build ## Runs dev/build

.PHONY: dev/run
dev/run: dev/mirror/api ## Runs dev process, mirrors api - TODO: trigger fron servers

.PHONY: dev/run/api
dev/run/api: ## Run the local api server
	@cd ./servers/api && go run main.go

.PHONY: dev/run/front
dev/run/front: ## Run the local fron server
	@cd ./servers/front && go run main.go

.PHONY: dev/build
dev/build: clean ## Builds all local binaries, calls clean first
	@go build -o ./builds/front/front_server ./servers/front/main.go
	@go build -o ./builds/api/api_server ./servers/api/main.go
	@go build -o ./builds/api/seed_cmd ./commands/seed/main.go
	@go build -o ./builds/commands/github_standards ./commands/github_standards/main.go
#--fork-remove-start
	@go build -o ./builds/commands/aws_costs ./commands/aws_costs/main.go
#--fork-remove-end

.PHONY: dev/seed/api
dev/seed/api: ## Seed local databases with data
#	github_standards
	@mkdir -p ./builds/api/github_standards/data
	@cp ./datastore/github_standards/github_standards*.sql ./builds/api/github_standards/
	@echo "seeding github_standards"
	@./builds/api/seed_cmd \
		-table github_standards \
		-db ./builds/api/github_standards.db \
		-schema ./builds/api/github_standards/github_standards.sql \
		-data "./builds/api/github_standards/data/*.json"
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
#--fork-remove-end

.PHONY: dev/mirror
dev/mirror: data/sqlc dev/build data/sync dev/seed/api dev/mirror/api ## Mirrors build setup

.PHONY: dev/mirror/api
dev/mirror/api: ## Mirror position of data files for the api server
	@mv ./builds/api/github_standards.db ./servers/api/github_standards.db
#--fork-remove-start
	@mv ./builds/api/aws_costs.db ./servers/api/aws_costs.db
#--fork-remove-end
