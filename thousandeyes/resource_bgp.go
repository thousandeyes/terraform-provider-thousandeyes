package thousandeyes

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceBGP() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.BgpTest{}, schemas, nil),
		Create: resourceBGPCreate,
		Read:   resourceBGPRead,
		Update: resourceBGPUpdate,
		Delete: resourceBGPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create a ThousandEyes BGP test. This test type collects BGP routing related information and presents a visualization of the route and relevants events on the timeline. For more information, see [BGP Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#bgp-test).",
	}
	return &resource
}

func resourceBGPRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.BGPTestsAPIService)(&apiClient.Common)

		req := api.GetBgpTest(id).Expand(tests.AllowedExpandBgpTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceBGPUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.BGPTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := ResourceUpdate(d, &tests.UpdateBgpTestRequest{})

	req := api.UpdateBgpTest(d.Id()).UpdateBgpTestRequest(*update).Expand(tests.AllowedExpandBgpTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}

	return resourceBGPRead(d, m)
}

func resourceBGPDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.BGPTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteBgpTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceBGPCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.BGPTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildBGPStruct(d)

	req := api.CreateBgpTest().BgpTestRequest(*local).Expand(tests.AllowedExpandBgpTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceBGPRead(d, m)
}

func buildBGPStruct(d *schema.ResourceData) *tests.BgpTestRequest {
	return ResourceBuildStruct(d, &tests.BgpTestRequest{})
}
