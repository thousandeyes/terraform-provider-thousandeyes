data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_agent" "arg_frankfurt" {
  agent_name = "Frankfurt, Germany"
}

resource "thousandeyes_account_group" "test" {
  account_group_name        = "User Acceptance Test - Account Group (Updated)"
  agents           =  [data.thousandeyes_agent.arg_amsterdam.agent_id]
}
