package thousandeyes

// import (
// 	"os"
// 	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/alerts"
)

func TestAccThousandEyesAlertRule(t *testing.T) {
	var alertRuleResourceName = "thousandeyes_alert_rule.test"
	testCases := []struct {
		name                 string
		createResourceFile   string
		updateResourceFile   string
		checkDestroyFunction func(*terraform.State) error
		checkCreateFunc      []resource.TestCheckFunc
		checkUpdateFunc      []resource.TestCheckFunc
	}{
		{
			name:               "test_association_maintained_after_alert_rule_update",
			createResourceFile: "acceptance_resources/alert_rule/create_test_with_alert_rule.tf",
			updateResourceFile: "acceptance_resources/alert_rule/update_alert_rule.tf",
			checkDestroyFunction: func(state *terraform.State) error {
				resourceList := []ResourceType{
					{
						Name:         "Agent To Server Test",
						ResourceName: "thousandeyes_agent_to_server",
						GetResource: func(id string) (interface{}, error) {
							return getAgentToServer(id)
						}},
					{
						Name:         "Agent To Server Alert Rule Test",
						ResourceName: "thousandeyes_alert_rule",
						GetResource: func(id string) (interface{}, error) {
							return getAlertRule(id)
						}},
				}
				return testAccCheckResourceDestroy(resourceList, state)
			},
			checkCreateFunc: []resource.TestCheckFunc{
				//alert rule is created
				//with 3 required violating rounds
				resource.TestCheckResourceAttr(alertRuleResourceName, "rounds_violating_required", "3"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				//alert rule is updated
				//to 4 required violating rounds
				resource.TestCheckResourceAttr(alertRuleResourceName, "rounds_violating_required", "4"),
				//and the test association is maintained
				resource.TestCheckResourceAttr(alertRuleResourceName, "test_ids.#", "1"),
			},
		},
	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			resource.Test(t, resource.TestCase{
// 				PreCheck:          func() { testAccPreCheck(t) },
// 				ProviderFactories: providerFactories,
// 				CheckDestroy:      tc.checkDestroyFunction,
// 				Steps: []resource.TestStep{
// 					{
// 						Config: testAccThousandEyesAlertRuleConfig(tc.createResourceFile),
// 						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
// 					},
// 					{
// 						Config: testAccThousandEyesAlertRuleConfig(tc.updateResourceFile),
// 						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
// 					},
// 				},
// 			})
// 		})
// 	}
// }

func testAccThousandEyesAlertRuleConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getAlertRule(id string) (interface{}, error) {
	api := (*alerts.AlertRulesAPIService)(&testClient.Common)
	req := api.GetAlertRule(id)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
