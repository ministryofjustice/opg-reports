SERVICES ?= api


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
		go test -count=1 -cover -covermode=atomic ./... && echo "" && echo "passed"
	@echo "==="
.PHONY: tests

test:
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
		go test -count=1 -cover -covermode=atomic ./... -run="$(name)" && echo "" && echo "passed"
	@echo "==="
.PHONY: tests

## Run the go code coverage tool
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
.PHONY: coverage

#========= DOCKER =========
## Build local development version of the docker image
docker/build:
	@env DOCKER_BUILDKIT=1 \
	docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		build ${SERVICES}
.PHONY: docker/build

## Build and run the local docker images
docker/up: docker/build
	@env DOCKER_BUILDKIT=1 \
	docker compose \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		up \
		-d ${SERVICES}
.PHONY: docker/up

## Clean any old docker images out
docker/clean: docker/down
	@docker image rm $(images) || echo "ok"
	@env DOCKER_BUILDKIT=1 \
	docker compose \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		rm ${SERVICES}
	@docker container prune -f
	@docker image prune -f --filter="dangling=true"
.PHONY: docker/clean

## run docker compose down, turning off all docker containers
docker/down:
	@env DOCKER_BUILDKIT=1 \
	docker compose \
		--verbose \
		-f docker-compose.yml \
		-f docker-compose.dev.yml \
		down
.PHONY: docker/down
