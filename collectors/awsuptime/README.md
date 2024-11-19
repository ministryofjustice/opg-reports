# AWS Uptime

This command returns uptime data.

### Generation

```bash
aws-vault exec <aws-profile> -- go run main.go \
    -unit "<unit>" \
    -id "<aws-account-id>" \
    -day "-"
```

### Upload

Upload the generated file by running:

```bash
aws-vault exec <aws-profile> -- aws s3 cp \
	--recursive ./data s3://<bucket>/aws_uptime \
	--sse AES256
```


### Download

```bash
aws-vault exec <aws-profile> -- aws s3 sync s3://<bucket>/aws_uptime ./bucket-data
```
