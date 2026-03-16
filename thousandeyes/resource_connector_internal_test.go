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
