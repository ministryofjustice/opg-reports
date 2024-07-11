module "workspace_manager" {
  source  = "git@github.com:ministryofjustice/opg-terraform-workspace-manager.git//terraform/workspace_cleanup?ref=0.3.2"
  enabled = local.environment.account_name == "development" ? true : false
}
