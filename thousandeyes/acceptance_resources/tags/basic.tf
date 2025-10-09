data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_alert_rule" "def_alert_rule1" {
  rule_name = "Default HTTP Alert Rule 2.0"
}

data "thousandeyes_alert_rule" "def_alert_rule2" {
  rule_name = "Default Network Alert Rule 2.0"
}

resource "thousandeyes_http_server" "test" {
  test_name      = "New User Acceptance Test - HTTP with tags"
  interval       = 120
  alerts_enabled = true
  url            = "https://www.thousandeyes.com"
  agents         = [data.thousandeyes_agent.amsterdam.agent_id]
  alert_rules    = [data.thousandeyes_alert_rule.def_alert_rule1.id]
}

resource "thousandeyes_agent_to_server" "test" {
  test_name      = "New User Acceptance Test - Agent To Server with tags"
  interval       = 120
  alerts_enabled = true
  server         = "api.stg.thousandeyes.com"
  protocol       = "tcp"
  probe_mode     = "sack"
  agents         = [data.thousandeyes_agent.amsterdam.agent_id]
  alert_rules    = [data.thousandeyes_alert_rule.def_alert_rule2.id]
}

resource "thousandeyes_tag" "tag1" {
  key         = "UAT Tag1 Key"
  value       = "UAT Tag1 Value"
  object_type = "test"
  color       = "#b3de69"
  access_type = "all"
  icon        = "LABEL"
}

resource "thousandeyes_tag" "tag2" {
  key         = "UAT Tag2 Key"
  value       = "UAT Tag2 Value"
  object_type = "test"
  color       = "#fdb462"
  access_type = "all"
  icon        = "LABEL"
}

resource "thousandeyes_tag" "tag3" {
  key         = "UAT Tag3 Key"
  value       = "UAT Tag3 Value"
  object_type = "test"
  color       = "#8dd3c7"
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

resource "thousandeyes_tag_assignment" "assign2" {
  tag_id = thousandeyes_tag.tag2.id
  assignments {
    id   = thousandeyes_agent_to_server.test.id
    type = "test"
  }
}

resource "thousandeyes_tag_assignment" "assign3" {
  tag_id = thousandeyes_tag.tag3.id
  assignments {
    id   = thousandeyes_http_server.test.id
    type = "test"
  }
  assignments {
    id   = thousandeyes_agent_to_server.test.id
    type = "test"
  }
}

