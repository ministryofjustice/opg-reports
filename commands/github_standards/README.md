# GitHub Standards

This command returns all the repositories for the organsiation and team set that the `GITHUB_ACCESS_TOKEN` has access to. These are then converted into local struct and written to scsv file.

### Generation

The command will generate a csv file (`./github_standards/github_standards.csv`):

```bash
env GITHUB_ACCESS_TOKEN=${GITHUB_TOKEN} go run main.go \
	-organisation ministryofjustice \
	-team opg
```

### Upload

Upload the generated file by running:

```bash
aws-vault exec shared-development-operator -- aws s3 cp \
	--recursive ./github_standards s3://report-data-development/github_standards \
	--sse AES256
```


### Download

```bash
aws-vault exec shared-development-operator -- aws s3 sync s3://report-data-development/github_standards ./bucket-data
```
