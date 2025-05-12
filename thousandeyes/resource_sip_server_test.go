package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestAccThousandEyesSIPServer(t *testing.T) {
	var resourceName = "thousandeyes_sip_server.test"
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
			name:                 "create_update_delete_sip_server_test",
			createResourceFile:   "acceptance_resources/sip_server/basic.tf",
			updateResourceFile:   "acceptance_resources/sip_server/update.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckSIPServerResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - SIP Server"),
				resource.TestCheckResourceAttr(resourceName, "target_sip_credentials.0.sip_registrar", "thousandeyes.com"),
				resource.TestCheckResourceAttr(resourceName, "target_sip_credentials.0.protocol", "tcp"),
				resource.TestCheckResourceAttr(resourceName, "target_sip_credentials.0.port", "5060"),
				resource.TestCheckResourceAttr(resourceName, "interval", "120"),
				resource.TestCheckResourceAttr(resourceName, "probe_mode", "sack"),
				resource.TestCheckResourceAttr(resourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(resourceName, "alert_rules.#", "2"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - SIP Server (Updated)"),
				resource.TestCheckResourceAttr(resourceName, "target_sip_credentials.0.sip_registrar", "thousandeyes.com"),
				resource.TestCheckResourceAttr(resourceName, "target_sip_credentials.0.protocol", "tcp"),
				resource.TestCheckResourceAttr(resourceName, "target_sip_credentials.0.port", "5065"),
				resource.TestCheckResourceAttr(resourceName, "interval", "300"),
				resource.TestCheckResourceAttr(resourceName, "probe_mode", "sack"),
				resource.TestCheckResourceAttr(resourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(resourceName, "alert_rules.#", "2"),
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
						Config: testAccThousandEyesSIPServerConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesSIPServerConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckSIPServerResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_sip_server",
			GetResource: func(id string) (interface{}, error) {
				return getSIPServer(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesSIPServerConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getSIPServer(id string) (interface{}, error) {
	api := (*tests.SIPServerTestsAPIService)(&testClient.Common)
	req := api.GetSipServerTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
