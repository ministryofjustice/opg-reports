
tests:
	@go clean -testcache
	@clear
	@echo "=== tests"
	@env env CGO_ENABLED=1 LOG_LEVEL="info" go test -count=1 -cover -covermode=atomic ./... && echo "" && echo "passed"
.PHONY: tests
