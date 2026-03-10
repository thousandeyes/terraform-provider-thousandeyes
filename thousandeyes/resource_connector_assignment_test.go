package thousandeyes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func TestAccThousandEyesConnectorAssignment(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	testAccPreCheck(t)

	operationID1, err := createAcceptanceWebhookOperation(testClient)
	if err != nil {
		t.Fatalf("failed creating first acceptance webhook operation: %v", err)
	}

	operationID2, err := createAcceptanceWebhookOperation(testClient)
	if err != nil {
		deleteAcceptanceWebhookOperation(testClient, operationID1)
		t.Fatalf("failed creating second acceptance webhook operation: %v", err)
	}

	connectorID, err := createAcceptanceConnector(testClient)
	if err != nil {
		deleteAcceptanceWebhookOperation(testClient, operationID1)
		deleteAcceptanceWebhookOperation(testClient, operationID2)
		t.Fatalf("failed creating acceptance connector: %v", err)
	}

	defer func() {
		deleteAcceptanceConnector(testClient, connectorID)
		deleteAcceptanceWebhookOperation(testClient, operationID1)
		deleteAcceptanceWebhookOperation(testClient, operationID2)
	}()

	resourceName := "thousandeyes_connector_assignment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(_ *terraform.State) error {
			assignments, err := getConnectorOperations(testClient, connectorID)
			if err != nil {
				return err
			}
			if assignments != nil && len(assignments.Items) != 0 {
				return fmt.Errorf("expected no operation assignments after destroy, found %d", len(assignments.Items))
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConnectorAssignmentConfig(
					"acceptance_resources/connector_assignment/basic.tf",
					connectorID,
					operationID1,
					operationID2,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "connector_id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "operation_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "operation_ids.*", operationID1),
				),
			},
			{
				Config: testAccConnectorAssignmentConfig(
					"acceptance_resources/connector_assignment/update.tf",
					connectorID,
					operationID1,
					operationID2,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "connector_id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "operation_ids.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "operation_ids.*", operationID1),
					resource.TestCheckTypeSetElemAttr(resourceName, "operation_ids.*", operationID2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccThousandEyesConnectorAssignment_threeWebhookOperations(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	testAccPreCheck(t)

	operationID1, err := createAcceptanceWebhookOperation(testClient)
	if err != nil {
		t.Fatalf("failed creating first acceptance webhook operation: %v", err)
	}

	operationID2, err := createAcceptanceWebhookOperation(testClient)
	if err != nil {
		deleteAcceptanceWebhookOperation(testClient, operationID1)
		t.Fatalf("failed creating second acceptance webhook operation: %v", err)
	}

	operationID3, err := createAcceptanceWebhookOperation(testClient)
	if err != nil {
		deleteAcceptanceWebhookOperation(testClient, operationID1)
		deleteAcceptanceWebhookOperation(testClient, operationID2)
		t.Fatalf("failed creating third acceptance webhook operation: %v", err)
	}

	connectorID, err := createAcceptanceConnector(testClient)
	if err != nil {
		deleteAcceptanceWebhookOperation(testClient, operationID1)
		deleteAcceptanceWebhookOperation(testClient, operationID2)
		deleteAcceptanceWebhookOperation(testClient, operationID3)
		t.Fatalf("failed creating acceptance connector: %v", err)
	}

	defer func() {
		deleteAcceptanceConnector(testClient, connectorID)
		deleteAcceptanceWebhookOperation(testClient, operationID1)
		deleteAcceptanceWebhookOperation(testClient, operationID2)
		deleteAcceptanceWebhookOperation(testClient, operationID3)
	}()

	resourceName := "thousandeyes_connector_assignment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(_ *terraform.State) error {
			assignments, err := getConnectorOperations(testClient, connectorID)
			if err != nil {
				return err
			}
			if assignments != nil && len(assignments.Items) != 0 {
				return fmt.Errorf("expected no operation assignments after destroy, found %d", len(assignments.Items))
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConnectorAssignmentThreeOperationsConfig(connectorID, operationID1, operationID2, operationID3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "connector_id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "operation_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "operation_ids.*", operationID1),
					resource.TestCheckTypeSetElemAttr(resourceName, "operation_ids.*", operationID2),
					resource.TestCheckTypeSetElemAttr(resourceName, "operation_ids.*", operationID3),
				),
			},
		},
	})
}

