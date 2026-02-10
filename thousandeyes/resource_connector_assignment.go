package thousandeyes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

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
		Description: "Manages all connector assignments for a webhook operation.\n\nThis resource is authoritative for the operation and uses PUT replace-all semantics:\n- To add a connector, Terraform sends the full list including the new ID.\n- To remove a connector, Terraform sends the full list without that ID.\n- To remove all connectors, Terraform sends an empty list (`[]`).\n\nAny connector assignments made outside Terraform for the same webhook operation will be removed on the next apply if they are not present in `connector_ids`.",
		Importer: &schema.ResourceImporter{
			StateContext: resourceConnectorAssignmentImport,
		},
	}
}

func resourceConnectorAssignmentImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	if err := d.Set("webhook_operation_id", d.Id()); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
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
	return executeWebhookOperationConnectorsRequest(apiClient, http.MethodGet, operationID, nil)
}

func setWebhookOperationConnectors(apiClient *client.APIClient, operationID string, connectorIDs []string) (*connectors.Assignments, error) {
	return executeWebhookOperationConnectorsRequest(apiClient, http.MethodPut, operationID, connectorIDs)
}

func executeWebhookOperationConnectorsRequest(apiClient *client.APIClient, method string, operationID string, requestBody interface{}) (*connectors.Assignments, error) {
	basePath := strings.TrimRight(apiClient.GetConfig().ServerURL, "/")
	path := fmt.Sprintf("%s/operations/%s/%s/connectors", basePath, operationTypeWebhooks, url.PathEscape(operationID))

	headers := map[string]string{
		"Accept": "application/problem+json, application/hal+json, application/json",
	}
	if method == http.MethodPut {
		headers["Content-Type"] = "application/json"
	}

	queryParams := url.Values{}
	if aid, ok := aidStringFromContext(apiClient.GetConfig().Context); ok {
		queryParams.Set("aid", aid)
	}

	req, err := apiClient.PrepareRequest(path, method, requestBody, headers, queryParams, url.Values{})
	if err != nil {
		return nil, err
	}

	httpResp, err := apiClient.CallAPI(req)
	if err != nil || httpResp == nil {
		return nil, err
	}

	body, err := io.ReadAll(httpResp.Body)
	httpResp.Body.Close()
	httpResp.Body = io.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("operation connectors request failed: %s: %s", httpResp.Status, strings.TrimSpace(string(body)))
	}

	resp := &connectors.Assignments{}
	if len(body) == 0 {
		return resp, nil
	}
	if err := apiClient.Decode(resp, body, httpResp.Header.Get("Content-Type")); err != nil {
		return nil, err
	}

	return resp, nil
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
