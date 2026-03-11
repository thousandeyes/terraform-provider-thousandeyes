package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestAccThousandEyesDashboard(t *testing.T) {
	var resourceName = "thousandeyes_dashboard.test_dashboard"
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
			name:                 "create_update_delete_dashboard_test",
			createResourceFile:   "acceptance_resources/dashboard/basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/update.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "title", "Test Dashboard"),
				resource.TestCheckResourceAttr(resourceName, "description", "Test Dashboard Description"),
				resource.TestCheckResourceAttr(resourceName, "global_filter_id", "123"),
				resource.TestCheckResourceAttr(resourceName, "is_private", "false"),
				resource.TestCheckResourceAttr(resourceName, "default_timespan.0.duration", "100"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "title", "Test Dashboard (Updated)"),
				resource.TestCheckResourceAttr(resourceName, "description", "Updated Test Dashboard Description"),
				resource.TestCheckResourceAttr(resourceName, "is_private", "true"),
				resource.TestCheckResourceAttr(resourceName, "default_timespan.0.duration", "3600"),
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
						Config: testAccThousandEyesDashboardConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesDashboardConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckDashboardResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_dashboard",
			GetResource: func(id string) (interface{}, error) {
				return getDashboard(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesDashboardConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getDashboard(id string) (interface{}, error) {
	api := (*dashboards.DashboardsAPIService)(&testClient.Common)
	req := api.GetDashboard(id)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
