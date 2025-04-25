resource "thousandeyes_web_transaction" "example_web_transaction_test" {
  test_name          = "Example Web Transaction test set from Terraform provider"
  interval           = 120
  alerts_enabled     = false
  url                = "https://www.thousandeyes.com"
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
  agents             = ["3"] # Singapore
} 
