module "network" {
  source = "./modules/network"

  cidr = "10.1.0.0/16"
  tags = local.default_tags
}