func testAccConnectorAssignmentConfig(resourceFile, connectorID, operationID1, operationID2 string) string {
	content, err := os.ReadFile(resourceFile)
	if err != nil {
		panic(err)
	}

	cfg := string(content)
	prefix := fmt.Sprintf(`
locals {
  connector_id   = %q
  operation_id_1 = %q
  operation_id_2 = %q
}
`, connectorID, operationID1, operationID2)

	return prefix + "\n" + cfg
}

func testAccConnectorAssignmentThreeOperationsConfig(connectorID, operationID1, operationID2, operationID3 string) string {
	return fmt.Sprintf(`
resource "thousandeyes_connector_assignment" "test" {
  connector_id = %q
  operation_ids = [
    %q,
    %q,
    %q
  ]
}
`, connectorID, operationID1, operationID2, operationID3)
}

func createAcceptanceWebhookOperation(apiClient *client.APIClient) (string, error) {
	name := fmt.Sprintf("UAT Connector Assignment Operation %d", time.Now().UnixNano())
	body := map[string]interface{}{
		"name":     name,
		"category": "alerts",
		"status":   "pending",
		"enabled":  false,
		"type":     "webhook",
		"headers": []map[string]string{
			{
				"name":  "Content-Type",
				"value": "application/json",
			},
		},
	}

	var operation struct {
		ID *string `json:"id"`
	}
	if err := executeAcceptanceRequest(apiClient, http.MethodPost, "/operations/webhooks", body, &operation); err != nil {
		return "", err
	}
	if operation.ID == nil || *operation.ID == "" {
		return "", fmt.Errorf("webhook operation creation returned empty id")
	}
	return *operation.ID, nil
}

func deleteAcceptanceWebhookOperation(apiClient *client.APIClient, operationID string) {
	if operationID == "" {
		return
	}
	if err := executeAcceptanceRequest(apiClient, http.MethodDelete, fmt.Sprintf("/operations/webhooks/%s", url.PathEscape(operationID)), nil, nil); err != nil {
		log.Printf("[WARN] failed deleting acceptance webhook operation %s: %v", operationID, err)
	}
}

func createAcceptanceConnector(apiClient *client.APIClient) (string, error) {
	name := fmt.Sprintf("UAT Connector Assignment Connector %d", time.Now().UnixNano())
	body := map[string]interface{}{
		"name":   name,
		"target": "https://webhook.site/terraform-provider-thousandeyes",
		"type":   "generic",
	}

	var connector struct {
		ID *string `json:"id"`
	}
	if err := executeAcceptanceRequest(apiClient, http.MethodPost, "/connectors/generic", body, &connector); err != nil {
		return "", err
	}
	if connector.ID == nil || *connector.ID == "" {
		return "", fmt.Errorf("connector creation returned empty id")
	}
	return *connector.ID, nil
}

func deleteAcceptanceConnector(apiClient *client.APIClient, connectorID string) {
	if connectorID == "" {
		return
	}
	if err := executeAcceptanceRequest(apiClient, http.MethodDelete, fmt.Sprintf("/connectors/generic/%s", url.PathEscape(connectorID)), nil, nil); err != nil {
		log.Printf("[WARN] failed deleting acceptance connector %s: %v", connectorID, err)
	}
}

func executeAcceptanceRequest(apiClient *client.APIClient, method, relativePath string, requestBody interface{}, out interface{}) error {
	base := strings.TrimRight(apiClient.GetConfig().ServerURL, "/")
	path := base + relativePath

	headers := map[string]string{
		"Accept": "application/problem+json, application/hal+json, application/json",
	}
	if requestBody != nil {
		headers["Content-Type"] = "application/json"
	}

	queryParams := url.Values{}
	if aid, ok := aidStringFromContext(apiClient.GetConfig().Context); ok {
		queryParams.Set("aid", aid)
	}

	req, err := apiClient.PrepareRequest(path, method, requestBody, headers, queryParams, url.Values{})
	if err != nil {
		return err
	}

	httpResp, err := apiClient.CallAPI(req)
	if err != nil || httpResp == nil {
		return err
	}

	body, err := io.ReadAll(httpResp.Body)
	httpResp.Body.Close()
	httpResp.Body = io.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if httpResp.StatusCode >= 300 {
		return fmt.Errorf("acceptance request failed: %s: %s", httpResp.Status, strings.TrimSpace(string(body)))
	}

	if out == nil || len(body) == 0 {
		return nil
	}

	return apiClient.Decode(out, body, httpResp.Header.Get("Content-Type"))
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
