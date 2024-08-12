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

##############################
# GO BUILD
# - build all go binaries at once and push to ./builds/go/
#   using goreleaser
##############################
go-build:
	@goreleaser build --clean --single-target --skip=validate
	@rm -f ./builds/binaries/*.json
	@rm -f ./builds/binaries/*.yml


go-run-api: go-build
	@cd ./builds/api/ && ./api_server

go-run-front: go-build
	@cd ./builds/front/ && ./front_server

