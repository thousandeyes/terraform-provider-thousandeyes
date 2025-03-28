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
		resourceFile         string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
		checkFunc            []resource.TestCheckFunc
	}{
		{
			name:                 "basic",
			resourceFile:         "acceptance_resources/sip_server/basic.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckSIPServerResourceDestroy,
			checkFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - SIP Server"),
				resource.TestCheckResourceAttr(resourceName, "target_sip_credentials.0.sip_registrar", "thousandeyes.com"),
				resource.TestCheckResourceAttr(resourceName, "target_sip_credentials.0.protocol", "tcp"),
				resource.TestCheckResourceAttr(resourceName, "target_sip_credentials.0.port", "5060"),
				resource.TestCheckResourceAttr(resourceName, "interval", "120"),
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
						Config: testAccThousandEyesSIPServerConfig(tc.resourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkFunc...),
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
