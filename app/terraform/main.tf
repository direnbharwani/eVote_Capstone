module "terraform_state_backend" {
  source = "./modules/terraform_state"
}

module "evote_poc_storage" {
  source = "./modules/storage"
}