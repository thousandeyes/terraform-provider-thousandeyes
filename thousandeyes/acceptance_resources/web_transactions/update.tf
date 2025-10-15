data "thousandeyes_agent" "arg_ny" {
  agent_name = "New York, NY"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default Web Transaction Alert Rule 2.0"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT Web Transactions Alert Rule (Updated)"
  alert_type                = "web-transactions"
  expression                = "((webPages((webTxResponseTime >= 100 ms) && (webTxPageLoadError != \"\") && (webTxOnLoadTime >= 200 ms))))"
  rounds_violating_out_of   = 3
  rounds_violating_required = 3
  minimum_sources           = 1
}

resource "thousandeyes_web_transaction" "test" {
  test_name          = "User Acceptance Test - Web Transactions (Updated)"
  interval           = 300
  alerts_enabled     = true
  url                = "https://www.thousandeyes.com"
  emulated_device_id = "1"
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

  agents      = [data.thousandeyes_agent.arg_ny.agent_id]
  alert_rules = [thousandeyes_alert_rule.test.id, data.thousandeyes_alert_rule.def_alert_rule.id]
}
