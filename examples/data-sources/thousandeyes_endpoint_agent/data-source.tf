data "thousandeyes_endpoint_agent" "laptop" {
  agent_name = "Alice's laptop"
}

output "endpoint_agent_id" {
  value = data.thousandeyes_endpoint_agent.laptop.agent_id
}

output "endpoint_agent_name" {
  value = data.thousandeyes_endpoint_agent.laptop.agent_name
}
