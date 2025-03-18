package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesDNSSEC(t *testing.T) {
	var resourceName = "thousandeyes_dnssec.test"
	var testCases = []struct {
		name                 string
		resourceFile         string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
		checkFunc            []resource.TestCheckFunc
	}{
		{
			name:                 "basic",
			resourceFile:         "acceptance_resources/dnssec/basic.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckDNSSECResourceDestroy,
			checkFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - DNSSEC"),
				resource.TestCheckResourceAttr(resourceName, "domain", "thousandeyes.com A"),
				resource.TestCheckResourceAttr(resourceName, "interval", "120"),
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
						Config: testAccThousandEyesDNSSECConfig(tc.resourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckDNSSECResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_dnssec",
			GetResource: func(id int64) (interface{}, error) {
				return testClient.GetDNSSec(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesDNSSECConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}