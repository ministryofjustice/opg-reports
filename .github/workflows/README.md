# GitHub Workflows

All workflows follow the same naming convention, please make sure to check this when creating any additional workflows.

For a workflow that runs as part of the [CI/CD pipeline](#pipeline-workflows), please name it following this pattern:

```
workflow_<optional_qualifier>_<pr|path_to_live>.yml
```

For the main pull request trigger that would look like:

```
workflow_pr.yml
```

But for a workflow that runs on a pr, but its not the primary task - say running document site generation, add some additional info to the name, such as:

```
workflow_generate_docs_pr.yml
```

For a workflow that runs a report / data generation command, the please follow the below convention:

```
report_<source>_<purpose>.yml
```



## Use of composite actions

Where sensible, keep the main workflows cleaner by utilising composite actions [within this repository](../actions/README.md), the shared actions from [opg-github-actions](https://github.com/ministryofjustice/opg-github-actions/) or other known / trusted sources.

## "Pipeline" workflows

"Pipeline" workflows are those that actually deploy the terraform and generate releases of this code base making use of github actions.

As the reporting tool has a limited scope and used for internal reporting only it runs with a simplified `development` and `production` only fixed environments, we do not utilise ephemeral environments or additional workspaces.

These are foucused on testing the code base and infrastructure to ensure released versions function correctly and always use `pull_request` and `push` triggers.

## Report workflows<a name="report-workflows"></a>

Each report workflow uses the checked out state of this repository to build the go binaries (note: this needs adjusting to fetch that latest released version) and then run the specific report.

The result of that report is then uploaded to an s3 bucket for storage and will be fetched and used in the next build of the [api](../../servers/api/README.md)

To help reduce the reptition of the code, there are composite actions to handle both steps ([`report`](../actions/report/README.md) and [`s3_upload`](../actions/s3_upload/README.md) ) that handle the tasks for you.

Please make sure the jobs have any authentication requirements setup before calling the action.
