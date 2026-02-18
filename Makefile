#========= IMPORT TEAMS =========
IMPORT_DIR ?= ${BUILD_DIR}/import
IMPORT_CMD ?= ${IMPORT_DIR}/import
.PHONY: import-teams
import-teams: CMD_LIST=import
import-teams: get-metadata build-cmds
	@echo "- importing teams "
	@env LOG_LEVEL=info ${IMPORT_CMD} teams \
		--db="${API_DB}" \
		--migration-file="${BUILD_DIR}/migrations.json" \
		--src-file="${METADATA_EX_DIR}/teams.json"

#========= RUN THE API =========
# api command variables
API_DIR ?= ${BUILD_DIR}/api
API_CMD ?= ${API_DIR}/api
API_DB ?= ${API_DIR}/database/api.db
.PHONY: api
api: CMD_LIST=api
api: build-cmds
	@echo "- starting api "
	@env LOG_LEVEL=info ${API_CMD} \
		--db="${API_DB}" \
		--migration-file="${BUILD_DIR}/migrations.json" \
		--api-host="localhost:8081" \

#========= GET the database from s3 =========
DB_BUCKET ?= opg-reports-development
DB_KEY ?= database/api.db
.PHONY: get-db
get-db:
	@echo "- downloading database "

#========= GET opg-metadata release =========
# metadata related variables
METADATA_REPO ?= ministryofjustice/opg-metadata
METADATA_TAG ?= v0.1.29
METADATA_FILE ?= metadata.zip
METADATA_DIR ?= ${BUILD_DIR}/metadata
METADATA_EX_DIR ?= ${METADATA_DIR}/extracted
## Very rarely pulled in, so we can run it from
## make presuming the github cli (gh) instead
## of creating code - commands can take a file
## path param for local development
.PHONY: get-metadata
get-metadata:
	@rm -Rf ${METADATA_DIR}/extracted
	@mkdir -p ${METADATA_DIR}/extracted
	@env GITHUB_TOKEN="${GITHUB_TOKEN}" \
		gh release download ${METADATA_TAG} \
			--clobber \
			--repo ${METADATA_REPO} \
			--dir ${METADATA_DIR} \
			--pattern ${METADATA_FILE}
	@unzip -qq ${METADATA_DIR}/${METADATA_FILE} \
		-d ${METADATA_EX_DIR}


#========= GO BUILDS =========
CMD_DIR = ./report/cmd
# list of commands
CMD_LIST = $(notdir $(wildcard ${CMD_DIR}/*))
# location to put all built files into
BUILD_DIR ?= ./builds
## build all commands based on folder structure within the ./reports/cmd
## directory but allow CMD_LIST changed to make it smarter and allow
## for updating just specific commands
.PHONY: build-cmds
build-cmds:
	@for cmd in ${CMD_LIST}; do \
		echo "- building command [$${cmd}] "; \
		mkdir -p ${BUILD_DIR}/$${cmd}/ ; \
		rm -f ${BUILD_DIR}/$${cmd}/$${cmd} ; \
		go build -ldflags="-w -s" -o ${BUILD_DIR}/$${cmd}/$${cmd} ${CMD_DIR}/$${cmd}; \
	done

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

