resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT BGP Alert Rule"
  alert_type                = "BGP"
  expression                = "[(((prefixLengthIPv4 >= 16) && (prefixLengthIPv4 <= 32)) || ((prefixLengthIPv6 >= 32) && (prefixLengthIPv6 <= 128)))]((reachability < 100%) && (changes > 1))"
  minimum_sources           = 1
  rounds_violating_required = 1
  rounds_violating_out_of   = 1
}

resource "thousandeyes_bgp" "test" {
  test_name      = "User Acceptance Test - BGP"
  alerts_enabled = true
  use_public_bgp = true
  prefix         = "192.0.2.0/24"

  alert_rules {
    rule_id = thousandeyes_alert_rule.test.id
  }

  alert_rules {
    rule_id = 921614 #BGP Default Alert Rule
  }
}