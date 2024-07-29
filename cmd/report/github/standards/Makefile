.DEFAULT_GOAL: all
.PHONY: all sample

all:
	@echo "Nothing to run, try using sample"

sample:
	@echo "[Go](github_standards) Running example fetching everything for a team"
	@env GITHUB_ACCESS_TOKEN=${GITHUB_TOKEN} go run main.go \
		-organisation ministryofjustice \
		-team opg
	@echo "[Go](github_standards) Running example fetching just one repository"
	@env GITHUB_ACCESS_TOKEN=${GITHUB_TOKEN} go run main.go \
		-repo ministryofjustice/opg-aws-whitepaper
