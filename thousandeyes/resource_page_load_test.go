package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesPageLoadServer(t *testing.T) {
	var httpResourceName = "thousandeyes_page_load.test"
	var testCases = []struct {
		name                 string
		resourceFile         string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
		checkFunc            []resource.TestCheckFunc
	}{
		{
			name:                 "basic",
			resourceFile:         "acceptance_resources/page_load/basic.tf",
			resourceName:         httpResourceName,
			checkDestroyFunction: testAccCheckDefaultPageLoadResourceDestroy,
			checkFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(httpResourceName, "test_name", "User Acceptance Test - Page Load"),
				resource.TestCheckResourceAttr(httpResourceName, "interval", "120"),
				resource.TestCheckResourceAttr(httpResourceName, "http_interval", "120"),
				resource.TestCheckResourceAttr(httpResourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(httpResourceName, "alert_rules.#", "2"),
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
						Config: testAccThousandEyesPageLoadConfig(tc.resourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckDefaultPageLoadResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_page_load",
			GetResource: func(id int64) (interface{}, error) {
				return testClient.GetPageLoad(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesPageLoadConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}
