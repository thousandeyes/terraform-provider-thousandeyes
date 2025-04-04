package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

const webTrEmulationDeviceIdKey emulationDeviceIdKeyType = "wt_emulation_device_id"

func resourceWebTransaction() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.WebTransactionTestRequest{}, schemas.CommonSchema, nil),
		Create: resourceWebTransactionCreate,
		Read:   resourceWebTransactionRead,
		Update: resourceWebTransactionUpdate,
		Delete: resourceWebTransactionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows users to create a transaction test. This test type is a scripted synthetic browser interaction that can traverse multiple pages and user actions. For more information, see [Transaction Tests](https://docs.thousandeyes.com/product-documentation/browser-synthetics/transaction-tests).",
	}
	return &resource
}

func resourceWebTransactionRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.WebTransactionTestsAPIService)(&apiClient.Common)

		req := api.GetWebTransactionsTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		edID := apiClient.GetConfig().Context.Value(webTrEmulationDeviceIdKey)
		if edID == nil {
			resp.EmulatedDeviceId = nil
		} else {
			apiClient.GetConfig().Context = GetContextWithAid(apiClient.GetConfig().Context)
		}
		return resp, err
	})
}

func resourceWebTransactionUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.WebTransactionTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := buildWebTransactionStruct(d)
	if update.EmulatedDeviceId != nil && len(*update.EmulatedDeviceId) > 0 {
		apiClient.GetConfig().Context = context.WithValue(
			apiClient.GetConfig().Context,
			webTrEmulationDeviceIdKey,
			struct{}{},
		)
	}

	req := api.UpdateWebTransactionsTest(d.Id()).WebTransactionTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceWebTransactionRead(d, m)
}

func resourceWebTransactionDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.WebTransactionTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteWebTransactionsTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceWebTransactionCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.WebTransactionTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildWebTransactionStruct(d)
	if local.EmulatedDeviceId != nil && len(*local.EmulatedDeviceId) > 0 {
		apiClient.GetConfig().Context = context.WithValue(
			apiClient.GetConfig().Context,
			webTrEmulationDeviceIdKey,
			struct{}{},
		)
	}

	req := api.CreateWebTransactionsTest().WebTransactionTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceWebTransactionRead(d, m)
}

func buildWebTransactionStruct(d *schema.ResourceData) *tests.WebTransactionTestRequest {
	return ResourceBuildStruct(d, &tests.WebTransactionTestRequest{})
}
