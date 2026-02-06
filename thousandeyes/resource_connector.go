package thousandeyes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

func resourceConnector() *schema.Resource {
	return &schema.Resource{
		Schema:        schemas.ConnectorSchema,
		Create:        resourceConnectorCreate,
		Read:          resourceConnectorRead,
		Update:        resourceConnectorUpdate,
		Delete:        resourceConnectorDelete,
		CustomizeDiff: validateConnectorAuthentication,
		Description:   "Manages a ThousandEyes connector. A connector defines where webhook notifications are sent, including the target URL and authentication configuration.",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectorCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*connectors.GenericConnectorsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Connector")

	connector := buildConnectorRequest(d)
	req := api.CreateGenericConnector().GenericConnector(*connector)
	req = SetAidFromContextFloat32(apiClient.GetConfig().Context, req)

	resp, httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode < 300 {
			if parsed, parseErr := decodeGenericConnectorFromResponse(httpResp); parseErr == nil && parsed != nil {
				connectorID, idErr := connectorIDOrError(parsed)
				if idErr != nil {
					return idErr
				}
				d.SetId(connectorID)
				return setConnectorResourceData(d, parsed)
			}
		}
		return err
	}

	connectorID, err := connectorIDOrError(resp)
	if err != nil {
		return err
	}
	d.SetId(connectorID)
	return resourceConnectorRead(d, m)
}

func resourceConnectorRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*connectors.GenericConnectorsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Reading ThousandEyes Connector %s", d.Id())

	req := api.GetGenericConnector(d.Id())
	req = SetAidFromContextFloat32(apiClient.GetConfig().Context, req)
	resp, httpResp, err := req.Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			log.Printf("[INFO] Connector %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		if httpResp != nil && httpResp.StatusCode < 300 {
			if parsed, parseErr := decodeGenericConnectorFromResponse(httpResp); parseErr == nil && parsed != nil {
				return setConnectorResourceData(d, parsed)
			}
		}
		return err
	}

	return setConnectorResourceData(d, resp)
}

func resourceConnectorUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*connectors.GenericConnectorsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Connector %s", d.Id())

	connector := buildConnectorRequest(d)
	req := api.UpdateGenericConnector(d.Id()).GenericConnector(*connector)
	req = SetAidFromContextFloat32(apiClient.GetConfig().Context, req)

	_, httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode < 300 {
			if parsed, parseErr := decodeGenericConnectorFromResponse(httpResp); parseErr == nil && parsed != nil {
				return setConnectorResourceData(d, parsed)
			}
		}
		return err
	}

	return resourceConnectorRead(d, m)
}

func resourceConnectorDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*connectors.GenericConnectorsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Connector %s", d.Id())

	req := api.DeleteGenericConnector(d.Id())
	req = SetAidFromContextFloat32(apiClient.GetConfig().Context, req)
	_, err := req.Execute()
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func buildConnectorRequest(d *schema.ResourceData) *connectors.GenericConnector {
	connector := &connectors.GenericConnector{
		Name:   d.Get("name").(string),
		Target: d.Get("target").(string),
		Type:   connectors.CONNECTORTYPE_GENERIC,
	}

	if v, ok := d.GetOk("headers"); ok {
		connector.Headers = expandHeaders(v.([]interface{}))
	}

	if v, ok := d.GetOk("authentication"); ok {
		authList := v.([]interface{})
		if len(authList) > 0 {
			connector.Authentication = expandAuthentication(authList[0].(map[string]interface{}))
		}
	}

	return connector
}

func expandHeaders(headersList []interface{}) []connectors.Header {
	headers := make([]connectors.Header, 0, len(headersList))
	for _, h := range headersList {
		headerMap := h.(map[string]interface{})
		headers = append(headers, connectors.Header{
			Name:  headerMap["name"].(string),
			Value: headerMap["value"].(string),
		})
	}
	return headers
}

func expandAuthentication(authMap map[string]interface{}) *connectors.GenericConnectorAuth {
	authType := getAuthString(authMap, "type")

	switch authType {
	case "basic":
		return &connectors.GenericConnectorAuth{
			BasicAuthentication: &connectors.BasicAuthentication{
				Type:     connectors.AuthenticationType(authType),
				Username: getAuthString(authMap, "username"),
				Password: getAuthString(authMap, "password"),
			},
		}
	case "bearer-token":
		return &connectors.GenericConnectorAuth{
			BearerTokenAuthentication: &connectors.BearerTokenAuthentication{
				Type:  connectors.AuthenticationType(authType),
				Token: getAuthString(authMap, "token"),
			},
		}
	case "other-token":
		return &connectors.GenericConnectorAuth{
			OtherTokenAuthentication: &connectors.OtherTokenAuthentication{
				Type:  connectors.AuthenticationType(authType),
				Token: getAuthString(authMap, "token"),
			},
		}
	case "oauth-client-credentials":
		return &connectors.GenericConnectorAuth{
			OauthClientCredentialsAuthentication: &connectors.OauthClientCredentialsAuthentication{
				Type:              connectors.AuthenticationType(authType),
				OauthClientId:     getAuthString(authMap, "oauth_client_id"),
				OauthClientSecret: getAuthString(authMap, "oauth_client_secret"),
				OauthTokenUrl:     getAuthString(authMap, "oauth_token_url"),
				Token:             getAuthStringPtr(authMap, "token"),
			},
		}
	case "oauth-auth-code":
		return &connectors.GenericConnectorAuth{
			OauthCodeAuthentication: &connectors.OauthCodeAuthentication{
				Type:              connectors.AuthenticationType(authType),
				OauthClientId:     getAuthString(authMap, "oauth_client_id"),
				OauthClientSecret: getAuthString(authMap, "oauth_client_secret"),
				OauthTokenUrl:     getAuthString(authMap, "oauth_token_url"),
				OauthAuthUrl:      getAuthString(authMap, "oauth_auth_url"),
				Code:              getAuthString(authMap, "code"),
				RedirectUri:       getAuthString(authMap, "redirect_uri"),
				Token:             getAuthStringPtr(authMap, "token"),
				RefreshToken:      getAuthStringPtr(authMap, "refresh_token"),
			},
		}
	}
	return nil
}

