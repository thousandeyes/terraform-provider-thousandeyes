package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesFTPServer(t *testing.T) {
	testCases := []struct {
		name                 string
		resourceFile         string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
	}{
		{
			name:                 "basic",
			resourceFile:         "acceptance_resources/ftp_server/basic.tf",
			resourceName:         "thousandeyes_ftp_server.test",
			checkDestroyFunction: testAccCheckDefaultResourceDestroy,
		},
		{
			name:                 "alerts_enabled",
			resourceFile:         "acceptance_resources/ftp_server/alerts_enabled.tf",
			resourceName:         "thousandeyes_ftp_server.test",
			checkDestroyFunction: testAccCheckDefaultResourceDestroy,
		},
		{
			name:         "alerts_enabled_multiple_alert_rules",
			resourceFile: "acceptance_resources/ftp_server/alerts_enabled_multiple_alert_rules.tf",
			resourceName: "thousandeyes_ftp_server.test",
			checkDestroyFunction: func(state *terraform.State) error {
				resourceList := []ResourceType{
					{
						Name:         "FTP Server Test",
						ResourceName: "thousandeyes_ftp_server",
						GetResource: func(id int64) (interface{}, error) {
							return testClient.GetFTPServer(id)
						}},
					{
						Name:         "Alert Rule",
						ResourceName: "thousandeyes_alert_rule",
						GetResource: func(id int64) (interface{}, error) {
							return testClient.GetAlertRule(id)
						}},
				}
				return testAccCheckResourceDestroy(resourceList, state)
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
						Config: testAccThousandEyesFTPServerConfig(tc.resourceFile),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(tc.resourceName, "password", "test_password"),
							resource.TestCheckResourceAttr(tc.resourceName, "username", "test_username"),
							// Add more checks based on the resource attributes
						),
					},
				},
			})
		})
	}
}

func testAccCheckDefaultResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			Name:         "FTP Server Test",
			ResourceName: "thousandeyes_ftp_server",
			GetResource: func(id int64) (interface{}, error) {
				return testClient.GetFTPServer(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesFTPServerConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}
