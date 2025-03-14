package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesHTTPServer(t *testing.T) {
	var httpResourceName = "thousandeyes_http_server.test"
	var testCases = []struct {
		name                 string
		resourceFile         string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
		checkFunc            []resource.TestCheckFunc
	}{
		{
			name:                 "basic",
			resourceFile:         "acceptance_resources/http_server/basic.tf",
			resourceName:         httpResourceName,
			checkDestroyFunction: testAccCheckDefaultHTTPResourceDestroy,
			checkFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(httpResourceName, "test_name", "User Acceptance Test - HTTP"),
				resource.TestCheckResourceAttr(httpResourceName, "interval", "120"),
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
						Config: testAccThousandEyesHTTPServerConfig(tc.resourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckDefaultHTTPResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_http_server",
			GetResource: func(id int64) (interface{}, error) {
				return testClient.GetHTTPServer(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesHTTPServerConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}
