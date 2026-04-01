data "thousandeyes_agent" "frankfurt" {
  agent_name = "Frankfurt, Germany"
}

data "thousandeyes_alert_rule" "def_alert_rule1" {
  rule_name = "Default HTTP Alert Rule 2.0"
}

data "thousandeyes_alert_rule" "def_alert_rule2" {
  rule_name = "Default Network Alert Rule 2.0"
}

resource "thousandeyes_http_server" "test" {
  test_name      = "UAT - Tag Assignment Update - HTTP"
  interval       = 120
  alerts_enabled = true
  url            = "https://www.thousandeyes.com"
  agents         = [data.thousandeyes_agent.frankfurt.agent_id]
  alert_rules    = [data.thousandeyes_alert_rule.def_alert_rule1.id]
}

resource "thousandeyes_agent_to_server" "test" {
  test_name      = "UAT - Tag Assignment Update - Agent To Server"
  interval       = 120
  alerts_enabled = true
  server         = "api.stg.thousandeyes.com"
  protocol       = "tcp"
  probe_mode     = "sack"
  agents         = [data.thousandeyes_agent.frankfurt.agent_id]
  alert_rules    = [data.thousandeyes_alert_rule.def_alert_rule2.id]
}

resource "thousandeyes_tag" "tag1" {
  key         = "UAT TagAssignUpdate Tag1 Key"
  value       = "UAT TagAssignUpdate Tag1 Value"
  object_type = "test"
  color       = "#b3de69"
  access_type = "all"
  icon        = "LABEL"
}

resource "thousandeyes_tag" "tag2" {
  key         = "UAT TagAssignUpdate Tag2 Key"
  value       = "UAT TagAssignUpdate Tag2 Value"
  object_type = "test"
  color       = "#fdb462"
  access_type = "all"
  icon        = "LABEL"
}

resource "thousandeyes_tag_assignment" "assign1" {
  tag_id = thousandeyes_tag.tag1.id
  assignments {
    id   = thousandeyes_http_server.test.id
    type = "test"
  }
}
