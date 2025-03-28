package thousandeyes

import (
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceFTPServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.FtpServerTestRequest{}, schemas.CommonSchema, nil),
		Create: resourceFTPServerCreate,
		Read:   resourceFTPServerRead,
		Update: resourceFTPServerUpdate,
		Delete: resourceFTPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create an FTP server test. This test type verifies the availability and performance of FTP servers. For more information, see [FTP Server Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#ftp-server-test).",
	}
	resource.Schema["password"] = schemas.CommonSchema["password-ftp"]
	resource.Schema["username"] = schemas.CommonSchema["username-ftp"]
	return &resource
}

func resourceFTPServerRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.FTPServerTestsAPIService)(&apiClient.Common)

		req := api.GetFtpServerTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceFTPServerUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.FTPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := ResourceUpdate(d, &tests.FtpServerTestRequest{})

	req := api.UpdateFtpServerTest(d.Id()).FtpServerTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceFTPServerRead(d, m)
}

func resourceFTPServerDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.FTPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteFtpServerTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceFTPServerCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.FTPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildFTPServerStruct(d)

	req := api.CreateFtpServerTest().FtpServerTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, _ := req.Execute()
	// if err != nil {
	// 	return err
	// }

	id := *resp.TestId
	d.SetId(id)
	return resourceFTPServerRead(d, m)
}

func buildFTPServerStruct(d *schema.ResourceData) *tests.FtpServerTestRequest {
	return ResourceBuildStruct(d, &tests.FtpServerTestRequest{})
}
