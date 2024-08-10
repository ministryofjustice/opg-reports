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
	@go clean -testcache
	@clear && env GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -race -count=1 -v ./... -run="$(name)"

tests:
	@go clean -testcache
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -race -count=1 -cover -covermode=atomic -v ./...

coverage:
	@rm -Rf ./code-coverage.out
	@go clean -testcache
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover -v ./...
	@go tool cover -html=code-coverage.out

benchmarks:
	@go clean -testcache
	@clear && env LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1  -v ./... -bench=. -run=xxx -benchmem


##############################
# SQLC
##############################
# trigger sqlc generate on all datastore sections
sqlc:
	@cd ./datastore/github_standards/ && sqlc generate

##############################
# DATABASE SETUPS
##############################
DB="./builds/dbs/github_standards.db"
CSV="./builds/csv/github_standards/github_standards.csv"
SCHEMA="./datastore/github_standards/schema.sql"
TABLE="github_standards"
# run db import for github standards
db-github-standards:
	@mkdir -p ./builds/dbs
	@if [ ! -f "${DB}" ]; then \
		sqlite3 "${DB}" "VACUUM;" ; \
		sqlite3 "${DB}" < "${SCHEMA}"  ; \
		sqlite3 "${DB}" ".import --csv ${CSV} ${TABLE}" ; \
	fi
	@sqlite3 "${DB}" "VACUUM;"

dbs: sqlc data-generator db-github-standards
	@ls -lth "./builds/dbs/"
##############################
# DATA
# - generate fake data in the build dir if not present
# - import present files into sqlite
##############################
data-generator:
	@./builds/go/demo_data_cmd -which all -dir ./builds/csv

##############################
# GO BUILD
# - build all go binaries at once and push to ./builds/go/
#   using goreleaser
##############################
go-build: sqlc
	@goreleaser build --clean --single-target --skip=validate
	@rm -f ./builds/go/*.json
	@rm -f ./builds/go/*.yaml

