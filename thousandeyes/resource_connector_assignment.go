package thousandeyes

import (
	"context"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

const operationTypeWebhooks = "webhooks"

func resourceConnectorAssignment() *schema.Resource {
	return &schema.Resource{
		Schema:      schemas.ConnectorAssignmentSchema,
		Create:      resourceConnectorAssignmentCreate,
		Read:        resourceConnectorAssignmentRead,
		Update:      resourceConnectorAssignmentUpdate,
		Delete:      resourceConnectorAssignmentDelete,
		Description: "Manages all connector assignments for a webhook operation.\n\nThis resource is authoritative for the operation and uses PUT replace-all semantics:\n- To assign a connector, Terraform sends the full list including the connector ID.\n- To remove all connectors, Terraform sends an empty list (`[]`).\n\nAny connector assignments made outside Terraform for the same webhook operation will be removed on the next apply if they are not present in `connector_ids`.\n\nThe current API supports one connector per webhook operation.",
		Importer: &schema.ResourceImporter{
			StateContext: resourceConnectorAssignmentImport,
		},
	}
}

func resourceConnectorAssignmentImport(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := d.Set("webhook_operation_id", d.Id()); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, resourceConnectorAssignmentRead(d, m)
}

func resourceConnectorAssignmentCreate(d *schema.ResourceData, m interface{}) error {
	operationID := d.Get("webhook_operation_id").(string)
	log.Printf("[INFO] Creating ThousandEyes Connector Assignment for webhook operation %s", operationID)

	if _, err := setWebhookOperationConnectors(m.(*client.APIClient), operationID, expandConnectorIDs(d)); err != nil {
		return err
	}

	d.SetId(operationID)
	return resourceConnectorAssignmentRead(d, m)
}

func resourceConnectorAssignmentRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	operationID := d.Get("webhook_operation_id").(string)
	if operationID == "" {
		operationID = d.Id()
	}

	log.Printf("[INFO] Reading ThousandEyes Connector Assignment for webhook operation %s", operationID)

	assignments, err := getWebhookOperationConnectors(apiClient, operationID)
	if err != nil {
		if IsNotFoundError(err) {
			log.Printf("[INFO] Webhook operation %s not found, removing connector assignment from state", operationID)
			d.SetId("")
			return nil
		}
		return err
	}

	items := []string{}
	if assignments != nil && assignments.Items != nil {
		items = assignments.Items
	}

	if err := d.Set("webhook_operation_id", operationID); err != nil {
		return err
	}
	if err := d.Set("connector_ids", items); err != nil {
		return err
	}

	d.SetId(operationID)
	return nil
}

func resourceConnectorAssignmentUpdate(d *schema.ResourceData, m interface{}) error {
	operationID := d.Get("webhook_operation_id").(string)
	log.Printf("[INFO] Updating ThousandEyes Connector Assignment for webhook operation %s", operationID)

	if _, err := setWebhookOperationConnectors(m.(*client.APIClient), operationID, expandConnectorIDs(d)); err != nil {
		return err
	}

	return resourceConnectorAssignmentRead(d, m)
}

func resourceConnectorAssignmentDelete(d *schema.ResourceData, m interface{}) error {
	operationID := d.Get("webhook_operation_id").(string)
	if operationID == "" {
		operationID = d.Id()
	}

	log.Printf("[INFO] Deleting ThousandEyes Connector Assignment for webhook operation %s", operationID)

	if _, err := setWebhookOperationConnectors(m.(*client.APIClient), operationID, []string{}); err != nil && !IsNotFoundError(err) {
		return err
	}

	d.SetId("")
	return nil
}

func expandConnectorIDs(d *schema.ResourceData) []string {
	raw := d.Get("connector_ids").(*schema.Set).List()
	connectorIDs := make([]string, 0, len(raw))
	for _, item := range raw {
		if id, ok := item.(string); ok && id != "" {
			connectorIDs = append(connectorIDs, id)
		}
	}
	sort.Strings(connectorIDs)
	return connectorIDs
}

func getWebhookOperationConnectors(apiClient *client.APIClient, operationID string) (*connectors.Assignments, error) {
	api := (*connectors.OperationConnectorsAPIService)(&apiClient.Common)

	req := api.GetOperationConnectors(operationTypeWebhooks, operationID)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	return resp, err
}

func setWebhookOperationConnectors(apiClient *client.APIClient, operationID string, connectorIDs []string) (*connectors.Assignments, error) {
	api := (*connectors.OperationConnectorsAPIService)(&apiClient.Common)

	req := api.SetOperationConnectors(operationTypeWebhooks, operationID).RequestBody(connectorIDs)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	return resp, err
}

func aidStringFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	aid, ok := ctx.Value(accountGroupIdKey).(string)
	if !ok || aid == "" {
		return "", false
	}
	return aid, true
}
