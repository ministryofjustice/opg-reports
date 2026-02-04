SERVICES ?= api front

# CMD_BUILD = ./builds/cmd
# API_BUILD = ./builds/api
# DB_BUILD = ./builds/databases
# DBP ?= ${DB_BUILD}/api.db
# FRONT_BUILD = ./builds/front

# SOURCE DIRECTORIES
SRC_CMDS = ./report/cmd
CMD_LIST = $(notdir $(wildcard ${SRC_CMDS}/*))
SRC_FRONT_DIR = ./report/cmd/front
# BUILT LOCATIONS
BUILT_ROOT = ./builds
## api locations
BUILT_API_DB_DIR = ${BUILT_ROOT}/api/databases
BUILT_API_DB_PATH = ${BUILT_API_DB_DIR}/api.db
BUILT_API_CMD = ${BUILT_ROOT}/api/api
## database downloader tool
BUILT_DB_DOWNLOADER = ${BUILT_ROOT}/db/db
## front
BUILT_FRONT_DIR = ${BUILT_ROOT}/front/
BUILT_FRONT_CMD = ${BUILT_ROOT}/front/front
## govuk related settings
BUILT_GOVUK_CMD = ${BUILT_ROOT}/govuk/govuk


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
.PHONY: tests
test:
	@go clean -testcache
	@clear
	@echo "=== test: $(name)"
	@env CGO_ENABLED=1 \
		LOG_LEVEL="WARN" \
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

#========= LOCAL BUILDS =========
## build all commands based on folder structure within the ./reports/cmd
## directory
.PHONY: build-cmds
build-cmds:
	@for cmd in ${CMD_LIST}; do \
		echo "- building command [$${cmd}] "; \
		mkdir -p ${BUILT_ROOT}/$${cmd}/ ; \
		rm -f ${BUILT_ROOT}/$${cmd}/$${cmd} ; \
		go build -ldflags="-w -s" -o ${BUILT_ROOT}/$${cmd}/$${cmd} ${SRC_CMDS}/$${cmd}; \
	done

#========= LOCAL SETUP =========

## download the development version of the database
.PHONY: db-download
db-download: CMD_LIST="db"
db-download: build-cmds
	@echo "- downloading development database with aws-vault"
	@mkdir -p ${BUILT_API_DB_DIR}
	@rm -f ${BUILT_API_DB_PATH}
	@aws-vault exec shared-development-operator -- \
		env DATABASE_PATH=${BUILT_API_DB_PATH} \
		DATABASE_BUCKET_NAME="opg-reports-development" \
		${BUILT_DB_DOWNLOADER} download

## run the api from the local ./build folder structure
.PHONY: api
api: CMD_LIST="api"
api: build-cmds
	@echo "- starting api binary"
	@env DATABASE_PATH=${BUILT_API_DB_PATH} \
	SERVERS_API_ADDR="localhost:8081" \
	SERVERS_FRONT_ADDR="localhost:8080" \
		${BUILT_API_CMD}

## run the front from the local ./build folders and setup templates
## and govuk assets as well
.PHONY: front
front: CMD_LIST="front"
front: build-cmds
	@rm -Rf ${BUILT_FRONT_DIR}/govuk
	@echo "- downloading govuk assets"
	@env LOG_LEVEL=ERROR GITHUB_TOKEN=${GITHUB_TOKEN} \
		${BUILT_GOVUK_CMD} --directory="${BUILT_FRONT_DIR}/govuk"

	@echo "- copying templates and local assets"
	@cp -r ${SRC_FRONT_DIR}/templates ${BUILT_FRONT_DIR}/
	@cp -r ${SRC_FRONT_DIR}/local-assets ${BUILT_FRONT_DIR}/

	@echo "- starting front server"
	@env SERVERS_API_ADDR="localhost:8081" \
		SERVERS_FRONT_ADDR="localhost:8080" \
		SERVERS_FRONT_DIRECTORY=${BUILT_FRONT_DIR} \
			${BUILT_FRONT_CMD}


# #========= DOCKER =========
# ## Build local development version of the docker image
# docker/build:
# 	@env DOCKER_BUILDKIT=1 \
# 	docker compose ${VERBOSE} \
# 		-f docker-compose.yml \
# 		-f docker-compose.dev.yml \
# 		build ${SERVICES}
# .PHONY: docker/build

# ## Build and run the local docker images
# docker/up: local/build docker/clean docker/build
# 	@env DOCKER_BUILDKIT=1 \
# 	docker compose ${VERBOSE} \
# 		-f docker-compose.yml \
# 		-f docker-compose.dev.yml \
# 		up \
# 		-d ${SERVICES}
# .PHONY: docker/up

# ## Clean any old docker images out
# docker/clean: docker/down
# 	@docker image rm $(shell docker images -a | grep 'opg-reports/*' | awk '{print $$1":"$$2}') || echo "ok"
# 	@env DOCKER_BUILDKIT=1 \
# 	docker compose ${VERBOSE} \
# 		-f docker-compose.yml \
# 		-f docker-compose.dev.yml \
# 		rm ${SERVICES}
# 	@docker container prune -f
# 	@docker image prune -f --filter="dangling=true"
# .PHONY: docker/clean

# ## run docker compose down, turning off all docker containers
# docker/down:
# 	@env DOCKER_BUILDKIT=1 \
# 	docker compose ${VERBOSE} \
# 		-f docker-compose.yml \
# 		-f docker-compose.dev.yml \
# 		down
# .PHONY: docker/down
