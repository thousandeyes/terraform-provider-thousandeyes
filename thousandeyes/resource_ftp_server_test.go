package thousandeyes

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesFTPServer_basic(t *testing.T) {
	testName := "tf-acc-test-ftp-server"
	resourceName := "thousandeyes_ftp_server.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckThousandEyesFTPServerDestroy,
		Steps:             testAccCheckThousandEyesFTPServerSteps(testName, resourceName),
	})
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

func testAccCheckThousandEyesFTPServerSteps(testName, resourceName string) []resource.TestStep {
	return []resource.TestStep{
		{
			Config: testAccThousandEyesFTPServerConfig(testName),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "password", "test_password"),
				resource.TestCheckResourceAttr(resourceName, "username", "test_username"),
				// Add more checks based on the resource attributes
			),
		},
	}
}

func testAccThousandEyesFTPServerConfig(testName string) string {
	return fmt.Sprintf(`
	data "thousandeyes_agent" "test"{
		agent_name = "Vancouver, Canada"
	}
		
	resource "thousandeyes_ftp_server" "test" {
		password       = "test_password"
		username       = "test_username"
		test_name      = "Acceptance Test - FTP"
  		description    = "description"
		request_type   = "Download"
  		ftp_time_limit = 10
		ftp_target_time = 1000
		interval       = 900
		alerts_enabled = false
        network_measurements = false
 		url 		   = "ftp://speedtest.tele2.net/"

	  	agents {
				agent_id = data.thousandeyes_agent.test.agent_id
	  	}
	}
	`)
}
