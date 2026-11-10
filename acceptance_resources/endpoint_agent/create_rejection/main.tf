terraform {
  required_providers {
    thousandeyes = {
      source = "thousandeyes/thousandeyes"
    }
  }
}

provider "thousandeyes" {}

resource "thousandeyes_endpoint_agent" "unimported" {
  name = "Terraform create rejection check"
}
