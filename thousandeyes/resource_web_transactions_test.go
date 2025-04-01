package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestAccThousandEyesWebTransactions(t *testing.T) {
	var httpResourceName = "thousandeyes_web_transaction.test"
	var testCases = []struct {
		name                 string
		createResourceFile   string
		updateResourceFile   string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
		checkCreateFunc      []resource.TestCheckFunc
		checkUpdateFunc      []resource.TestCheckFunc
	}{
		{
			name:                 "create_update_delete_web_transactions_test",
			createResourceFile:   "acceptance_resources/web_transactions/basic.tf",
			updateResourceFile:   "acceptance_resources/web_transactions/update.tf",
			resourceName:         httpResourceName,
			checkDestroyFunction: testAccCheckDefaultWebTransactionsResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(httpResourceName, "test_name", "User Acceptance Test - Web Transactions"),
				resource.TestCheckResourceAttr(httpResourceName, "interval", "120"),
				resource.TestCheckResourceAttr(httpResourceName, "transaction_script", "  import { By, Key, until } from 'selenium-webdriver'; \n  import { driver, markers, credentials, downloads, transaction, test } from 'thousandeyes'; \n  runScript(); \n  async function runScript() \n  { const settings = test.getSettings();\n  // Load page\n  await driver.get(settings.url);\n  await driver.wait(until.titleIs('Digital Experience Monitoring | ThousandEyes'), 1000);\n  await driver.takeScreenshot();\n};\n"),
				resource.TestCheckResourceAttr(httpResourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(httpResourceName, "alert_rules.#", "2"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(httpResourceName, "test_name", "User Acceptance Test - Web Transactions (Updated)"),
				resource.TestCheckResourceAttr(httpResourceName, "interval", "300"),
				resource.TestCheckResourceAttr(httpResourceName, "transaction_script", "  import { By, Key, until } from 'selenium-webdriver'; \n  import { driver, markers, credentials, downloads, transaction, test } from 'thousandeyes'; \n  runScript(); \n  async function runScript() \n  { const settings = test.getSettings();\n  // Load page\n  await driver.get(settings.url);\n  await driver.wait(until.titleIs('Digital Experience Monitoring | ThousandEyes'), 1000);\n  await driver.takeScreenshot();\n};\n"),
				resource.TestCheckResourceAttr(httpResourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(httpResourceName, "alert_rules.#", "2"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      tc.checkDestroyFunction,
				Steps: []resource.TestStep{
					{
						Config: testAccThousandEyesWebTransactionsConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesWebTransactionsConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckDefaultWebTransactionsResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_web_transaction",
			GetResource: func(id string) (interface{}, error) {
				return getWebTransactions(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesWebTransactionsConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getWebTransactions(id string) (interface{}, error) {
	api := (*tests.WebTransactionTestsAPIService)(&testClient.Common)
	req := api.GetWebTransactionsTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