func getAuthString(authMap map[string]interface{}, key string) string {
	if v, ok := authMap[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getAuthStringPtr(authMap map[string]interface{}, key string) *string {
	if v := getAuthString(authMap, key); v != "" {
		return &v
	}
	return nil
}

func connectorIDOrError(connector *connectors.GenericConnector) (string, error) {
	if connector == nil || connector.Id == nil || *connector.Id == "" {
		return "", fmt.Errorf("create connector response missing id")
	}
	return *connector.Id, nil
}

func decodeGenericConnectorFromResponse(httpResp *http.Response) (*connectors.GenericConnector, error) {
	if httpResp == nil || httpResp.Body == nil {
		return nil, fmt.Errorf("no response body to decode")
	}
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	httpResp.Body = io.NopCloser(bytes.NewBuffer(body))
	return decodeGenericConnectorFromBody(body)
}

func decodeGenericConnectorFromBody(body []byte) (*connectors.GenericConnector, error) {
	type genericConnectorFallback struct {
		ID               *string                  `json:"id"`
		Type             connectors.ConnectorType `json:"type"`
		Name             string                   `json:"name"`
		Target           string                   `json:"target"`
		LastModifiedDate json.RawMessage          `json:"lastModifiedDate"`
	}

	var raw genericConnectorFallback
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	connector := &connectors.GenericConnector{
		Id:     raw.ID,
		Type:   raw.Type,
		Name:   raw.Name,
		Target: raw.Target,
	}

	if ts, ok := parseConnectorLastModifiedDate(raw.LastModifiedDate); ok {
		connector.LastModifiedDate = &ts
	}

	return connector, nil
}

func parseConnectorLastModifiedDate(raw json.RawMessage) (time.Time, bool) {
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return time.Time{}, false
	}

	var rfc3339 string
	if err := json.Unmarshal(raw, &rfc3339); err == nil {
		parsed, parseErr := time.Parse(time.RFC3339, rfc3339)
		if parseErr != nil {
			return time.Time{}, false
		}
		return parsed, true
	}

	var millis int64
	if err := json.Unmarshal(raw, &millis); err == nil {
		return time.UnixMilli(millis), true
	}

	var asFloat float64
	if err := json.Unmarshal(raw, &asFloat); err == nil {
		return time.UnixMilli(int64(asFloat)), true
	}

	return time.Time{}, false
}

func validateConnectorAuthentication(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	v, ok := d.GetOk("authentication")
	if !ok {
		return nil
	}
	authList, ok := v.([]interface{})
	if !ok || len(authList) == 0 {
		return nil
	}
	authMap, ok := authList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	authType := getAuthString(authMap, "type")
	if authType == "" {
		return nil
	}

	var required []string
	switch authType {
	case "basic":
		required = []string{"username", "password"}
	case "bearer-token", "other-token":
		required = []string{"token"}
	case "oauth-client-credentials":
		required = []string{"oauth_client_id", "oauth_client_secret", "oauth_token_url"}
	case "oauth-auth-code":
		required = []string{"oauth_client_id", "oauth_client_secret", "oauth_token_url", "oauth_auth_url", "code", "redirect_uri"}
	default:
		return nil
	}

	var missing []string
	for _, key := range required {
		if getAuthString(authMap, key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("authentication.type %q requires %s", authType, strings.Join(missing, ", "))
	}

	return nil
}

func setConnectorResourceData(d *schema.ResourceData, connector *connectors.GenericConnector) error {
	if err := d.Set("name", connector.Name); err != nil {
		return err
	}
	if err := d.Set("target", connector.Target); err != nil {
		return err
	}
	if err := d.Set("type", string(connector.Type)); err != nil {
		return err
	}

	if connector.LastModifiedDate != nil {
		if err := d.Set("last_modified_date", connector.LastModifiedDate.Format(time.RFC3339)); err != nil {
			return err
		}
	} else {
		if err := d.Set("last_modified_date", nil); err != nil {
			return err
		}
	}

	return nil
}
