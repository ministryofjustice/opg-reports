# Report

Build the local `go` code and then run the binary whose name is passed along as the `cmd` input.

The outputed variable `data_folder` can then be used by `s3_upload` action to push data to the s3 bucket.

Please ensure any authentication (AWS creds, github tokens etc) are done before calling this action within your job.


## Usage

Here is a typical usage of the action where the report requires the environment variable set in this step as well:

```
- name: "Run report"
  id: run
  uses: ./.github/actions/report
  env:
    GITHUB_ACCESS_TOKEN: ${{ secrets.GH_ORG_ACCESS_TOKEN }}
  with:
    name: "Repository standards"
    cmd: "github_standards"
    arguments: '-organisation ministryofjustice -team opg'

```
