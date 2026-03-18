package thousandeyes

import (
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
	apiClient := m.(*client.APIClient)
	api := (*connectors.WebhookOperationsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Reading ThousandEyes Webhook Operation %s", d.Id())

	req := api.GetWebhookOperation(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
	resp, httpResp, err := req.Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			log.Printf("[INFO] Webhook Operation %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	return setWebhookOperationResourceData(d, resp)
}

func resourceWebhookOperationUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*connectors.WebhookOperationsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Webhook Operation %s", d.Id())
	update := buildWebhookOperationStruct(d)

	req := api.UpdateWebhookOperation(d.Id()).WebhookOperation(*update)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

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
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

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
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

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

func setWebhookOperationResourceData(d *schema.ResourceData, webhook *connectors.WebhookOperation) error {
	if err := d.Set("name", webhook.Name); err != nil {
		return err
	}
	if err := d.Set("category", string(webhook.Category)); err != nil {
		return err
	}
	if err := d.Set("status", string(webhook.Status)); err != nil {
		return err
	}
	if err := d.Set("enabled", webhook.GetEnabled()); err != nil {
		return err
	}
	if webhook.Path != nil {
		if err := d.Set("path", *webhook.Path); err != nil {
			return err
		}
	} else {
		if err := d.Set("path", nil); err != nil {
			return err
		}
	}
	if webhook.Payload != nil {
		if err := d.Set("payload", *webhook.Payload); err != nil {
			return err
		}
	} else {
		if err := d.Set("payload", nil); err != nil {
			return err
		}
	}
	if webhook.QueryParams != nil {
		if err := d.Set("query_params", *webhook.QueryParams); err != nil {
			return err
		}
	} else {
		if err := d.Set("query_params", nil); err != nil {
			return err
		}
	}
	if webhook.Type != nil {
		if err := d.Set("type", string(*webhook.Type)); err != nil {
			return err
		}
	} else {
		if err := d.Set("type", nil); err != nil {
			return err
		}
	}
	if webhook.Links != nil && webhook.Links.Self != nil {
		if err := d.Set("link", webhook.Links.Self.Href); err != nil {
			return err
		}
	} else {
		if err := d.Set("link", nil); err != nil {
			return err
		}
	}

	if len(webhook.Headers) > 0 {
		if err := d.Set("headers", flattenWebhookOperationHeaders(webhook.Headers)); err != nil {
			return err
		}
	} else {
		if err := d.Set("headers", nil); err != nil {
			return err
		}
	}

	return nil
}

func flattenWebhookOperationHeaders(headers []connectors.Header) []interface{} {
	out := make([]interface{}, 0, len(headers))
	for _, header := range headers {
		out = append(out, map[string]interface{}{
			"name":  header.Name,
			"value": header.Value,
		})
	}
	return out
}
