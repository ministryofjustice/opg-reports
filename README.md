# opg-reports

This repository acts as a central hub to generate, combine and display series of key data that we report on frequently to both internal and external parties.

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


# Forking<a name="forking"></a>

While this code bae make use of a lot of configuration, the default and checked in versions are set for OPG versions. Therefore, for most other users it is sensible to fork this repository to make your reguired changes.

We'll try to walk you though all the changes you need to make after forking, including infrastructure / access needs that are used but not included in this repository.

In order to start capturing github repository standards for your organisation you will need to focus on the following areas fo code:

- Remove / disable reports that won't be used
- CI user for terraform
- Using this projects terraform
- S3 bucket
- Role for local development
- Role or user for workflow access
- Repository secrets
- Repository standards report workflow
- ECR

## Current state

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


## Remove or disable reports<a name="forking-remove-reports"></a>

By default, this repository includes a series of reports that may not be relevant for your needs. The guide presumes you will be only using the [github repository standards report](./cmd/report/github/standards/README.md) and its [workflow](./.github/workflows/report_repository_standards.yml).

You can leave the other `go` reports, but please do either remove the `on` triggers or delete the files for the other workflows starting with `report_`.

## CI user for terraform<a name="forking-ci-user"></a>

If you are planning on using the terraform within this you will need to create a new role / user for your accounts. One you have, please find all references to our role `docs-and-metadata-ci` and change that.

Should be in the following files:
- [PR workflow](./.github/workflows/workflow_pr.yml)
- [Path to live workflow](./.github/workflows/workflow_path_to_live.yml)
- [terraform account default role](./terraform/account/variables.tf)
- [terraform environment default role](./terraform/environment/variables.tf)
- [s3 access policy](./terraform/environment/s3.tf)

## Using this projects terraform<a name="forking-this-terraform"></a>

The terraform included in this project is split between `account` and `environment` levels and utilises 2 fixed terraform workspaces - `development` and `production`.

Within the [terraform folders](./terraform/) you will find tfvars files that are mapped based upon the workspace name.

If you want to use this terraform you will need to update the state file locations that are configured as well as the ci user it runs with ([see above](#forking-ci-user)).

If you want to configure your terraform elsewhere, you will need to remove sections of the ci/cd workflows relating to `terraform`. Please check [workflow_path_to_live.yml](./.github/workflows/workflow_path_to_live.yml) and [workflow_pr.yml](./.github/workflows/workflow_pr.yml) and remove suitable references, all jobs will start with `terraform_` prefix.

*Note: We're very early in the process (< v1), so for the moment we are only utilising development workspace.*

## S3 bucket<a name="forking-bucket"></a>

An S3 bucket will need to be created (via this projects terraform or elsewhere); once done so you will need to update all references to the bucket in the code base to reflect the new names.

For development:
- [`Makefile`](./Makefile) - change the value of the variable `BUCKET`
- [`report_repository_standards.yml`](./.github/workflows/report_repository_standards.yml) - in the `upload_to_s3` step, change the `bucket` value

## Role for local development<a name="forking-local-dev-access"></a>

OPG apply have an operator role for developers that covers the permissions they need to work within their team. You can see in our [`s3.tf`](./terraform/environment/s3.tf) the operator role is being granted permissions for our bucket in the `allow_data_role_access` policy.

Local development assumes the use of `aws-vault` and mutliple profiles within a persons `/.aws/config` file. This is reflected within the `Makefile` when downloading assets from the s3 bucket in the build process, swapping to a specific operator profile to do so.

You will need to update the AWS profile (`shared-development-operator`) used within the code base to match how your uses can access the data.

- [`Makefile`](./Makefile) - change the value of `AWS_PROFILE` variable to one suitable for your team.

## Role or user for workflow access<a name="forking-workflow-role-access"></a>

Our workflows use OIDC roles to upload their result data to and s3 bucket. This role (`opg-reports-github-actions-s3`) is defined within [our s3 terraform](./terraform/environment/s3.tf) and will need to be replaced with a suitable you have create role that has the following permissions:

```
s3:GetObject
s3:PutObject
s3:ListBucket
```

You will need to update the following files with the new role:

- [`report_repository_standards.yml`](./.github/workflows/report_repository_standards.yml) - in the `configure_aws_creds_s3_upload` step, change the role to assume
  - If you are not using an OIDC role, you will need to expand this and add secrets etc to authenticate


## Repository secrets<a name="forking-secrets"></a>

The repository uses a series of secrets for various forms of access and permissions.

### AWS secrets<a name="forking-secrets-aws"></a>

For running terraform as a CI user, the repository has that users keys stored as secrets:

- `AWS_ACCESS_KEY_ID_ACTIONS`
- `AWS_SECRET_ACCESS_KEY_ACTIONS`

These are configured outside of this repositories terraform code, so will need to be replaced and matched to the [CI user you create](#forking-ci-user).

### GitHub Token<a name="forking-secrets-github"></a>

The github repository standards report needs a token that can access public as well as private repositories for your team. this is stored in a secret:

- `GH_ORG_ACCESS_TOKEN`

Again, this is configured in terraform elsewhere, so please do replace this value with a suitable one when forking.

### SSH<a name="forking-secrets-ssh"></a>

The terraform code within this repository calls external modules, so we store a secret containing a private ssh key to allow access:

- `SSH_PRIVATE_KEY_EXTERNAL_MODULES`

Replace when forking to allow the terraform to run, but if you are not using this projects terraform it can be ignored.


## Repository standards reporting workflow

This [workflow](./.github/workflows/report_repository_standards.yml) calls [`go` code](./cmd/report/github/standards/README.md) to go and fetch data.

The go code expects arguments to be either an organisation and team or the name of a repository:

- `-organisation <org-slug> -team <team-slug>`
- `-repo <full-name-slug>`

Within the workflow you will need to replace the env variable `arguments` with versions that work for your team.

The workflow is currently configured to run every Saturday at 9am - feel free to adjust for your own needs.
