# Reports

All data for the reporting website is generated via the commands within this folder.

## Assets<a name="assets"></a>

A report requires multiple assets to function, they are normally:

1. `go` code to fetch the data
2. [Github workflow](../../.github/workflows/README.md#report-workflows) to run the go code on a set interval / trigger
3. [S3 bucket](../../.github/workflows/README.md) and path to store outputed data
4. [A series of api end point handlers](../../services/api/README.md) to read the data and filter / transform it
5. [A navigation item setup in front end services](../../services/front/README.md) config.json to display the data

This README addresses the 1st point in detail and provides some details for the 2nd and 3rd, but please read their respective README files for more details.


## Structure<a name="structure"></a>

The `./cmd/report` folder is organised by the source of the data (aws / github etc), then the type of data (costs / standards), matching the pattern below:

```
./cmd/report/<provider-name>/<report-category>/<optional-time-period-group>/
```

This would look like the below for AWS uptime data:

```
./cmd/report/aws/uptime/daily/
```

<a name="report-path"></a>
The `<provider-name>/<report-category>/<optional-time-period-group>/` segment is referred to as `report-path`

Each command is built as a single binary from a `main.go` using the common build process [`Makefile`](../../Makefile) to work. Any new reports must have a new target for them added to that `Makefile` - see existing targets with `go-report-*` prefix.

Each report should be scoped to a singular entity and be called as many times as required for the number of instances you have. For example, to generate costs for AWS, run once per account and time period.

Each report binary is built within the github workflow, under the `go_test_and_build` job. The binary will be included in the release artifact for the workflow.

When creating your new report command, please include a `Makefile` within the same directory with a `sample` target showing typical usage.

Please see an existing report for more detail.

## Report output<a name="report-output"></a>

Each report should output its resulting data into a single json file with a unique name based on a combination of its inputs names and values (using `^` and `.` as seperators), located in a sub folder `./data/`.

<a name="filename-pattern"></a>
```
./cmd/report/aws/uptime/daily/data/<argument-a>^<value-1>.<argument-b>^<value-2>.json
```

Therefore a report that fetched AWS uptime status being run with `--month=2024-01 --account=test` arguments should produce a filename like:

```
./cmd/report/aws/uptime/daily/data/account^test.month=2024-01.json
```

When run from a [github workflow](../../.github/workflows/README.md), the `./data/` directory will be uploaded to the s3 bucket with a path matching the [report-path](#report-path)

```
 s3://<bucket-name>/aws/uptime/daily/
```

The s3 bucket upload will overwrite any file with the same name, so please ensure filenames are unique and versioning is enabled on the bucket.

**Note:** The entire bucket is downloaded as part of the [API image build process](../../services/api/README.md) and included within the docker image. Please check your docker image registery visiblity settings to ensure any sensitive data is not being exposed.


### Access

Please make sure to restrict any access / permissions being used to the minimal amount required to run your workflow. This includes making use if OIDC roles for AWS access and limiting the permissions of the workflow itself.
