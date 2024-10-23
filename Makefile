.PHONY: tests test coverage openapi sdk
## run all tests
tests:
	@go clean -testcache
	@clear
	@echo "============== tests =============="
	@env env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -cover -covermode=atomic ./... && echo "" && echo "passed ✅"

## run go suite tests that match the file pattern
## usage:
## `make tests name=<pattern>`
test:
	@go clean -testcache
	@clear
	@echo "============== test: [$(name)] =============="
	@env CGO_ENABLED=1 GITHUB_ACCESS_TOKEN="${GITHUB_TOKEN}" LOG_LEVEL="info" LOG_TO="stdout" go test -count=1 -v ./... -run="$(name)"

## Run the go code coverage tool
coverage:
	@rm -Rf ./code-coverage.out
	@clear
	@echo "============== coverage =============="
	@env CGO_ENABLED=1 LOG_LEVEL="warn" LOG_TO="stdout" go test -count=1 -covermode=count -coverprofile=code-coverage.out -cover ./...
	@go tool cover -html=code-coverage.out


## Output the open api spec
openapi:
	@env CGO_ENABLED=1 go run ./api/main.go openapi > openapi.yaml

## Generate local SDK from the spec using oapi-codegen v2.4.1 (pinned to sha)
sdk:
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@9c09ef9e9d4be639bd3feff31ff2c06961421272
	@mkdir -p ./sdk
	@oapi-codegen -generate "types,client" -package sdk openapi.yaml > sdk/sdk.go
	@go mod tidy

