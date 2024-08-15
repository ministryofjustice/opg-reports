# Forking this repository

While there is an automated script, we also try to provide detailed breakdown on how to manaully make the required changes to this repository.

We'll try to walk you though all the changes you need to make after forking, including infrastructure / access needs that are used but not included in this repository.

Please test your changes as you go - see the [running locally section on how](./README.md#running-locally).

In order to start capturing github repository standards for your organisation you will need to focus on the following areas fo code:

- [Remove / disable reports that won't be used](#forking-remove-reports)
- [CI user for terraform](#forking-ci-user)
- [Using this projects terraform](#forking-this-terraform)
- [S3 bucket](#forking-bucket)
- [Role for local development](#forking-local-dev-access)
- [Role or user for workflow access](#forking-workflow-role-access)
- [Repository secrets](#forking-secrets)
- [Repository standards report workflow](#forking-standards-workflow)
- [ECR](#forking-ecr)
- [Front end configuration](#forking-front-config)

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

This project makes use of OIDC roles to allow the workflows to pull / push data to sources such as s3 and ecr.

## ECR<a name="forking-workflow-role-access-ecr"></a>

*Only relevant if you are using this projects terraform.*

Within the projects workflows there is an ecr specifc OIDC role (`opg-reports-github-actions-ecr-push`) configured to push/pull docker images with ECR. Its replacement will require the following permissions:

```
ecr:CompleteLayerUpload
ecr:UploadLayerPart
ecr:InitiateLayerUpload
ecr:BatchCheckLayerAvailability
ecr:PutImage
ecr:BatchGetImage
ecr:GetDownloadUrlForLayer
```

You will need to update all references:

- [`workflow_path_to_live.yml`](./.github/workflows/workflow_path_to_live.yml)
- [`workflow_pr.yml`](./.github/workflows/workflow_pr.yml)

## S3<a name="forking-workflow-role-access-s3"></a>

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

## AWS secrets<a name="forking-secrets-aws"></a>

For running terraform as a CI user, the repository has that users keys stored as secrets:

- `AWS_ACCESS_KEY_ID_ACTIONS`
- `AWS_SECRET_ACCESS_KEY_ACTIONS`

These are configured outside of this repositories terraform code, so will need to be replaced and matched to the [CI user you create](#forking-ci-user).

## GitHub Token<a name="forking-secrets-github"></a>

The github repository standards report needs a token that can access public as well as private repositories for your team. this is stored in a secret:

- `GH_ORG_ACCESS_TOKEN`

Again, this is configured in terraform elsewhere, so please do replace this value with a suitable one when forking.

## SSH<a name="forking-secrets-ssh"></a>

The terraform code within this repository calls external modules, so we store a secret containing a private ssh key to allow access:

- `SSH_PRIVATE_KEY_EXTERNAL_MODULES`

Replace when forking to allow the terraform to run, but if you are not using this projects terraform it can be ignored.


## Repository standards reporting workflow<a name="forking-standards-workflow"></a>

This [workflow](./.github/workflows/report_repository_standards.yml) calls [`go` code](./cmd/report/github/standards/README.md) to go and fetch data.

The go code expects arguments to be either an organisation and team or the name of a repository:

- `-organisation <org-slug> -team <team-slug>`
- `-repo <full-name-slug>`

Within the workflow you will need to replace the env variable `arguments` with versions that work for your team.

The workflow is currently configured to run every Saturday at 9am - feel free to adjust for your own needs.

## ECR<a name="forking-ecr"></a>

During the docker build process we utilise a private registry (AWS ECR) within one of our accounts. As the built version of the images contains a copy of all data from the s3 bucket, it is advisable keep using a private registry.

You will need to update the code in the following places to change the registry:

- [`docker-compose.yml`](./docker-compose.yml)


## Front end configuration<a name="forking-front-config"></a>

The front end (display layer) service utilises a configuration file to setup the navigation, organisation name and which pages use api data or static content.

Which configuration file to use is controlled by an environment variable - `CONFIG_FILE` - for the docker container. By default, this is set to the [`config.json`](./services/front/config.json) version, which is in turn a symlink to [`config.opg.json`](./services/front/config.opg.json).

We recommend you create a copy of [`config.simple.json`](./services/front/config.simple.json) to start with and rename it something relevant (`config.<business-unit>.json`). You can then either change the symlink to point to this new file, or change the environment varibles in the following places:

- [`docker-compose.yml`](./docker-compose.yml) - change the env var (`CONFIG_FILE`) value
- [`Makefile`](./Makefile) - in the `dev-run-front` target, swap the value of the `CONFIG_FILE` variable.

The Dockerfile will copy files matching `config*.json` pattern building the build.

For more details on how the configuration file is used, please see its [package details](./services/front/cnf/cnf.go).

## Organisation

The front end uses the `organsiation` value from the [configuration file](./services/front/config.json) in the html title and various page headings - please adjust this to something suitable for yourself.

## Sections

The `sections` block is a recursive structure that is used to render the navigation levels of the front end. The [simple example included](./services/front/config.simple.json) shows a single level of navigation with one main page displaying the github standards report details.

You are unlikey to need to change this.

## Standards

The properties that are checked for the baseline and extended tests are configured in this block. This are the normal values and unlikely they need changing.
