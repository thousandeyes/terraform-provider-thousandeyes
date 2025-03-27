package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestAccThousandEyesBGP(t *testing.T) {
	var resourceName = "thousandeyes_bgp.test"
	var testCases = []struct {
		name                 string
		resourceFile         string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
		checkFunc            []resource.TestCheckFunc
	}{
		{
			name:                 "basic",
			resourceFile:         "acceptance_resources/bgp/basic.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckBGPResourceDestroy,
			checkFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - BGP"),
				resource.TestCheckResourceAttr(resourceName, "use_public_bgp", "true"),
				resource.TestCheckResourceAttr(resourceName, "prefix", "192.0.2.0/24"),
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
						Config: testAccThousandEyesBGPConfig(tc.resourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckBGPResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_bgp",
			GetResource: func(id string) (interface{}, error) {
				return getBGP(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesBGPConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getBGP(id string) (interface{}, error) {
	api := (*tests.BGPTestsAPIService)(&testClient.Common)
	req := api.GetBgpTest(id)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
