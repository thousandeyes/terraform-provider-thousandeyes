package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestAccThousandEyesAPI(t *testing.T) {
	var resourceName = "thousandeyes_api.test"
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
			name:               "create_update_delete_api_test",
			createResourceFile: "acceptance_resources/api/basic.tf",
			// updateResourceFile:   "acceptance_resources/api/update.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckDefaultAPIResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - API"),
				resource.TestCheckResourceAttr(resourceName, "interval", "120"),
				resource.TestCheckResourceAttr(resourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(resourceName, "alert_rules.#", "2"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - API (Updated)"),
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
						Config:             testAccThousandEyesAPIConfig(tc.createResourceFile),
						Check:              resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
						ExpectNonEmptyPlan: true,
					},
					// {
					// 	Config: testAccThousandEyesAPIConfig(tc.updateResourceFile),
					// 	Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					// },
				},
			})
		})
	}
}

func testAccCheckDefaultAPIResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_api",
			GetResource: func(id string) (interface{}, error) {
				return getAPI(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesAPIConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getAPI(id string) (interface{}, error) {
	api := (*tests.APITestsAPIService)(&testClient.Common)
	req := api.GetApiTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
