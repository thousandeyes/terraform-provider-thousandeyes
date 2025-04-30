data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_agent" "arg_frankfurt" {
  agent_name = "Frankfurt, Germany"
}

resource "thousandeyes_account_group" "example_account_group" {
  account_group_name = "Example Account Group"
  agents             = [data.thousandeyes_agent.arg_amsterdam.agent_id, data.thousandeyes_agent.arg_frankfurt.agent_id]
}
