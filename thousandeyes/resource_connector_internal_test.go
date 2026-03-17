package thousandeyes

import (
	"testing"

	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

func TestFlattenConnectorHeadersPreservesExistingSensitiveValues(t *testing.T) {
	headers := []connectors.Header{
		{Name: "X-Trace-Id", Value: "********"},
	}
	prior := []interface{}{
		map[string]interface{}{
			"name":  "X-Trace-Id",
			"value": "local-secret-value",
		},
	}

	got := flattenConnectorHeaders(headers, prior)
	if len(got) != 1 {
		t.Fatalf("expected 1 header, got %d", len(got))
	}
	headerMap := got[0].(map[string]interface{})
	if headerMap["name"] != "X-Trace-Id" {
		t.Fatalf("unexpected header name: %#v", headerMap["name"])
	}
	if headerMap["value"] != "local-secret-value" {
		t.Fatalf("expected preserved local header value, got %#v", headerMap["value"])
	}
}

func TestFlattenConnectorAuthenticationUsesServerFieldsAndKeepsSecret(t *testing.T) {
	auth := &connectors.GenericConnectorAuth{
		BasicAuthentication: &connectors.BasicAuthentication{
			Type:     connectors.AUTHENTICATIONTYPE_BASIC,
			Username: "remote-user",
			Password: "********",
		},
	}
	prior := []interface{}{
		map[string]interface{}{
			"type":     "basic",
			"username": "old-user",
			"password": "local-secret",
		},
	}

	got := flattenConnectorAuthentication(auth, prior)
	if len(got) != 1 {
		t.Fatalf("expected 1 auth block, got %d", len(got))
	}
	authMap := got[0].(map[string]interface{})
	if authMap["type"] != "basic" {
		t.Fatalf("unexpected auth type: %#v", authMap["type"])
	}
	if authMap["username"] != "remote-user" {
		t.Fatalf("expected server username, got %#v", authMap["username"])
	}
	if authMap["password"] != "local-secret" {
		t.Fatalf("expected preserved local password, got %#v", authMap["password"])
	}
}

func TestFlattenConnectorAuthenticationFallsBackWhenValuesAreMaskedOrMissing(t *testing.T) {
	auth := &connectors.GenericConnectorAuth{
		OauthCodeAuthentication: &connectors.OauthCodeAuthentication{
			Type:              connectors.AUTHENTICATIONTYPE_OAUTH_AUTH_CODE,
			OauthClientId:     "client-id",
			OauthClientSecret: "*****",
			OauthTokenUrl:     "https://auth.example.com/token",
			OauthAuthUrl:      "https://auth.example.com/authorize",
			Code:              "*****",
			RedirectUri:       "",
		},
	}
	prior := []interface{}{
		map[string]interface{}{
			"type":                "oauth-auth-code",
			"oauth_client_id":     "client-id",
			"oauth_client_secret": "local-client-secret",
			"oauth_token_url":     "https://auth.example.com/token",
			"oauth_auth_url":      "https://auth.example.com/authorize",
			"code":                "local-auth-code",
			"redirect_uri":        "https://app.example.com/callback",
		},
	}

	got := flattenConnectorAuthentication(auth, prior)
	if len(got) != 1 {
		t.Fatalf("expected 1 auth block, got %d", len(got))
	}
	authMap := got[0].(map[string]interface{})
	if authMap["oauth_client_secret"] != "local-client-secret" {
		t.Fatalf("expected fallback client secret, got %#v", authMap["oauth_client_secret"])
	}
	if authMap["code"] != "local-auth-code" {
		t.Fatalf("expected fallback code, got %#v", authMap["code"])
	}
	if authMap["redirect_uri"] != "https://app.example.com/callback" {
		t.Fatalf("expected fallback redirect uri, got %#v", authMap["redirect_uri"])
	}
}

func TestShouldFallbackAuthField(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "empty", value: "", want: true},
		{name: "masked", value: "*****", want: true},
		{name: "plain", value: "testuser", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldFallbackAuthField(tt.value)
			if got != tt.want {
				t.Fatalf("shouldFallbackAuthField(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
