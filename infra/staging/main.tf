
module "core" {
  source = "../core"

  project = {
    name           = "blinders"
    environment    = "staging"
    default_region = "ap-south-1"
  }

  domains = {
    http      = "staging.api.peakee.co"
    websocket = "staging.ws.peakee.co"
  }

  env_filename = "../../.env.staging"
}

