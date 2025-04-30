package thousandeyes

// import (
// 	"os"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// 	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
// )

// func TestAccThousandEyesAccountGroup(t *testing.T) {
// 	var httpResourceName = "thousandeyes_account_group.test"
// 	var testCases = []struct {
// 		name                 string
// 		createResourceFile   string
// 		updateResourceFile   string
// 		resourceName         string
// 		checkDestroyFunction func(*terraform.State) error
// 		checkCreateFunc      []resource.TestCheckFunc
// 		checkUpdateFunc      []resource.TestCheckFunc
// 	}{
// 		{
// 			name:                 "create_update_delete_account_group_test",
// 			createResourceFile:   "acceptance_resources/account_group/basic.tf",
// 			updateResourceFile:   "acceptance_resources/account_group/update.tf",
// 			resourceName:         httpResourceName,
// 			checkDestroyFunction: testAccCheckAccountGroupResourceDestroy,
// 			checkCreateFunc: []resource.TestCheckFunc{
// 				resource.TestCheckResourceAttr(httpResourceName, "account_group_name", "User Acceptance Test - Account Group 1"),
// 				resource.TestCheckResourceAttr(httpResourceName, "agents.#", "2"),
// 			},
// 			checkUpdateFunc: []resource.TestCheckFunc{
// 				resource.TestCheckResourceAttr(httpResourceName, "account_group_name", "User Acceptance Test - Account Group (Updated)"),
// 				resource.TestCheckResourceAttr(httpResourceName, "agents.#", "1"),
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			resource.Test(t, resource.TestCase{
// 				PreCheck:          func() { testAccPreCheck(t) },
// 				ProviderFactories: providerFactories,
// 				CheckDestroy:      tc.checkDestroyFunction,
// 				Steps: []resource.TestStep{
// 					{
// 						Config: testAccThousandEyesAccountGroupConfig(tc.createResourceFile),
// 						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
// 					},
// 					{
// 						Config: testAccThousandEyesAccountGroupConfig(tc.updateResourceFile),
// 						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
// 					},
// 				},
// 			})
// 		})
// 	}
// }

// func testAccCheckAccountGroupResourceDestroy(s *terraform.State) error {
// 	resourceList := []ResourceType{
// 		{
// 			ResourceName: "thousandeyes_account_group",
// 			GetResource: func(id string) (interface{}, error) {
// 				return getAccountGroup(id)
// 			}},
// 	}
// 	return testAccCheckResourceDestroy(resourceList, s)
// }

// func testAccThousandEyesAccountGroupConfig(testResource string) string {
// 	content, err := os.ReadFile(testResource)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return string(content)
// }

// func getAccountGroup(id string) (interface{}, error) {
// 	api := (*administrative.AccountGroupsAPIService)(&testClient.Common)
// 	req := api.GetAccountGroup(id).Expand(administrative.AllowedExpandAccountGroupOptionsEnumValues)
// 	resp, _, err := req.Execute()
// 	return resp, err
// }
