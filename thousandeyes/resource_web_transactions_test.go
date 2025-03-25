package thousandeyes

// import (
// 	"os"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// )

// func TestAccThousandEyesWebTransactions(t *testing.T) {
// 	var httpResourceName = "thousandeyes_web_transaction.test"
// 	var testCases = []struct {
// 		name                 string
// 		resourceFile         string
// 		resourceName         string
// 		checkDestroyFunction func(*terraform.State) error
// 		checkFunc            []resource.TestCheckFunc
// 	}{
// 		{
// 			name:                 "basic",
// 			resourceFile:         "acceptance_resources/web_transactions/basic.tf",
// 			resourceName:         httpResourceName,
// 			checkDestroyFunction: testAccCheckDefaultWebTransactionsResourceDestroy,
// 			checkFunc: []resource.TestCheckFunc{
// 				resource.TestCheckResourceAttr(httpResourceName, "url", "https://www.thousandeyes.com"),
// 				resource.TestCheckResourceAttr(httpResourceName, "test_name", "User Acceptance Test - Web Transactions"),
// 				resource.TestCheckResourceAttr(httpResourceName, "interval", "120"),
// 				resource.TestCheckResourceAttr(httpResourceName, "transaction_script", "  import { By, Key, until } from 'selenium-webdriver'; \n  import { driver, markers, credentials, downloads, transaction, test } from 'thousandeyes'; \n  runScript(); \n  async function runScript() \n  { const settings = test.getSettings();\n  // Load page\n  await driver.get(settings.url);\n  await driver.wait(until.titleIs('Digital Experience Monitoring | ThousandEyes'), 1000);\n  await driver.takeScreenshot();\n};\n"),
// 				resource.TestCheckResourceAttr(httpResourceName, "alerts_enabled", "true"),
// 				resource.TestCheckResourceAttr(httpResourceName, "alert_rules.#", "2"),
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			resource.Test(t, resource.TestCase{
// 				PreCheck:          func() { testAccPreCheck(t) },
// 				ProviderFactories: providerFactories,
// 				CheckDestroy:      tc.checkDestroyFunction,
// 				Steps: []resource.TestStep{
// 					{
// 						Config: testAccThousandEyesWebTransactionsConfig(tc.resourceFile),
// 						Check:  resource.ComposeTestCheckFunc(tc.checkFunc...),
// 					},
// 				},
// 			})
// 		})
// 	}
// }

// func testAccCheckDefaultWebTransactionsResourceDestroy(s *terraform.State) error {
// 	resourceList := []ResourceType{
// 		{
// 			ResourceName: "thousandeyes_web_transaction",
// 			GetResource: func(id int64) (interface{}, error) {
// 				return testClient.GetWebTransaction(id)
// 			}},
// 	}
// 	return testAccCheckResourceDestroy(resourceList, s)
// }

// func testAccThousandEyesWebTransactionsConfig(testResource string) string {
// 	content, err := os.ReadFile(testResource)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return string(content)
// }