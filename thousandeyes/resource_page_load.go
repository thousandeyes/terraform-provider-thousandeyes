package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

type emulationDeviceIdKeyType string

const emulationDeviceIdKey emulationDeviceIdKeyType = "emulation_device_id"

func resourcePageLoad() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.PageLoadTestRequest{}, schemas.CommonSchema, nil),
		Create: resourcePageLoadCreate,
		Read:   resourcePageLoadRead,
		Update: resourcePageLoadUpdate,
		Delete: resourcePageLoadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create a page load test. This test type obtains in-browser site performance metrics. For more information, see [Page Load Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#page-load-test).",
	}
	return &resource
}

func resourcePageLoadRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.PageLoadTestsAPIService)(&apiClient.Common)

		req := api.GetPageLoadTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		edID, ok := apiClient.GetConfig().Context.Value(emulationDeviceIdKey).(*string)
		if !ok || edID == nil {
			resp.EmulatedDeviceId = nil
		}
		apiClient.GetConfig().Context = context.WithValue(apiClient.GetConfig().Context, emulationDeviceIdKey, nil)
		return resp, err
	})
}

func resourcePageLoadUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.PageLoadTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := buildPageLoadStruct(d)
	if update.EmulatedDeviceId != nil && len(*update.EmulatedDeviceId) > 0 {
		apiClient.GetConfig().Context = context.WithValue(
			apiClient.GetConfig().Context,
			emulationDeviceIdKey,
			update.EmulatedDeviceId,
		)
	}

	req := api.UpdatePageLoadTest(d.Id()).PageLoadTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourcePageLoadRead(d, m)
}

func resourcePageLoadDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.PageLoadTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeletePageLoadTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourcePageLoadCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.PageLoadTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildPageLoadStruct(d)
	if local.EmulatedDeviceId != nil && len(*local.EmulatedDeviceId) > 0 {
		apiClient.GetConfig().Context = context.WithValue(
			apiClient.GetConfig().Context,
			emulationDeviceIdKey,
			local.EmulatedDeviceId,
		)
	}

	req := api.CreatePageLoadTest().PageLoadTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourcePageLoadRead(d, m)
}

func buildPageLoadStruct(d *schema.ResourceData) *tests.PageLoadTestRequest {
	return ResourceBuildStruct(d, &tests.PageLoadTestRequest{})
}
