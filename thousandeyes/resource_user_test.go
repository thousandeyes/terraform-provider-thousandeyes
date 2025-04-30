package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
)

func TestAccThousandEyesUser(t *testing.T) {
	var httpResourceName = "thousandeyes_user.test"
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
			name:                 "create_update_delete_user_test",
			createResourceFile:   "acceptance_resources/user/basic.tf",
			updateResourceFile:   "acceptance_resources/user/update.tf",
			resourceName:         httpResourceName,
			checkDestroyFunction: testAccCheckUserResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "name", "User Acceptance Test - User"),
				resource.TestCheckResourceAttr(httpResourceName, "email", "example@test.com"),
				resource.TestCheckResourceAttr(httpResourceName, "all_account_group_role_ids.#", "1"),
				resource.TestCheckResourceAttr(httpResourceName, "account_group_roles.#", "1"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "name", "User Acceptance Test - User (Updated)"),
				resource.TestCheckResourceAttr(httpResourceName, "email", "example@test.com"),
				resource.TestCheckResourceAttr(httpResourceName, "all_account_group_role_ids.#", "1"),
				resource.TestCheckResourceAttr(httpResourceName, "account_group_roles.#", "1"),
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
						Config: testAccThousandEyesUserConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesUserConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckUserResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_user",
			GetResource: func(id string) (interface{}, error) {
				return getUser(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesUserConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getUser(id string) (interface{}, error) {
	api := (*administrative.UsersAPIService)(&testClient.Common)
	req := api.GetUser(id)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
