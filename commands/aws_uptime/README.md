# AWS Uptime

This command returns uptime date for the .

### Generation

```bash
aws-vault exec <aws-profile> -- go run main.go \
    -unit "<unit>" \
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
