.DEFAULT_GOAL: all
.PHONY: test tests benchmarks coverage go-build demo-data

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
	@go clean -testcache
	@clear
	@echo "============== coverage =============="
	@env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover -v ./...
	@go tool cover -html=code-coverage.out

benchmarks:
	@go clean -testcache
	@clear
	@echo "============== benchmarks =============="
	@echo " WARNING: CAN BE SLOW"
	@env LOG_LEVEL="warn" LOG_TO="stdout" go test -v ./... -bench=. -run=xxx -benchmem -benchtime=10s

benchmark:
	@go clean -testcache
	@clear
	@echo "============== benchmark: [$(name)] =============="
	@echo " WARNING: CAN BE SLOW"
	@env LOG_LEVEL="info" LOG_TO="stdout" go test -v ./... -bench=$(name) -run=xxx -benchmem -benchtime=10s

##############################
# FRONT END ASSETS
##############################
GOVUK_FRONT_VERSION := "v5.4.0"
GOVUK_DOWNLOAD_FOLDER := ./builds/govuk-frontend

assets-front:
	@echo "-----"
	@echo "[Assets](front) Building..."
	@echo "	source: [alphagov/govuk-frontend@${GOVUK_FRONT_VERSION}]"
	@echo "	target: [./servers//front/assets/]"
	@rm -Rf ${GOVUK_DOWNLOAD_FOLDER}
	@rm -Rf ./servers//front/assets/css/
	@rm -Rf ./servers//front/assets/fonts/
	@rm -Rf ./servers//front/assets/images/
	@rm -Rf ./servers//front/assets/manifest.json
	@mkdir -p ${GOVUK_DOWNLOAD_FOLDER}
	@cd ${GOVUK_DOWNLOAD_FOLDER} && gh release download ${GOVUK_FRONT_VERSION} -R alphagov/govuk-frontend
	@cd ${GOVUK_DOWNLOAD_FOLDER} && unzip -qq release-${GOVUK_FRONT_VERSION}.zip
	@cd ${GOVUK_DOWNLOAD_FOLDER} && mkdir -p ./assets/css/ && mv govuk-frontend-*.css* ./assets/css/
	@mv ${GOVUK_DOWNLOAD_FOLDER}/assets/css/ ./servers//front/assets/
	@mv ${GOVUK_DOWNLOAD_FOLDER}/assets/fonts/ ./servers//front/assets/
	@mv ${GOVUK_DOWNLOAD_FOLDER}/assets/images/ ./servers//front/assets/
	@mv ${GOVUK_DOWNLOAD_FOLDER}/assets/manifest.json ./servers//front/assets/
	@rm -Rf ${GOVUK_DOWNLOAD_FOLDER}
	@echo "[Assets](front) Built."
##############################
# SQLC
##############################
# trigger sqlc generate on all datastore sections
sqlc:
	@cd ./datastore/github_standards/ && sqlc generate

##############################
# GO BUILD
# - build all go binaries at once and push to ./builds/go/
#   using goreleaser
##############################
go-build: sqlc
	@goreleaser build --clean --single-target --skip=validate
	@mkdir -p ./builds/go/dbs
	@rm -f ./builds/go/*.json
	@rm -f ./builds/go/*.yaml
# 	copy assets over to the build directory
	@cp ./datastore/github_standards/schema.sql ./builds/go/dbs/github_standards.sql
	@cp ./servers/front/config.json ./builds/go/
# @ls -lth ./builds/go/
	@./builds/go/seeder_cmd -which all -dir ./builds/go
# @ls -lth ./builds/go/dbs

go-run-api:
	@if [ ! -f "./builds/go/api" ]; then \
		make go-build; \
	fi
	@cd ./builds/go/ && ./api

go-run-front:
	@if [ ! -f "./builds/go/front" ]; then \
		make go-build; \
	fi
	@cd ./builds/go/ && ./front

