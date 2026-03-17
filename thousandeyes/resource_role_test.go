package thousandeyes

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
)

func TestAccThousandEyesRole(t *testing.T) {
	var httpResourceName = "thousandeyes_role.test"
	roleName, roleNameUpdated := testAccRoleNames()
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
			name:                 "create_update_delete_role_test",
			createResourceFile:   "acceptance_resources/role/basic.tf",
			updateResourceFile:   "acceptance_resources/role/update.tf",
			resourceName:         httpResourceName,
			checkDestroyFunction: testAccCheckRoleResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "name", roleName),
				resource.TestCheckResourceAttr(httpResourceName, "permissions.#", "1"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "name", roleNameUpdated),
				resource.TestCheckResourceAttr(httpResourceName, "permissions.#", "2"),
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
						Config: testAccThousandEyesRoleConfig(tc.createResourceFile, roleName, roleNameUpdated),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesRoleConfig(tc.updateResourceFile, roleName, roleNameUpdated),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccRoleNames() (string, string) {
	suffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("User Acceptance Test - Role-%s", suffix)
	return name, name + "-Updated"
}

func testAccCheckRoleResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_role",
			GetResource: func(id string) (interface{}, error) {
				return getRole(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesRoleConfig(testResource, roleName, roleNameUpdated string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}

	config := string(content)
	updatedRolePlaceholder := "__THOUSANDEYES_ROLE_UPDATED_NAME__"
	config = strings.ReplaceAll(config, "User Acceptance Test - Role (Updated)", updatedRolePlaceholder)
	config = strings.ReplaceAll(config, "User Acceptance Test - Role", roleName)
	config = strings.ReplaceAll(config, updatedRolePlaceholder, roleNameUpdated)
	return config
}

func getRole(id string) (interface{}, error) {
	api := (*administrative.RolesAPIService)(&testClient.Common)
	req := api.GetRole(id)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
