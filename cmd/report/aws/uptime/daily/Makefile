.DEFAULT_GOAL: all
.PHONY: all sample

all:
	@echo "Nothing to run, try using sample"

sample:
	@mkdir -p ./data/
	@echo "[Go](aws_uptime_daily) Running example using Digideps"
	@aws-vault exec digideps-production-breakglass -- go run main.go \
		-account_unit "Digideps" \
		-day "-"
	@ls -larth ./data/
