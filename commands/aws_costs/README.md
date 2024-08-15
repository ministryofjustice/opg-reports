# AWS Costs

This command returns cost data for the month provided at a daily granularity.

### Generation

```bash
aws-vault exec digideps-development-operator -- go run main.go \
    -id 248804316466 \
    -name "Digideps test" \
    -environment "test" \
    -unit "Digideps" \
    -label "Digideps" \
    -month "-"
```

### Upload

Upload the generated file by running:

```bash
aws-vault exec shared-development-operator -- aws s3 cp \
	--recursive ./data s3://report-data-development/aws_costs \
	--sse AES256
```


### Download

```bash
aws-vault exec shared-development-operator -- aws s3 sync s3://report-data-development/aws_costs ./bucket-data
```
