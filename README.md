# opg-reports

This repository acts as a central hub to generate, combine and display series of key data that we report on frequently to both internal and external parties.

- [Current functionality](#functionality)
- [Forking](#forking)
- [Detailed breakdown](#structure)

# Current functionality<a name="functionality"></a>

At the moment, the CI/CD workflows (with `workflow_` prefix) within this project run on pull_requests or pushes to main. The workflow does the following actions:

- Determines terraform versions to use
- Creates a semver tag
- Runs the go test suites
- Creates docker image
  - Fetches assets (from `s3` and the govuk front end) using the makefile (`make assets`)
  - Login to ECR
  - Builds and pushes containers (both the `api` and `front` found with ./services/ folder)
- Runs terraform (both ./terraform/account and ./terraform/environment) for the development workspace
- Creates a release (or pre-release) using the semver tag generated, attaching the go code binaries as an artifact

# Forking<a name="forking"></a>

While this code base makes use of a lot of configuration, the default and checked in versions are set for OPGs requirements. Therefore, for most other users it is sensible to fork this repository to make your required changes.

## Automated

To make this process as simple as we can, we have included a [bash script](./scripts/fork.sh) that should be run after you have forked the original repository to update the settings.

Run this by calling from the root of the forked repository:

```bash
./scripts/fork.sh
```

The script will ask for following details:

- Name of business unit (`--business-unit`)
- Name of the *DEVELOPMENT* S3 bucket used for data storage (`--development-bucket-name`)
- The *DEVELOPMENT* aws profile for local S3 download  (`--aws-profile`)
- The *DEVELOPMENT* role ARN to use for *DOWNLOADING* from the S3 bucket in workflow (`--development-bucket-download-arn`)
- The *DEVELOPMENT* OIDC role ARN to use for *UPLOADING* to the S3 bucket in workflows (`--development-bucket-upload-arn`)
- The AWS ECR registry id (`--ecr-registry-id`)
- The *DEVELOPMENT* OIDC role ARN to use for pushing to ECR (`--ecr-login-push-arn`)
- The github organisation slug (`--gh-org`)
- The github team slug (`--gh-team`)

You can also pass these values as cli arguments and the script will use the passed value instead of asking, like:

```bash
./scripts/fork.sh \
  --business-unit OPG \
  --gh-team opg \
  --aws-profile operator
```

It will then go over the code base and update all the relevent files - you should then be able to commit those changes and have a functional version.

You can use the argument names

*Note:* Assumes infrastructure will be handled external to this repository.

## Manual

Please review [the guide](./MANUAL_FORK_GUIDE.md) for how to fork this repository in more detail.

# Structure and organisation<a name="structure">

The repository has 3 core areas, data gathering, api layer and a display layer.

## Data gathering<a name="data-gathering-intro"></a>

All code responsible for generating data is based within the `./cmd/` folder path, mostly within `./cmd/report/` folder.

Currently, this repository has reports to fetch the following data:

- [Github repository standards](./cmd/report/github/standards/README.md)
- [AWS monthly costs](./cmd/report/aws/cost/monthly/README.md)

Each report is run via a [github workflow](./.github/workflows/README.md#report-workflows) using a series of arguments to fetch and then generate a single [data file](./cmd/report/README.md#filename-pattern). This data file is then uploaded to a shared s3 bucket for use by the api layer.

Each report folder contains a `Makefile` with examples of using the command.

Please see the [report readme](./cmd/report/README.md) for more details how reports are configured and utilised.

## API layer<a name="api-layer-intro"></a>

The [api layer](./services/api/README.md) creates a single webserver and delgates the handling of various endpoints to `go` code within sub-folders.

The folder paths of the handlers [should match the report-path](./cmd/report/README.md#report-path) for consistency.

The api layer contains a `./data/` directory which is where the s3 bucket is synch'd into (via a make command).

## Display layer<a name="display-layer-intro"></a>

The [display / front layer](./services/front/README.md) uses `go` webserver and its built in templating to generate html views. It has two request handlers - static and dynamic.

The static handler handle pages that are simple html / markup that don't utilise any data from the api.

The dynamic handler fetches information from the api and will then call the template stack to generate the output using teh api data.

This is configured within the [front's config.json file](./services/front/config.json).

### Configuration<a name="display-layer-configuration"></a>

The front uses its config file to set the visiable organisation name rendered in the html (default: OPG) as well as the navigation structure with where each page gets its data from.

Additionally, what repository values are checked for the baseline and extended standard checks are configured here.

Please see the [package details for more info](./services/front/cnf/cnf.go)


## Running locally<a name="running-locally"></a>

We utilise `Makefile`s to provide a simple and consistently available method to share common tasks for the project.

The [`Makefile`](./Makefile) at the project root contains targets to build and spin up the project for local development to make the process easier.

Currently, neither method for running the project have have realtime code updates to a running process.

### Test suite<a name="running-locally-tests"></a>

While not complete, there is a good amount of tests within the go code that you can run to check any changes against. These are setup in the `Makefile` as targets.

Run all tests:

```
make tests
```

Run a specifc set of tests based on their name:

```
make test name="<pattern>"
```

Run code coverage checks:

```
make coverage
```

### With Docker<a name="running-locally-with-docker"></a>

This project is setup to use `docker compose` files ([docker-compose.yml](./docker-compose.yml) and [docker-compose.dev.yml](./docker/docker-compose.dev.yml) ), with the base file referencing the `latest` images from our private registry.

Assuming you have access, you will to login to ECR using the below command:

```
aws-vault exec management-operator -- aws ecr get-login-password --region eu-west-1 | docker login --username AWS --password-stdin 311462405659.dkr.ecr.eu-west-1.amazonaws.com
```

Once you have authenticated, you can run the latest images directly by calling:

```
make up
```

Otherwise, you can build and run the images in development setup by running:

```
make dev
```

Note: This will fetch data from s3 and the gov uk front end assets and will trigger further authentication calls to AWS.

### Without docker<a name="running-locally-without-docker"></a>

This project is built using `go` ([see here for version](./go.mod)) and both the api and front are run directly from the go binaries. This means you can run these directly if you wish to instead buildign the docker images.

You will need to fetch the data stored in s3 and the govuk front assets before running the code by calling:

```
make assets
```

You can then run the api in the foreground by calling:

```
make dev-run-api
```

And then, in a different terminal, the front end:

```
make dev-run-front
```

The api will then be visible on `http://localhost:8081` and the front end on `http://localhost:8080`.
