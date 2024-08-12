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
	@echo "============== test =============="
	@env CGO_ENABLED=1 GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -race -count=1 -v ./... -run="$(name)"

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
	@echo " SLOW PROCESS - LARGE DATABASE SEEDS"
	@env LOG_LEVEL="INFO" LOG_TO="stdout" go test -count=1  -v ./... -bench=. -run=xxx -benchmem


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
	@cp ./datastore/github_standards/schema.sql ./builds/go/github_standards.sql
# @ls -lth ./builds/go/
	@./builds/go/seeder_cmd -which all -dir ./builds/go
# @ls -lth ./builds/go/dbs

go-run-api:
	@cd ./builds/go/ && ./api_server

