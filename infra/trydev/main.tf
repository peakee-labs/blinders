provider "aws" {
  region  = var.region
  profile = var.profile
}

terraform {
  backend "s3" {
    key = "blinders-trydev-state"
  }
}

module "core" {
  source = "../core"

  project = {
    name           = "blinders-try"
    environment    = "dev"
    default_region = "ap-south-1"
  }

  domains = {
    http      = "trydev.api.peakee.co"
    websocket = "trydev.ws.peakee.co"
  }

  env_filename = "../../.env.dev"
}

