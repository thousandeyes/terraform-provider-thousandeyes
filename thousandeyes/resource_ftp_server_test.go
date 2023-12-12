package thousandeyes

import (
	"fmt"
	"os"
	"strconv"
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
			checkDestroyFunction: testAccCheckThousandEyesFTPServerDestroy,
		},
		{
			name:                 "alerts_enabled",
			resourceFile:         "acceptance_resources/ftp_server/alerts_enabled.tf",
			resourceName:         "thousandeyes_ftp_server.test",
			checkDestroyFunction: testAccCheckThousandEyesFTPServerDestroy,
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

func testAccCheckThousandEyesFTPServerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "thousandeyes_ftp_server" {
			id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
			if err != nil {
				return err
			}
			_, err = testClient.GetFTPServer(id)
			if err == nil {
				return fmt.Errorf("FTPServer %s still exists", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckThousandEyesFTPServerSteps(testResource, resourceName string) []resource.TestStep {
	return []resource.TestStep{
		{
			Config: testAccThousandEyesFTPServerConfig(testResource),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "password", "test_password"),
				resource.TestCheckResourceAttr(resourceName, "username", "test_username"),
				// Add more checks based on the resource attributes
			),
		},
	}
}

func testAccThousandEyesFTPServerConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}
