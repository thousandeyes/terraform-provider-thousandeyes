package thousandeyes

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesEndpointAgentToServerDynamic(t *testing.T) {
	const address = "thousandeyes_endpoint_agent_to_server_dynamic.test"
	name := "User Acceptance Test - Endpoint Agent To Server Dynamic-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	var originalID string
	testCases := []struct {
		name               string
		createResourceFile string
		updateResourceFile string
	}{{
		name:               "create_update_import_replace_delete_endpoint_agent_to_server_dynamic_test",
		createResourceFile: "acceptance_resources/endpoint_agent_to_server_dynamic/basic.tf",
		updateResourceFile: "acceptance_resources/endpoint_agent_to_server_dynamic/update.tf",
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			basic := testAccThousandEyesEndpointAgentToServerDynamicConfig(tc.createResourceFile, name)
			replacement := strings.Replace(basic, "max_machines        = 25", "max_machines        = 12", 1)
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testAccCheckEndpointAgentToServerDynamicResourceDestroy,
				Steps: []resource.TestStep{
					{Config: basic, Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExistsAndStoreID(address, &originalID),
						resource.TestCheckResourceAttr(address, "test_name", name),
						resource.TestCheckResourceAttr(address, "application", "webex"),
						resource.TestCheckResourceAttr(address, "agent_selector_type", "all-agents"),
						resource.TestCheckResourceAttr(address, "max_machines", "25"),
						resource.TestCheckResourceAttr(address, "interval", "60"),
						resource.TestCheckResourceAttr(address, "is_enabled", "false"),
					)},
					{Config: testAccThousandEyesEndpointAgentToServerDynamicConfig(tc.updateResourceFile, name), Check: resource.ComposeTestCheckFunc(
						testAccCheckDependentResourceHasCorrectID(address, "test_id", &originalID),
						resource.TestCheckResourceAttr(address, "test_name", name+" (Updated)"),
						resource.TestCheckResourceAttr(address, "application", "zoom"),
						resource.TestCheckResourceAttr(address, "interval", "120"),
						resource.TestCheckResourceAttr(address, "is_enabled", "true"),
						resource.TestCheckResourceAttr(address, "protocol", "tcp"),
					)},
					{ResourceName: address, ImportState: true, ImportStateVerify: true},
					{Config: basic, Check: resource.ComposeTestCheckFunc(
						testAccCheckDependentResourceHasCorrectID(address, "test_id", &originalID),
						resource.TestCheckResourceAttr(address, "interval", "60"),
						resource.TestCheckResourceAttr(address, "is_enabled", "false"),
						resource.TestCheckResourceAttr(address, "protocol", "icmp"),
					)},
					{Config: replacement, Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(address, "max_machines", "12"),
						func(s *terraform.State) error {
							if s.RootModule().Resources[address].Primary.ID == originalID {
								return fmt.Errorf("changing max_machines did not replace the test")
							}
							req := endpointAgentToServerDynamicAPI(testClient).GetAgentToServerEndpointDynamicTest(originalID)
							_, response, err := SetAidFromContext(testClient.GetConfig().Context, req).Execute()
							if response == nil || response.StatusCode != 404 {
								return fmt.Errorf("replaced test %s was not deleted: %v", originalID, err)
							}
							return nil
						},
					)},
				},
			})
		})
	}
}

func testAccCheckEndpointAgentToServerDynamicResourceDestroy(s *terraform.State) error {
	return testAccCheckResourceDestroy([]ResourceType{{
		ResourceName: "thousandeyes_endpoint_agent_to_server_dynamic",
		GetResource:  getEndpointAgentToServerDynamic,
	}}, s)
}

func getEndpointAgentToServerDynamic(id string) (interface{}, error) {
	req := endpointAgentToServerDynamicAPI(testClient).GetAgentToServerEndpointDynamicTest(id)
	resp, _, err := SetAidFromContext(testClient.GetConfig().Context, req).Execute()
	return resp, err
}

func testAccThousandEyesEndpointAgentToServerDynamicConfig(testResource, name string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(string(content), "User Acceptance Test - Endpoint Agent To Server Dynamic", name)
}
