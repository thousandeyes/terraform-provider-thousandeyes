package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceHTTPServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.HttpServerTestRequest{}, schemas.CommonSchema, nil),
		Create: resourceHTTPServerCreate,
		Read:   resourceHTTPServerRead,
		Update: resourceHTTPServerUpdate,
		Delete: resourceHTTPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create an HTTP server test. This test type measures the availability and performance of an HTTP service. For more information, see [HTTP Server Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#http-server-test).",
	}
	resource.Schema["oauth"] = schemas.CommonSchema["oauth"]
	return &resource
}

func resourceHTTPServerRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)

		req := api.GetHttpServerTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceHTTPServerUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := buildHTTPServerStruct(d)

	req := api.UpdateHttpServerTest(d.Id()).HttpServerTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceHTTPServerRead(d, m)
}

func resourceHTTPServerDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteHttpServerTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceHTTPServerCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildHTTPServerStruct(d)

	req := api.CreateHttpServerTest().HttpServerTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceHTTPServerRead(d, m)
}

func buildHTTPServerStruct(d *schema.ResourceData) *tests.HttpServerTestRequest {
	return ResourceBuildStruct(d, &tests.HttpServerTestRequest{})
}
