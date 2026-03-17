package thousandeyes

import (
	"context"
	"fmt"
	"log"
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
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
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
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
	resp, httpResp, err := req.Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			log.Printf("[INFO] Connector %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
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
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}

	return resourceConnectorRead(d, m)
}

func resourceConnectorDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*connectors.GenericConnectorsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Connector %s", d.Id())

	req := api.DeleteGenericConnector(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
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
	if connector == nil {
		return fmt.Errorf("connector response is nil")
	}

	var priorHeaders []interface{}
	if v, ok := d.Get("headers").([]interface{}); ok {
		priorHeaders = v
	}

	var priorAuthentication []interface{}
	if v, ok := d.Get("authentication").([]interface{}); ok {
		priorAuthentication = v
	}

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
		ts := time.UnixMilli(*connector.LastModifiedDate).UTC().Format(time.RFC3339)
		if err := d.Set("last_modified_date", ts); err != nil {
			return err
		}
	} else {
		if err := d.Set("last_modified_date", nil); err != nil {
			return err
		}
	}

	if len(connector.Headers) > 0 {
		if err := d.Set("headers", flattenConnectorHeaders(connector.Headers, priorHeaders)); err != nil {
			return err
		}
	} else {
		if err := d.Set("headers", nil); err != nil {
			return err
		}
	}

	if auth := flattenConnectorAuthentication(connector.Authentication, priorAuthentication); auth != nil {
		if err := d.Set("authentication", auth); err != nil {
			return err
		}
	} else {
		if err := d.Set("authentication", nil); err != nil {
			return err
		}
	}

	return nil
}

func flattenConnectorHeaders(headers []connectors.Header, prior []interface{}) []interface{} {
	priorByName := map[string]string{}
	for _, raw := range prior {
		headerMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := headerMap["name"].(string)
		value, _ := headerMap["value"].(string)
		if name != "" && value != "" {
			priorByName[name] = value
		}
	}

	out := make([]interface{}, 0, len(headers))
	for _, header := range headers {
		value := header.Value
		if prior, ok := priorByName[header.Name]; ok && prior != "" {
			value = prior
		}
		out = append(out, map[string]interface{}{
			"name":  header.Name,
			"value": value,
		})
	}
	return out
}

func flattenConnectorAuthentication(auth *connectors.GenericConnectorAuth, prior []interface{}) []interface{} {
	if auth == nil {
		return nil
	}

	authMap := map[string]interface{}{}

	switch {
	case auth.BasicAuthentication != nil:
		a := auth.BasicAuthentication
		authMap["type"] = string(a.Type)
		authMap["username"] = a.Username
		authMap["password"] = a.Password
	case auth.BearerTokenAuthentication != nil:
		a := auth.BearerTokenAuthentication
		authMap["type"] = string(a.Type)
		authMap["token"] = a.Token
	case auth.OtherTokenAuthentication != nil:
		a := auth.OtherTokenAuthentication
		authMap["type"] = string(a.Type)
		authMap["token"] = a.Token
	case auth.OauthClientCredentialsAuthentication != nil:
		a := auth.OauthClientCredentialsAuthentication
		authMap["type"] = string(a.Type)
		authMap["oauth_client_id"] = a.OauthClientId
		authMap["oauth_token_url"] = a.OauthTokenUrl
		authMap["oauth_client_secret"] = a.OauthClientSecret
		if a.Token != nil {
			authMap["token"] = *a.Token
		}
	case auth.OauthCodeAuthentication != nil:
		a := auth.OauthCodeAuthentication
		authMap["type"] = string(a.Type)
		authMap["oauth_client_id"] = a.OauthClientId
		authMap["oauth_token_url"] = a.OauthTokenUrl
		authMap["oauth_auth_url"] = a.OauthAuthUrl
		authMap["oauth_client_secret"] = a.OauthClientSecret
		authMap["code"] = a.Code
		authMap["redirect_uri"] = a.RedirectUri
		if a.Token != nil {
			authMap["token"] = *a.Token
		}
		if a.RefreshToken != nil {
			authMap["refresh_token"] = *a.RefreshToken
		}
	default:
		return nil
	}

	if priorMap := firstAuthMap(prior); priorMap != nil {
		currentType, _ := authMap["type"].(string)
		priorType, _ := priorMap["type"].(string)
		if currentType == priorType {
			for key, remoteValue := range authMap {
				if key == "type" {
					continue
				}
				remoteString, ok := remoteValue.(string)
				if !ok || (!shouldFallbackAuthField(remoteString)) {
					continue
				}
				if priorValue, ok := priorMap[key].(string); ok && priorValue != "" {
					authMap[key] = priorValue
				}
			}
		}
	}

	return []interface{}{authMap}
}

func shouldFallbackAuthField(remoteValue string) bool {
	if remoteValue == "" {
		return true
	}

	// Connector API can obfuscate auth fields (for example "*****").
	for _, c := range remoteValue {
		if c != '*' {
			return false
		}
	}
	return true
}

func firstAuthMap(prior []interface{}) map[string]interface{} {
	if len(prior) == 0 {
		return nil
	}
	authMap, ok := prior[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return authMap
}
