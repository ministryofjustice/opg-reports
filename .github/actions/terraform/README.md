# Terraform plan and apply

Run a plan using the terraform version passed along and the from the directory using the workspace set.

Output the plan data to a local file.

When `apply` is set to `true`, then run the apply using the out from the plan step.

Both the plan and apply uses state lock timeout so if there is an error or clash the state file will unlock after 300 seconds.

## Usage

This is typical version of the how to use the action to apply the terraform within a folder.

Normally you would get the workspace and version values dynamically.

```
- name: "Account level terraform"
  uses: ./.github/actions/terraform
  with:
    apply: true
    directory: ./terraform/account
    workspace: "default"
    version: "1.9.1"

```
