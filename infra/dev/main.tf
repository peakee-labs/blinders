provider "aws" {
  region  = var.region
  profile = var.profile
}

terraform {
  backend "s3" {
    key = "blinders-dev-state"
  }
}

module "core" {
  source = "../core"

  project = {
    name           = "blinders"
    environment    = "dev"
    default_region = "ap-south-1"
  }

  domains = {
    http      = "dev.api.peakee.co"
    websocket = "dev.ws.peakee.co"
  }

  env_filename = "../../.env.dev"
}
