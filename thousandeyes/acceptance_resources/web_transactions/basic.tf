data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT Web Transactions Alert Rule"
  alert_type                = "web-transactions"
  expression                = "((webPages((webTxResponseTime >= 100 ms) && (webTxPageLoadError != \"\") && (webTxOnLoadTime >= 200 ms))))"
  rounds_violating_out_of   = 1
  rounds_violating_required = 1
  minimum_sources           = 1
}

resource "thousandeyes_web_transaction" "test" {
  test_name          = "User Acceptance Test - Web Transactions"
  interval           = 120
  alerts_enabled     = true
  url                = "https://www.thousandeyes.com"
  use_public_bgp     = true
  transaction_script = <<EOF
  import { By, Key, until } from 'selenium-webdriver'; 
  import { driver, markers, credentials, downloads, transaction, test } from 'thousandeyes'; 
  runScript(); 
  async function runScript() 
  { const settings = test.getSettings();
  // Load page
  await driver.get(settings.url);
  await driver.wait(until.titleIs('Digital Experience Monitoring | ThousandEyes'), 1000);
  await driver.takeScreenshot();
};
EOF

  agents      = [data.thousandeyes_agent.arg_amsterdam.agent_id]
  alert_rules = [thousandeyes_alert_rule.test.id, "921619"]
}
