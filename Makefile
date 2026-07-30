

#========= TESTS =========
## Run all tests but without tokens - so should only run mocks
.PHONY: tests
tests:
	@go clean -testcache
	@clear
	@echo "=== tests"
	@env CGO_ENABLED=1 \
		LOG_LEVEL="INFO" \
		GITHUB_TOKEN="" \
		GH_TOKEN="" \
		AWS_REGION="${AWS_REGION}" \
		AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}" \
		AWS_ACCESS_KEY_ID="" \
		AWS_SECRET_ACCESS_KEY="" \
		AWS_SESSION_TOKEN="" \
		go test -count=1 -cover -covermode=atomic ./... && echo "" && echo "passed ✅" || echo "failed ❌"
	@echo "==="


## Run a test but with all tokens passed so will try to get real data
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
