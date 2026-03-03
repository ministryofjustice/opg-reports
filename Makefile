LOG_LEVEL ?= info
GITHUBTOKEN ?= ${GH_TOKEN}
#========= IMPORT TEAMS =========
IMPORT_CMD ?= ${BUILD_DIR}/import
.PHONY: import-teams
import-teams: CMD_LIST=import
import-teams: get-metadata build-cmds
	@echo "- importing teams "
	@env LOG_LEVEL=${LOG_LEVEL} ${IMPORT_CMD} teams \
		--db="${API_DB}" \
		--src-file="${METADATA_EX_DIR}/teams.json"

#========= IMPORT ACCOUNTS =========
.PHONY: import-accounts
import-accounts: CMD_LIST=import
import-accounts: get-metadata build-cmds
	@echo "- importing accounts "
	@env LOG_LEVEL=${LOG_LEVEL} ${IMPORT_CMD} accounts \
		--db="${API_DB}" \
		--src-file="${METADATA_EX_DIR}/aws.accounts.json"

#========= IMPORT CODEBASES =========
.PHONY: import-codebases
import-codebases: CMD_LIST=import
import-codebases: build-cmds
	@echo "- importing codebases "
	@env GH_TOKEN="${GITHUBTOKEN}" \
		LOG_LEVEL=${LOG_LEVEL} \
		${IMPORT_CMD} codebases \
		--db="${API_DB}"

#========= IMPORT CODEOWNERS =========
.PHONY: import-codeowners
import-codeowners: CMD_LIST=import
import-codeowners: build-cmds
	@echo "- importing codeowners "
	@env GH_TOKEN="${GITHUBTOKEN}" \
		LOG_LEVEL=${LOG_LEVEL} \
		${IMPORT_CMD} codeowners \
		--db="${API_DB}"

#========= IMPORT UPTIME =========
UPTIME_PROFILE ?= use-production-operator
.PHONY: import-uptime
import-uptime: CMD_LIST=import
import-uptime: build-cmds
	@echo "- importing uptime data "
	@aws-vault exec ${UPTIME_PROFILE} -- \
		env LOG_LEVEL=${LOG_LEVEL} \
		${IMPORT_CMD} uptime \
		--db="${API_DB}"

#========= IMPORT COSTS =========
COSTS_PROFILE ?= use-production-operator
.PHONY: import-costs
import-costs: CMD_LIST=import
import-costs: build-cmds
	@echo "- importing cost data "
	@aws-vault exec ${COSTS_PROFILE} -- \
		env LOG_LEVEL=${LOG_LEVEL} \
		${IMPORT_CMD} costs \
		--db="${API_DB}"


#========= RUN THE API =========
# api command variables
API_DB_DIR ?= ${BUILD_DIR}/database
API_CMD ?= ${BUILD_DIR}/api
API_DB ?= ${API_DB_DIR}/api.db
.PHONY: api
api: CMD_LIST=api
api: build-cmds
	@echo "- starting api "
	@env LOG_LEVEL=${LOG_LEVEL} ${API_CMD} \
		--db="${API_DB}" \
		--api-host="localhost:8081"

#========= RUN THE FRONT END =========
# api command variables
FRONT_CMD ?= ${BUILD_DIR}/front
LOCAL_ASSETS ?= ./web
TEMPLATE_DIR ?= ./report/internal/front/templates
.PHONY: front-assets
front-assets: get-govuk
	@echo "- copying templates and local assets"
	@rm -Rf ${BUILD_DIR}/templates
	@rm -Rf ${BUILD_DIR}/web
	@cp -r ${LOCAL_ASSETS} ${BUILD_DIR}
	@cp -r ${TEMPLATE_DIR} ${BUILD_DIR}

.PHONY: front
front: CMD_LIST=front
front: build-cmds front-assets
	@echo "- starting front "
	@env LOG_LEVEL=${LOG_LEVEL} ${FRONT_CMD} \
		--api-host="localhost:8081" \
		--front-host="localhost:8080" \
		--root-dir="${BUILD_DIR}" \
		--govuk-version="${GOVUK_TAG}"

