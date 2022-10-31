terraform {
  required_providers {
    thousandeyes = {
      source = "thousandeyes/thousandeyes"
      version = ">= 1.3.1"
    }
  }
}

provider "thousandeyes" {
  token            = "your-token"
  account_group_id = "your-account-id"
}

data "thousandeyes_agent" "arg_cordoba" {
  agent_name = "Cordoba, Argentina"
}

resource "thousandeyes_http_server" "www_thousandeyes_http_test" {
  test_name      = "Example HTTP test set from Terraform provider"
  interval       = 120
  alerts_enabled = false

  url = "https://www.thousandeyes.com"

  agents {
    agent_id = data.thousandeyes_agent.arg_cordoba.agent_id
  }
}
