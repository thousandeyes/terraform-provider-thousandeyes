data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT DNS Server Alert Rule"
  alert_type                = "dns-server"
  expression                = "((probDetail != \"\") && (Auto(delay >= Medium sensitivity)))"
  minimum_sources           = 1
  rounds_violating_required = 1
  rounds_violating_out_of   = 1
}

resource "thousandeyes_dns_server" "test" {
  test_name      = "User Acceptance Test - DNS Server"
  interval       = 120
  alerts_enabled = true
  domain         = "thousandeyes.com A"
  agents         = [data.thousandeyes_agent.amsterdam.agent_id]
  alert_rules    = ["921612", thousandeyes_alert_rule.test.id]
  dns_servers    = ["ns-cloud-d1.googledomains.com", "ns-1458.awsdns-54.org", "ns-597.awsdns-10.net"]
}
