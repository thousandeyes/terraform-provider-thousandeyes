package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestAccThousandEyesDNSServer(t *testing.T) {
	var resourceName = "thousandeyes_dns_server.test"
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
			name:                 "create_update_delete_dns_server_test",
			createResourceFile:   "acceptance_resources/dns_server/basic.tf",
			updateResourceFile:   "acceptance_resources/dns_server/update.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckDNSServerResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - DNS Server"),
				resource.TestCheckResourceAttr(resourceName, "domain", "thousandeyes.com A"),
				resource.TestCheckResourceAttr(resourceName, "interval", "120"),
				resource.TestCheckResourceAttr(resourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(resourceName, "alert_rules.#", "2"),
				resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "ns-1458.awsdns-54.org"),
				resource.TestCheckResourceAttr(resourceName, "dns_servers.2", "ns-cloud-d1.googledomains.com"),
				resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "ns-597.awsdns-10.net"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - DNS Server (Updated)"),
				resource.TestCheckResourceAttr(resourceName, "domain", "thousandeyes.com A"),
				resource.TestCheckResourceAttr(resourceName, "interval", "300"),
				resource.TestCheckResourceAttr(resourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(resourceName, "alert_rules.#", "2"),
				resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "ns-1458.awsdns-54.org"),
				resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "ns-cloud-d1.googledomains.com"),
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
						Config: testAccThousandEyesDNSServerConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesDNSServerConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckDNSServerResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_dns_server",
			GetResource: func(id string) (interface{}, error) {
				return getDNSServer(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesDNSServerConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getDNSServer(id string) (interface{}, error) {
	api := (*tests.DNSServerTestsAPIService)(&testClient.Common)
	req := api.GetDnsServerTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
