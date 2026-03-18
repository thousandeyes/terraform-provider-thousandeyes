package thousandeyes

import (
	"testing"

	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

func TestFlattenWebhookOperationHeadersUsesRemoteValue(t *testing.T) {
	headers := []connectors.Header{
		{Name: "Authorization", Value: "Bearer remote-token"},
	}

	got := flattenWebhookOperationHeaders(headers)
	if len(got) != 1 {
		t.Fatalf("expected 1 header, got %d", len(got))
	}
	headerMap := got[0].(map[string]interface{})
	if headerMap["name"] != "Authorization" {
		t.Fatalf("unexpected header name: %#v", headerMap["name"])
	}
	if headerMap["value"] != "Bearer remote-token" {
		t.Fatalf("expected remote value, got %#v", headerMap["value"])
	}
}

func TestFlattenWebhookOperationHeadersUsesRemoteShapeForDriftDetection(t *testing.T) {
	headers := []connectors.Header{
		{Name: "X-New-Header", Value: "*****"},
	}

	got := flattenWebhookOperationHeaders(headers)
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

func TestFlattenWebhookOperationHeadersKeepsRemoteValueWhenNotMasked(t *testing.T) {
	headers := []connectors.Header{
		{Name: "X-Trace-Id", Value: "remote-visible"},
	}

	got := flattenWebhookOperationHeaders(headers)
	if len(got) != 1 {
		t.Fatalf("expected 1 header, got %d", len(got))
	}
	headerMap := got[0].(map[string]interface{})
	if headerMap["value"] != "remote-visible" {
		t.Fatalf("expected remote value, got %#v", headerMap["value"])
	}
}
