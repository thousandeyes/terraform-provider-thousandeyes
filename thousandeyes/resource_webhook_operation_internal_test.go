package thousandeyes

import (
	"testing"

	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

func TestFlattenWebhookOperationHeadersUsesPriorValue(t *testing.T) {
	headers := []connectors.Header{
		{Name: "Authorization", Value: "*****"},
	}
	prior := []interface{}{
		map[string]interface{}{
			"name":  "Authorization",
			"value": "Bearer local-token",
		},
	}

	got := flattenWebhookOperationHeaders(headers, prior)
	if len(got) != 1 {
		t.Fatalf("expected 1 header, got %d", len(got))
	}
	headerMap := got[0].(map[string]interface{})
	if headerMap["name"] != "Authorization" {
		t.Fatalf("unexpected header name: %#v", headerMap["name"])
	}
	if headerMap["value"] != "Bearer local-token" {
		t.Fatalf("expected prior value, got %#v", headerMap["value"])
	}
}

func TestFlattenWebhookOperationHeadersUsesRemoteShapeForDriftDetection(t *testing.T) {
	headers := []connectors.Header{
		{Name: "X-New-Header", Value: "*****"},
	}

	got := flattenWebhookOperationHeaders(headers, nil)
	if len(got) != 1 {
		t.Fatalf("expected 1 header, got %d", len(got))
	}
	headerMap := got[0].(map[string]interface{})
	if headerMap["name"] != "X-New-Header" {
		t.Fatalf("expected remote header name, got %#v", headerMap["name"])
	}
	if headerMap["value"] != "*****" {
		t.Fatalf("expected remote masked value for unknown header, got %#v", headerMap["value"])
	}
}

func TestFlattenWebhookOperationHeadersKeepsEmptyPriorValue(t *testing.T) {
	headers := []connectors.Header{
		{Name: "X-Trace-Id", Value: "*****"},
	}
	prior := []interface{}{
		map[string]interface{}{
			"name":  "X-Trace-Id",
			"value": "",
		},
	}

	got := flattenWebhookOperationHeaders(headers, prior)
	if len(got) != 1 {
		t.Fatalf("expected 1 header, got %d", len(got))
	}
	headerMap := got[0].(map[string]interface{})
	if headerMap["value"] != "" {
		t.Fatalf("expected empty prior value, got %#v", headerMap["value"])
	}
}
