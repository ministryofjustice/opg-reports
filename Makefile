SERVICES ?= api front

API_BUILD = ./builds/api
DBP = ${API_BUILD}/databases/api.db
FRONT_BUILD =./builds/front
VERBOSE ?= # --verbose

tests:
	@go clean -testcache
	@clear
	@echo "=== tests"
	@env CGO_ENABLED=1 \
		LOG_LEVEL="WARN" \
		GITHUB_TOKEN="${GITHUB_TOKEN}" \
		AWS_REGION="${AWS_REGION}" \
		AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}" \
		AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
		AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
		AWS_SESSION_TOKEN="${AWS_SESSION_TOKEN}" \
		go test -count=1 -cover -covermode=atomic ./... && echo "" && echo "passed"
	@echo "==="
.PHONY: tests

test:
	@go clean -testcache
	@clear
	@echo "=== test: $(name)"
	@env CGO_ENABLED=1 \
		LOG_LEVEL="WARN" \
		GITHUB_TOKEN="${GITHUB_TOKEN}" \
		AWS_REGION="${AWS_REGION}" \
		AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}" \
		AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
		AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
		AWS_SESSION_TOKEN="${AWS_SESSION_TOKEN}" \
		go test -count=1 -cover -covermode=atomic ./... -run="$(name)" && echo "" && echo "passed"
	@echo "==="
.PHONY: tests

## Run the go code coverage tool
coverage:
	@rm -Rf ./code-coverage.out
	@clear
	@echo "=== coverage"
	@env CGO_ENABLED=1 \
		LOG_LEVEL="WARN" \
		GITHUB_TOKEN="${GITHUB_TOKEN}" \
		AWS_REGION="${AWS_REGION}" \
		AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}" \
		AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
		AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
		AWS_SESSION_TOKEN="${AWS_SESSION_TOKEN}" \
		go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover ./...
	@go tool cover -html=code-coverage.out
.PHONY: coverage

#========= LOCAL =========
.PHONY: local/build/api
local/build/api:
	@rm -Rf ${API_BUILD}
	@mkdir -p ${API_BUILD} ${API_BUILD}/databases ${API_BUILD}/bin
	@go build -o ${API_BUILD}/bin/db ./report/cmd/db/
	@aws-vault exec shared-development-operator -- \
   		env DATABASE_PATH=${DBP} \
   		${API_BUILD}/bin/db download
	@go build -o ${API_BUILD}/bin/api ./report/cmd/api

.PHONY: local/build/front
local/build/front:
	@rm -Rf ${FRONT_BUILD}
	@mkdir -p ${FRONT_BUILD} ${FRONT_BUILD}/bin
	@go build -o ${FRONT_BUILD}/bin/govuk ./report/cmd/govuk/
	@env SERVERS_FRONT_DIRECTORY=${FRONT_BUILD} \
  		env GITHUB_TOKEN=${GITHUB_TOKEN} \
  		${FRONT_BUILD}/bin/govuk frontend
	@go build -o ${FRONT_BUILD}/bin/front ./report/cmd/front/
	@cp -r ./report/cmd/front/templates ${FRONT_BUILD}/
	@cp -r ./report/cmd/front/local-assets ${FRONT_BUILD}/

.PHONY: local/build
local/build: local/build/api local/build/front

#========= DOCKER =========
## Build local development version of the docker image
docker/build:
	@env DOCKER_BUILDKIT=1 \
	docker compose ${VERBOSE} \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		build ${SERVICES}
.PHONY: docker/build

## Build and run the local docker images
docker/up: local/build docker/clean docker/build
	@env DOCKER_BUILDKIT=1 \
	docker compose ${VERBOSE} \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		up \
		-d ${SERVICES}
.PHONY: docker/up

## Clean any old docker images out
docker/clean: docker/down
	@docker image rm $(shell docker images -a | grep 'opg-reports/*' | awk '{print $$1":"$$2}') || echo "ok"
	@env DOCKER_BUILDKIT=1 \
	docker compose ${VERBOSE} \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		rm ${SERVICES}
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"
.PHONY: docker/clean

## run docker compose down, turning off all docker containers
docker/down:
	@env DOCKER_BUILDKIT=1 \
	docker compose ${VERBOSE} \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		down
.PHONY: docker/down