#========= GET the database from s3 =========
GET_DB_BUCKET ?= opg-reports-development
GET_DB_PROFILE ?= shared-development-operator
MIGRATE_CMD ?= ${BUILD_DIR}/migrate
.PHONY: get-db
get-db: CMD_LIST=migrate
get-db: build-cmds
	@rm -Rf ${API_DB}
	@mkdir -p ${API_DB_DIR}
	@echo "- downloading database (${GET_DB_BUCKET})"
	@aws-vault exec ${GET_DB_PROFILE} -- aws s3 cp \
		--sse AES256 \
    	s3://${GET_DB_BUCKET}/database/api.db \
    	${API_DB_DIR}/
	@echo "- migrating & covnerting database "
	@env LOG_LEVEL=${LOG_LEVEL} ${MIGRATE_CMD} \
		--db="${API_DB}" \
		--convert

# .PHONY: upload-db
# upload-db:
# 	aws-vault exec ${GET_DB_PROFILE} -- aws s3 cp \
# 		--sse AES256 \
#     	${API_DB_DIR}/api.db \
# 		s3://${GET_DB_BUCKET}/database/api.db


#========= GET opg-metadata release =========
## Very rarely pulled in, so we can run it from
## make instead of creating code - commands can
## take a file path param for local development
##
## presumed gh client installed
METADATA_REPO ?= ministryofjustice/opg-metadata
METADATA_TAG ?= v0.1.29
METADATA_ZIP ?= metadata.zip
METADATA_EX_DIR ?= ${BUILD_DIR}/metadata-extracted
.PHONY: get-metadata
get-metadata:
	@rm -Rf ${METADATA_EX_DIR}
	@mkdir -p ${METADATA_EX_DIR}
	@env GH_TOKEN="${GITHUBTOKEN}" \
		gh release download ${METADATA_TAG} \
			--clobber \
			--repo ${METADATA_REPO} \
			--dir ${BUILD_DIR} \
			--pattern ${METADATA_ZIP}
	@unzip -qq ${BUILD_DIR}/${METADATA_ZIP} \
		-d ${METADATA_EX_DIR}
	@rm -f ${BUILD_DIR}/${METADATA_ZIP}
#========= GET gov-uk release =========
## Run during the build process; will
## download marked release from govuk
## and setup folder structure to work
## for front end
##
## presumed gh client installed
GOVUK_REPO ?= alphagov/govuk-frontend
GOVUK_TAG ?= 5.14.0
GOVUK_ZIP ?= release-v${GOVUK_TAG}.zip
GOVUK_EX_DIR ?= ${BUILD_DIR}/govuk
.PHONY: get-govuk
get-govuk:
	@rm -Rf ${GOVUK_EX_DIR}
	@mkdir -p ${GOVUK_EX_DIR}
	@env GH_TOKEN="${GITHUBTOKEN}" \
		gh release download v${GOVUK_TAG} \
			--clobber \
			--repo ${GOVUK_REPO} \
			--dir ${BUILD_DIR} \
			--pattern ${GOVUK_ZIP}
	@unzip -o -qq ${BUILD_DIR}/${GOVUK_ZIP} \
		-d ${GOVUK_EX_DIR}
	@rm -f ${BUILD_DIR}/${GOVUK_ZIP}



#========= GO BUILDS =========
BUILD_DIR ?= ./builds
## creates the dir
.PHONY: build-prep
build-prep:
	rm -Rf ${BUILD_DIR}/
	mkdir -p ${BUILD_DIR}/

## build all commands based on folder structure
## within the ./reports/cmd directory but allow
## CMD_LIST changed to make it smarter and allow
## for updating just specific commands
CMD_DIR = ./report/cmd
CMD_LIST = $(notdir $(wildcard ${CMD_DIR}/*))
.PHONY: build-cmds
build-cmds:
	@for cmd in ${CMD_LIST}; do \
		echo "- building command [$${cmd}] "; \
		mkdir -p ${BUILD_DIR}/ ; \
		rm -f ${BUILD_DIR}/$${cmd} ; \
		go build -ldflags="-w -s" -o ${BUILD_DIR}/$${cmd} ${CMD_DIR}/$${cmd}; \
	done

##========= DOCKER =========
SERVICES ?= front api
VERBOSE ?=
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
docker-up: CMD_LIST=front api
docker-up: build-cmds get-govuk docker-build
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
		LOG_LEVEL="${LOG_LEVEL}" \
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
		LOG_LEVEL="${LOG_LEVEL}" \
		GITHUB_TOKEN="${GITHUB_TOKEN}" \
		AWS_REGION="${AWS_REGION}" \
		AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}" \
		AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
		AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
		AWS_SESSION_TOKEN="${AWS_SESSION_TOKEN}" \
		go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover ./...
	@go tool cover -html=code-coverage.out



