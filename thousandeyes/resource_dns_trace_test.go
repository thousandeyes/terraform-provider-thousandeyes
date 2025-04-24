package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestAccThousandEyesDNSTrace(t *testing.T) {
	var resourceName = "thousandeyes_dns_trace.test"
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
			name:                 "create_update_delete_dns_trace_test",
			createResourceFile:   "acceptance_resources/dns_trace/basic.tf",
			updateResourceFile:   "acceptance_resources/dns_trace/update.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckDNSTraceResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - DNS Trace"),
				resource.TestCheckResourceAttr(resourceName, "domain", "thousandeyes.com A"),
				resource.TestCheckResourceAttr(resourceName, "interval", "120"),
				resource.TestCheckResourceAttr(resourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(resourceName, "alert_rules.#", "2"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - DNS Trace (Updated)"),
				resource.TestCheckResourceAttr(resourceName, "domain", "thousandeyes.com ANY"),
				resource.TestCheckResourceAttr(resourceName, "interval", "300"),
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
						Config: testAccThousandEyesDNSTraceConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesDNSTraceConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
						//ExpectNonEmptyPlan: true,
					},
				},
			})
		})
	}
}

func testAccCheckDNSTraceResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_dns_trace",
			GetResource: func(id string) (interface{}, error) {
				return getDNSTrace(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesDNSTraceConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getDNSTrace(id string) (interface{}, error) {
	api := (*tests.DNSTraceTestsAPIService)(&testClient.Common)
	req := api.GetDnsTraceTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
