package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceAPI() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.ApiTestRequest{}, schemas.CommonSchema, nil),
		Create: resourceAPICreate,
		Read:   resourceAPIRead,
		Update: resourceAPIUpdate,
		Delete: resourceAPIDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "This resource provides visibility into the performance of critical web API endpoints within your application ecosystem. For more information, see [API Tests](https://docs.thousandeyes.com/product-documentation/api-test).",
	}
	return &resource
}

func resourceAPIRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.APITestsAPIService)(&apiClient.Common)

		req := api.GetApiTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceAPIUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.APITestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := buildAPITestStruct(d)

	req := api.UpdateApiTest(d.Id()).ApiTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceAPIRead(d, m)
}

func resourceAPIDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.APITestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteApiTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAPICreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.APITestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildAPITestStruct(d)

	req := api.CreateApiTest().ApiTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceAPIRead(d, m)
}

func buildAPITestStruct(d *schema.ResourceData) *tests.ApiTestRequest {
	return ResourceBuildStruct(d, &tests.ApiTestRequest{})
}
