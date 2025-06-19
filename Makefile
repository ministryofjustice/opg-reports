
tests:
	@go clean -testcache
	@clear
	@echo "=== tests"
	@env CGO_ENABLED=1 LOG_LEVEL="info" GH_TOKEN="${GITHUB_TOKEN}" go test -count=1 -cover -covermode=atomic ./... && echo "" && echo "passed"
.PHONY: tests
