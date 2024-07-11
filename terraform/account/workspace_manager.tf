module "workspace_manager" {
  source  = "git@github.com:ministryofjustice/opg-terraform-workspace-manager.git//terraform/workspace_cleanup"
  enabled = local.environment.account_name == "development" ? true : false
}
