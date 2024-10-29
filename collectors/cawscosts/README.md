# AWS Costs

This command returns cost data for the month provided at a daily granularity.

### Generation

```bash
aws-vault exec <aws-profile> -- go run main.go \
    -id <account-id> \
    -name "<name>" \
    -environment "development" \
    -unit "<unit>" \
    -label "<label>" \
    -month "-"
```

### Upload

Upload the generated file by running:

```bash
aws-vault exec <aws-profile> -- aws s3 cp \
	--recursive ./data s3://<bucket>/aws_costs \
	--sse AES256
```


### Download

```bash
aws-vault exec <aws-profile> -- aws s3 sync s3://<bucket>/aws_costs ./bucket-data
```
