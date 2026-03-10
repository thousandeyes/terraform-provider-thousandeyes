package thousandeyes

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

func resourceConnectorAssignment() *schema.Resource {
	return &schema.Resource{
		Schema:      schemas.ConnectorAssignmentSchema,
		Create:      resourceConnectorAssignmentCreate,
		Read:        resourceConnectorAssignmentRead,
		Update:      resourceConnectorAssignmentUpdate,
		Delete:      resourceConnectorAssignmentDelete,
		Description: "Manages all webhook operation assignments for a connector.\n\nThis resource is authoritative for the connector and uses PUT replace-all semantics:\n- To add an operation, Terraform sends the full list including the new operation ID.\n- To remove an operation, Terraform sends the full list without that ID.\n- To remove all operation assignments, Terraform sends an empty list (`[]`).\n\nAny operation assignments made outside Terraform for the same connector will be removed on the next apply if they are not present in `operation_ids`.",
		Importer: &schema.ResourceImporter{
			StateContext: resourceConnectorAssignmentImport,
		},
	}
}

func resourceConnectorAssignmentImport(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := d.Set("connector_id", d.Id()); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, resourceConnectorAssignmentRead(d, m)
}

func resourceConnectorAssignmentCreate(d *schema.ResourceData, m interface{}) error {
	connectorID := d.Get("connector_id").(string)
	log.Printf("[INFO] Creating ThousandEyes Connector Assignment for connector %s", connectorID)

	if _, err := setConnectorOperations(m.(*client.APIClient), connectorID, expandOperationIDs(d)); err != nil {
		return err
	}

	d.SetId(connectorID)
	return resourceConnectorAssignmentRead(d, m)
}

func resourceConnectorAssignmentRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	connectorID := d.Get("connector_id").(string)
	if connectorID == "" {
		connectorID = d.Id()
	}

	log.Printf("[INFO] Reading ThousandEyes Connector Assignment for connector %s", connectorID)

	assignments, err := getConnectorOperations(apiClient, connectorID)
	if err != nil {
		if IsNotFoundError(err) {
			log.Printf("[INFO] Connector %s not found, removing connector assignment from state", connectorID)
			d.SetId("")
			return nil
		}
		return err
	}

	items := []string{}
	if assignments != nil && assignments.Items != nil {
		items = assignments.Items
	}

	if err := d.Set("connector_id", connectorID); err != nil {
		return err
	}
	if err := d.Set("operation_ids", items); err != nil {
		return err
	}

	d.SetId(connectorID)
	return nil
}

func resourceConnectorAssignmentUpdate(d *schema.ResourceData, m interface{}) error {
	connectorID := d.Get("connector_id").(string)
	log.Printf("[INFO] Updating ThousandEyes Connector Assignment for connector %s", connectorID)

	if _, err := setConnectorOperations(m.(*client.APIClient), connectorID, expandOperationIDs(d)); err != nil {
		return err
	}

	return resourceConnectorAssignmentRead(d, m)
}

func resourceConnectorAssignmentDelete(d *schema.ResourceData, m interface{}) error {
	connectorID := d.Get("connector_id").(string)
	if connectorID == "" {
		connectorID = d.Id()
	}

	log.Printf("[INFO] Deleting ThousandEyes Connector Assignment for connector %s", connectorID)

	if _, err := setConnectorOperations(m.(*client.APIClient), connectorID, []string{}); err != nil && !IsNotFoundError(err) {
		return err
	}

	d.SetId("")
	return nil
}

func expandOperationIDs(d *schema.ResourceData) []string {
	raw := d.Get("operation_ids").(*schema.Set).List()
	operationIDs := make([]string, 0, len(raw))
	for _, item := range raw {
		operationIDs = append(operationIDs, item.(string))
	}
	return operationIDs
}

func getConnectorOperations(apiClient *client.APIClient, connectorID string) (*connectors.Assignments, error) {
	api := (*connectors.GenericConnectorsAPIService)(&apiClient.Common)

	req := api.ListGenericConnectorOperations(connectorID)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	return resp, err
}

func setConnectorOperations(apiClient *client.APIClient, connectorID string, operationIDs []string) (*connectors.Assignments, error) {
	api := (*connectors.GenericConnectorsAPIService)(&apiClient.Common)

	req := api.SetGenericConnectorOperations(connectorID).RequestBody(operationIDs)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	return resp, err
}
