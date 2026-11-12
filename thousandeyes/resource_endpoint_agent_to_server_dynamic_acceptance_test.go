package thousandeyes

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesEndpointAgentToServerDynamic(t *testing.T) {
	const address = "thousandeyes_endpoint_agent_to_server_dynamic.test"
	testCases := []struct {
		name               string
		selectorType       string
		createResourceFile string
		updateResourceFile string
	}{
		{name: "all-agents", selectorType: "all-agents", createResourceFile: "acceptance_resources/endpoint_agent_to_server_dynamic/basic.tf", updateResourceFile: "acceptance_resources/endpoint_agent_to_server_dynamic/update.tf"},
		{name: "specific-agents", selectorType: "specific-agents", createResourceFile: "acceptance_resources/endpoint_agent_to_server_dynamic/specific_agents_basic.tf", updateResourceFile: "acceptance_resources/endpoint_agent_to_server_dynamic/specific_agents_update.tf"},
		{name: "agent-tags", selectorType: "agent-tags", createResourceFile: "acceptance_resources/endpoint_agent_to_server_dynamic/agent_tags_basic.tf", updateResourceFile: "acceptance_resources/endpoint_agent_to_server_dynamic/agent_tags_update.tf"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name := "User Acceptance Test - Endpoint Dynamic-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
			var originalID string
			selectorChecks := []resource.TestCheckFunc{resource.TestCheckResourceAttr(address, "agent_selector_type", tc.selectorType)}
			switch tc.selectorType {
			case "all-agents":
				selectorChecks = append(selectorChecks, resource.TestCheckResourceAttr(address, "agents.#", "0"), resource.TestCheckResourceAttr(address, "tag_ids.#", "0"))
			case "specific-agents":
				selectorChecks = append(selectorChecks, resource.TestCheckResourceAttr(address, "agents.#", "1"), resource.TestCheckResourceAttr(address, "tag_ids.#", "0"), resource.TestCheckResourceAttrPair(address, "agents.0", "data.thousandeyes_endpoint_agent.selected", "agent_id"))
			case "agent-tags":
				selectorChecks = append(selectorChecks, resource.TestCheckResourceAttr(address, "agents.#", "0"), resource.TestCheckResourceAttr(address, "tag_ids.#", "1"), resource.TestCheckResourceAttrPair(address, "tag_ids.0", "thousandeyes_tag.endpoint", "id"))
			}
			checkSelector := resource.ComposeTestCheckFunc(selectorChecks...)
			basic := testAccThousandEyesEndpointAgentToServerDynamicConfig(tc.createResourceFile, name)
			replacement := strings.Replace(basic, "max_machines        = 25", "max_machines        = 12", 1)
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					if tc.selectorType == "specific-agents" && os.Getenv("TE_ENDPOINT_AGENT_NAME") == "" {
						t.Fatal("TE_ENDPOINT_AGENT_NAME must identify an existing Endpoint Agent in the test account")
					}
					testAccPreCheck(t)
				},
				ProviderFactories: providerFactories,
				CheckDestroy:      testAccCheckEndpointAgentToServerDynamicResourceDestroy,
				Steps: []resource.TestStep{
					{Config: basic, Check: resource.ComposeTestCheckFunc(
						checkSelector,
						testAccCheckResourceExistsAndStoreID(address, &originalID),
						resource.TestCheckResourceAttr(address, "test_name", name),
						resource.TestCheckResourceAttr(address, "application", "webex"),
						resource.TestCheckResourceAttr(address, "max_machines", "25"),
						resource.TestCheckResourceAttr(address, "interval", "60"),
						resource.TestCheckResourceAttr(address, "is_enabled", "false"),
					)},
					{Config: testAccThousandEyesEndpointAgentToServerDynamicConfig(tc.updateResourceFile, name), Check: resource.ComposeTestCheckFunc(
						checkSelector,
						testAccCheckDependentResourceHasCorrectID(address, "test_id", &originalID),
						resource.TestCheckResourceAttr(address, "test_name", name+" (Updated)"),
						resource.TestCheckResourceAttr(address, "application", "zoom"),
						resource.TestCheckResourceAttr(address, "interval", "120"),
						resource.TestCheckResourceAttr(address, "is_enabled", "true"),
						resource.TestCheckResourceAttr(address, "protocol", "tcp"),
					)},
					{ResourceName: address, ImportState: true, ImportStateVerify: true},
					{Config: basic, Check: resource.ComposeTestCheckFunc(
						checkSelector,
						testAccCheckDependentResourceHasCorrectID(address, "test_id", &originalID),
						resource.TestCheckResourceAttr(address, "interval", "60"),
						resource.TestCheckResourceAttr(address, "is_enabled", "false"),
						resource.TestCheckResourceAttr(address, "protocol", "icmp"),
					)},
					{Config: replacement, Check: resource.ComposeTestCheckFunc(
						checkSelector,
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
	}, {ResourceName: "thousandeyes_tag", GetResource: func(id string) (interface{}, error) { return getTag(id) }}}, s)
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
	config := strings.ReplaceAll(string(content), "User Acceptance Test - Endpoint Agent To Server Dynamic", name)
	return strings.ReplaceAll(config, "\"__ENDPOINT_AGENT_NAME__\"", strconv.Quote(os.Getenv("TE_ENDPOINT_AGENT_NAME")))
}
