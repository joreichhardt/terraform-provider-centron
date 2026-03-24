terraform {
  required_providers {
    centron = {
      source = "registry.terraform.io/joreichhardt/centron"
    }
  }
}

provider "centron" {
  # Add client_id and client_secret here, or set CENTRON_CLIENT_ID / CENTRON_CLIENT_SECRET environment variables.
}

resource "centron_server" "test" {
  pool        = "pool_aw3_ssd"
  hostname    = "test-vm-tf"
  description = "Provisioned via Terraform"
  cores       = 1
  memory      = 1024
  disks       = [20]
  image       = "ubuntu-22"
  type        = "unmanaged"
}
