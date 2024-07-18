.DEFAULT_GOAL: all
.PHONY: all sample

all:
	@echo "Nothing to run, try using sample"

sample:
	@echo "[Go](aws_cost_monthly) Running example using Digideps Dev"
	@aws-vault exec digideps-development-operator -- go run main.go \
		-account_id 248804316466 \
		-account_name "Digideps Development" \
		-account_environment "Development" \
		-account_unit "Digideps" \
		-account_label "Digideps" \
		-month "-"
	@ls -larth ./data/
