package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestAccThousandEyesFTPServer(t *testing.T) {
	var ftpResourceName = "thousandeyes_ftp_server.test"
	testCases := []struct {
		name                 string
		createResourceFile   string
		updateResourceFile   string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
		checkCreateFunc      []resource.TestCheckFunc
		checkUpdateFunc      []resource.TestCheckFunc
	}{
		{
			name:                 "basic",
			createResourceFile:   "acceptance_resources/ftp_server/basic.tf",
			updateResourceFile:   "acceptance_resources/ftp_server/update.tf",
			resourceName:         ftpResourceName,
			checkDestroyFunction: testAccCheckDefaultResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(ftpResourceName, "url", "ftp://speedtest.tele2.net/"),
				resource.TestCheckResourceAttr(ftpResourceName, "test_name", "User Acceptance Test - FTP Server"),
				resource.TestCheckResourceAttr(ftpResourceName, "interval", "120"),
				resource.TestCheckResourceAttr(ftpResourceName, "password", ""),
				resource.TestCheckResourceAttr(ftpResourceName, "username", "test_username"),
				resource.TestCheckResourceAttr(ftpResourceName, "description", "description"),
				resource.TestCheckResourceAttr(ftpResourceName, "request_type", "download"),
				resource.TestCheckResourceAttr(ftpResourceName, "ftp_time_limit", "10"),
				resource.TestCheckResourceAttr(ftpResourceName, "ftp_target_time", "1000"),
				resource.TestCheckResourceAttr(ftpResourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(ftpResourceName, "alert_rules.#", "2"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(ftpResourceName, "url", "ftp://speedtest.tele2.net/"),
				resource.TestCheckResourceAttr(ftpResourceName, "test_name", "User Acceptance Test - FTP Server (Updated)"),
				resource.TestCheckResourceAttr(ftpResourceName, "interval", "300"),
				resource.TestCheckResourceAttr(ftpResourceName, "password", ""),
				resource.TestCheckResourceAttr(ftpResourceName, "username", "test_username"),
				resource.TestCheckResourceAttr(ftpResourceName, "description", "description"),
				resource.TestCheckResourceAttr(ftpResourceName, "request_type", "download"),
				resource.TestCheckResourceAttr(ftpResourceName, "ftp_time_limit", "10"),
				resource.TestCheckResourceAttr(ftpResourceName, "ftp_target_time", "1000"),
				resource.TestCheckResourceAttr(ftpResourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(ftpResourceName, "alert_rules.#", "2"),
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
						Config:             testAccThousandEyesFTPServerConfig(tc.createResourceFile),
						Check:              resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
						ExpectNonEmptyPlan: true,
					},
					{
						Config:             testAccThousandEyesFTPServerConfig(tc.updateResourceFile),
						Check:              resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
						ExpectNonEmptyPlan: true,
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
			GetResource: func(id string) (interface{}, error) {
				return getFTPServer(id)
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

func getFTPServer(id string) (interface{}, error) {
	api := (*tests.FTPServerTestsAPIService)(&testClient.Common)
	req := api.GetFtpServerTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
