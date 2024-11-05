# Github Standards

This command returns information about all repositories in github org and team.

### Generation

```bash
env GITHUB_ACCESS_TOKEN=${GITHUB_TOKEN} go run main.go \
    -organsiation="<github-org-slug>" \
    -team="<github-team-slug>"
```

### Upload

Upload the generated file by running:

```bash
aws-vault exec <aws-profile> -- aws s3 cp \
	--recursive ./data s3://<bucket>/github_standards \
	--sse AES256
```


### Download

```bash
aws-vault exec <aws-profile> -- aws s3 sync s3://<bucket>/github_standards ./bucket-data
```
