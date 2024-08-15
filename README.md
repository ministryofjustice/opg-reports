# opg-reports

This repository acts as a central hub to generate, combine and display series of key data that we report on frequently to both internal and external parties.

## Quick startup

As the codebase will auto-generate databases if they are not present, you can get a local version up and running in docker by:

```bash
make clean && make up
```

You can then view the site at:

- [front](http://localhost:8080)
- [api](http://localhost:8081)


## Forking

This repository comes with a bash script to update values, so afte ryou ahve forked please run:

```bash
./scripts/fork.sh
```

The script will ask for following details:

- Name of business unit (`--business-unit`)
- Name of the DEVELOPMENT S3 bucket used for data storage (`--development-bucket-name`)
- The DEVELOPMENT aws profile for local S3 download (`--aws-profile`)
- The DEVELOPMENT role ARN to use for DOWNLOADING from the S3 bucket in workflow (`--development-bucket-download-arn`)
- The DEVELOPMENT OIDC role ARN to use for UPLOADING to the S3 bucket in workflows (`--development-bucket-upload-arn`)
- The AWS ECR registry id (`--ecr-registry-id`)
- The DEVELOPMENT OIDC role ARN to use for pushing to ECR (`--ecr-login-push-arn`)
- The github organisation slug (`--gh-org`)
- The github team slug (`--gh-team`)

You can also pass these values as cli arguments and the script will use the passed value instead of asking, like:

```bash
./scripts/fork.sh \
  --business-unit OPG \
  --gh-team opg \
  --aws-profile operator
```

After this completes, you can re-run the app and see the updates ([see here](#quick-startup)).

Please then commit any changes made into your fork.

If you want to make changes manually, there is a more [detailed guide available here](./docs/usage/manual-forking-guide.md).

*Note:* Assumes infrastructure will be handled external to this repository.


## Docs


- [Development enviroment](./docs/usage/development-environment.md)
