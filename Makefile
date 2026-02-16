SERVICES ?= api front
# SOURCE DIRECTORIES
SRC_CMD_DIR = ./report/cmd
CMD_LIST = $(notdir $(wildcard ${SRC_CMD_DIR}/*))
SRC_FRONT_DIR = ./report/cmd/front
# BUILT LOCATIONS
BUILT_ROOT = ./builds
## api locations
BUILT_API_DIR = ${BUILT_ROOT}/api
BUILT_API_DB_PATH = ${BUILT_API_DIR}/database/api.db
BUILT_API_CMD = ${BUILT_API_DIR}/api
## db locations
BUILT_DB_CMD = ${BUILT_ROOT}/db/db
## front
BUILT_FRONT_DIR = ${BUILT_ROOT}/front/
BUILT_FRONT_CMD = ${BUILT_ROOT}/front/front
## govuk related settings
BUILT_GOVUK_CMD = ${BUILT_ROOT}/govuk/govuk

#========= LOCAL =========
## build all commands based on folder structure within the ./reports/cmd
## directory but allow CMD_LIST changed to make it smarter and allow
## for updating just specific commands
.PHONY: build-cmds
build-cmds:
	@for cmd in ${CMD_LIST}; do \
		echo "- building command [$${cmd}] "; \
		mkdir -p ${BUILT_ROOT}/$${cmd}/ ; \
		rm -f ${BUILT_ROOT}/$${cmd}/$${cmd} ; \
		go build -ldflags="-w -s" -o ${BUILT_ROOT}/$${cmd}/$${cmd} ${SRC_CMD_DIR}/$${cmd}; \
	done

## download the development version of the database
db_bucket = opg-reports-development
db_key = database/api.db
.PHONY: db-download
db-download: CMD_LIST=db
db-download: build-cmds
	@echo "- downloading development database with aws-vault"
	@rm -f ${BUILT_API_DB_PATH}
	@aws-vault exec shared-development-operator -- \
		${BUILT_DB_CMD} download \
			--bucket="${db_bucket}" \
			--key="${db_key}" \
			--directory="${BUILT_API_DIR}"

## seed and migrate a local database at fixed path
.PHONY: db-seed
db-seed: CMD_LIST=db
db-seed: build-cmds
	@echo "- seeding local database [${BUILT_API_DB_PATH}]"
	@env LOG_LEVEL=error ${BUILT_DB_CMD} seed --db="${BUILT_API_DB_PATH}"

## migrate a local database at fixed path
.PHONY: db-migrate
db-migrate: CMD_LIST=db
db-migrate: build-cmds
	@echo "- migrating local database [${BUILT_API_DB_PATH}]"
	@env LOG_LEVEL=error ${BUILT_DB_CMD} migrate --db="${BUILT_API_DB_PATH}"

## run the api from the local ./build folder structure
.PHONY: api
api: CMD_LIST=api
api: build-cmds
	@echo "- starting api "
	@env LOG_LEVEL=info ${BUILT_API_CMD} \
		--db="${BUILT_API_DB_PATH}" \
		--address="localhost:8081"

## run the front end
.PHONY: front
front: CMD_LIST=govuk front
front: build-cmds govuk
	@echo "- copying templates and local assets"
	@cp -r ${SRC_FRONT_DIR}/templates ${BUILT_FRONT_DIR}/
	@cp -r ${SRC_FRONT_DIR}/local-assets ${BUILT_FRONT_DIR}/

	@env LOG_LEVEL=info ${BUILT_FRONT_CMD} \
		--root-dir="${BUILT_FRONT_DIR}" \
		--api="localhost:8081" \
		--address="localhost:8080"

# build and download the govuk front end assets
.PHONY: govuk
govuk: CMD_LIST=govuk
govuk: build-cmds
	@rm -Rf ${BUILT_FRONT_DIR}/govuk
	@echo "- downloading govuk assets"
	@env LOG_LEVEL=ERROR GITHUB_TOKEN=${GITHUB_TOKEN} \
		${BUILT_GOVUK_CMD} --directory="${BUILT_FRONT_DIR}/govuk"


# #========= DOCKER =========
## Clean any old docker images out
.PHONY: docker-clean
docker-clean: docker-down
	@docker image rm $(shell docker images -a | grep 'opg-reports/*' | awk '{print $$1":"$$2}') || echo "ok"
	@env DOCKER_BUILDKIT=1 \
	docker compose ${VERBOSE} \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		rm ${SERVICES}
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"
.PHONY: docker-clean

.PHONY: docker-down
docker-down:
	@env DOCKER_BUILDKIT=1 \
	docker compose ${VERBOSE} \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		down

## Build local development version of the docker image
.PHONY: docker-build
docker-build:
	@env DOCKER_BUILDKIT=1 \
	docker compose ${VERBOSE} \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		build ${SERVICES}


## Build and run the local docker images
.PHONY: docker-up
docker-up: build-cmds  docker-build
	@env DOCKER_BUILDKIT=1 \
	docker compose ${VERBOSE} \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		up \
		-d ${SERVICES}


#========= TESTS =========
## Run all tests
.PHONY: tests
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
		go test -count=1 -cover -covermode=atomic ./... && echo "" && echo "passed ✅" || echo "failed ❌"
	@echo "==="

## Run specific test via named param
.PHONY: test
test:
	@go clean -testcache
	@clear
	@echo "=== test: $(name)"
	@env CGO_ENABLED=1 \
		LOG_LEVEL="INFO" \
		GITHUB_TOKEN="${GITHUB_TOKEN}" \
		GH_TOKEN="${GITHUB_TOKEN}" \
		AWS_REGION="${AWS_REGION}" \
		AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}" \
		AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
		AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
		AWS_SESSION_TOKEN="${AWS_SESSION_TOKEN}" \
		go test -count=1 -cover -covermode=atomic ./... -run="$(name)" && echo "" && echo "passed ✅" || echo "failed ❌"
	@echo "==="


## Run the go code coverage tool
.PHONY: coverage
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

