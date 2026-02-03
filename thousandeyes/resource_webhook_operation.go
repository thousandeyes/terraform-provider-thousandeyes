package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

func resourceWebhookOperation() *schema.Resource {
	resource := schema.Resource{
		Schema: schemas.WebhookOperationSchema,
		Create: resourceWebhookOperationCreate,
		Read:   resourceWebhookOperationRead,
		Update: resourceWebhookOperationUpdate,
		Delete: resourceWebhookOperationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "This resource allows you to create and manage webhook operations for ThousandEyes integrations. Webhook operations define what is sent (payload template, headers, path) when a webhook is triggered.",
	}
	return &resource
}

func resourceWebhookOperationRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*connectors.WebhookOperationsAPIService)(&apiClient.Common)

		req := api.GetWebhookOperation(id)
		req = SetAidFloatFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceWebhookOperationUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*connectors.WebhookOperationsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Webhook Operation %s", d.Id())
	update := buildWebhookOperationStruct(d)

	req := api.UpdateWebhookOperation(d.Id()).WebhookOperation(*update)
	req = SetAidFloatFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceWebhookOperationRead(d, m)
}

func resourceWebhookOperationDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*connectors.WebhookOperationsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Webhook Operation %s", d.Id())

	req := api.DeleteWebhookOperation(d.Id())
	req = SetAidFloatFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceWebhookOperationCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*connectors.WebhookOperationsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Webhook Operation")
	local := buildWebhookOperationStruct(d)

	req := api.CreateWebhookOperation().WebhookOperation(*local)
	req = SetAidFloatFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.Id
	d.SetId(id)
	return resourceWebhookOperationRead(d, m)
}

func buildWebhookOperationStruct(d *schema.ResourceData) *connectors.WebhookOperation {
	return ResourceBuildStruct(d, &connectors.WebhookOperation{})
}
