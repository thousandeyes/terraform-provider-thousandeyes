package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestAccThousandEyesPageLoadServer(t *testing.T) {
	var httpResourceName = "thousandeyes_page_load.test"
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
			name:                 "create_update_delete_page_load_test",
			createResourceFile:   "acceptance_resources/page_load/basic.tf",
			updateResourceFile:   "acceptance_resources/page_load/update.tf",
			resourceName:         httpResourceName,
			checkDestroyFunction: testAccCheckDefaultPageLoadResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(httpResourceName, "test_name", "New User Acceptance Test - Page Load"),
				resource.TestCheckResourceAttr(httpResourceName, "interval", "120"),
				resource.TestCheckResourceAttr(httpResourceName, "http_interval", "120"),
				resource.TestCheckResourceAttr(httpResourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(httpResourceName, "alert_rules.#", "2"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(httpResourceName, "test_name", "New User Acceptance Test - Page Load (Updated)"),
				resource.TestCheckResourceAttr(httpResourceName, "interval", "300"),
				resource.TestCheckResourceAttr(httpResourceName, "http_interval", "300"),
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
						Config: testAccThousandEyesPageLoadConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesPageLoadConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
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
			GetResource: func(id string) (interface{}, error) {
				return getPageLoad(id)
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

func getPageLoad(id string) (interface{}, error) {
	api := (*tests.PageLoadTestsAPIService)(&testClient.Common)
	req := api.GetPageLoadTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
