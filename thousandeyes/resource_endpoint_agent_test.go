package thousandeyes

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceEndpointAgentRead_ImportAndRefresh(t *testing.T) {
	const agentID = "agent-123"
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet || r.URL.Path != "/v7/endpoint/agents/"+agentID {
			t.Errorf("expected GET for imported agent, got %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("aid"); got != "12345" {
			t.Errorf("expected account group query aid=12345, got %q", got)
		}
		return endpointAgentJSONResponse(http.StatusOK, `{"id":"agent-123","name":"Alice laptop","licenseType":"advantage","status":"enabled"}`), nil
	})
	apiClient := newEndpointAgentTestClient(transport, "12345")
	d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{})
	d.SetId(agentID)

	if err := resourceEndpointAgentRead(d, apiClient); err != nil {
		t.Fatalf("read after import failed: %v", err)
	}
	if d.Id() != agentID {
		t.Fatalf("expected imported ID %q, got %q", agentID, d.Id())
	}
	for key, want := range map[string]string{"name": "Alice laptop", "license_type": "advantage", "status": "enabled"} {
		if got := d.Get(key).(string); got != want {
			t.Errorf("expected %s %q, got %q", key, want, got)
		}
	}
}

func TestResourceEndpointAgentRead_PreservesWirelessLicenseTypes(t *testing.T) {
	for _, license := range []string{"embed-wireless-advantage", "embed-wireless-essentials"} {
		t.Run(license, func(t *testing.T) {
			transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				return endpointAgentJSONResponse(http.StatusOK, `{"id":"agent-123","name":"Wireless agent","licenseType":"`+license+`","status":"disabled"}`), nil
			})
			d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{})
			d.SetId("agent-123")
			if err := resourceEndpointAgentRead(d, newEndpointAgentTestClient(transport, "12345")); err != nil {
				t.Fatalf("expected wireless license to be readable, got %v", err)
			}
			if got := d.Get("license_type").(string); got != license {
				t.Fatalf("expected license_type %q, got %q", license, got)
			}
		})
	}
}

func TestResourceEndpointAgentRead_NotFoundRemovesState(t *testing.T) {
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		return endpointAgentJSONResponse(http.StatusNotFound, `{"error":"not found"}`), nil
	})
	d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{})
	d.SetId("missing-agent")
	if err := resourceEndpointAgentRead(d, newEndpointAgentTestClient(transport, "12345")); err != nil {
		t.Fatalf("expected missing remote agent to be removed from state, got %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("expected state ID to be cleared, got %q", d.Id())
	}
}

func TestResourceEndpointAgentUpdate_UpdatesNameAndSupportedLicense(t *testing.T) {
	var patchCount, readCount int
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodPatch && r.URL.Path == "/v7/endpoint/agents/agent-123":
			patchCount++
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read update body: %v", err)
			}
			var payload map[string]string
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatalf("decode update body %q: %v", body, err)
			}
			if payload["name"] != "Renamed laptop" || payload["licenseType"] != "essentials" {
				t.Errorf("expected name and supported license update, got %s", body)
			}
			return endpointAgentJSONResponse(http.StatusOK, `{"id":"agent-123","name":"Renamed laptop","licenseType":"essentials","status":"enabled"}`), nil
		case r.Method == http.MethodGet && r.URL.Path == "/v7/endpoint/agents/agent-123":
			readCount++
			return endpointAgentJSONResponse(http.StatusOK, `{"id":"agent-123","name":"Renamed laptop","licenseType":"essentials","status":"enabled"}`), nil
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			return endpointAgentJSONResponse(http.StatusInternalServerError, `{}`), nil
		}
	})
	d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{
		"name":         "Renamed laptop",
		"license_type": "essentials",
	})
	d.SetId("agent-123")
	if err := resourceEndpointAgentUpdate(d, newEndpointAgentTestClient(transport, "12345")); err != nil {
		t.Fatalf("expected supported property update to succeed, got %v", err)
	}
	if patchCount != 1 {
		t.Fatalf("expected one update request, got %d", patchCount)
	}
	if readCount == 0 {
		t.Fatal("expected the resource to refresh remote properties after update")
	}
}

func TestResourceEndpointAgentUpdate_RejectsWirelessLicenseAssignmentBeforeRequest(t *testing.T) {
	for _, license := range []string{"embed-wireless-advantage", "embed-wireless-essentials"} {
		t.Run(license, func(t *testing.T) {
			requests := 0
			transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				requests++
				return endpointAgentJSONResponse(http.StatusOK, `{}`), nil
			})
			d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{
				"license_type": license,
			})
			d.SetId("agent-123")
			err := resourceEndpointAgentUpdate(d, newEndpointAgentTestClient(transport, "12345"))
			if err == nil || !strings.Contains(err.Error(), license) {
				t.Fatalf("expected clear diagnostic naming unsupported license %q, got %v", license, err)
			}
			if requests != 0 {
				t.Fatalf("unsupported license must be rejected before an API request, got %d", requests)
			}
		})
	}
}

func TestResourceEndpointAgentUpdate_ChangesStatus(t *testing.T) {
	for _, tc := range []struct {
		status string
		path   string
	}{
		{status: "enabled", path: "/v7/endpoint/agents/agent-123/enable"},
		{status: "disabled", path: "/v7/endpoint/agents/agent-123/disable"},
	} {
		t.Run(tc.status, func(t *testing.T) {
			var actionCount int
			transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				switch {
				case r.Method == http.MethodPost && r.URL.Path == tc.path:
					actionCount++
					return endpointAgentJSONResponse(http.StatusOK, `{"id":"agent-123","name":"Laptop","licenseType":"advantage","status":"`+tc.status+`"}`), nil
				case r.Method == http.MethodGet && r.URL.Path == "/v7/endpoint/agents/agent-123":
					return endpointAgentJSONResponse(http.StatusOK, `{"id":"agent-123","name":"Laptop","licenseType":"advantage","status":"`+tc.status+`"}`), nil
				default:
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
					return endpointAgentJSONResponse(http.StatusInternalServerError, `{}`), nil
				}
			})
			d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{"status": tc.status})
			d.SetId("agent-123")
			if err := resourceEndpointAgentUpdate(d, newEndpointAgentTestClient(transport, "12345")); err != nil {
				t.Fatalf("status update failed: %v", err)
			}
			if actionCount != 1 {
				t.Fatalf("expected one status operation at %s, got %d", tc.path, actionCount)
			}
			if got := d.Get("status").(string); got != tc.status {
				t.Fatalf("expected refreshed status %q, got %q", tc.status, got)
			}
		})
	}
}

func TestResourceEndpointAgentDelete_ClearsStateAfterSuccess(t *testing.T) {
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v7/endpoint/agents/agent-123" {
			t.Errorf("expected DELETE for agent, got %s %s", r.Method, r.URL.Path)
		}
		return endpointAgentJSONResponse(http.StatusNoContent, ``), nil
	})
	d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{})
	d.SetId("agent-123")
	if err := resourceEndpointAgentDelete(d, newEndpointAgentTestClient(transport, "12345")); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("expected ID to clear after successful delete, got %q", d.Id())
	}
}

func TestResourceEndpointAgentDelete_RetainsStateAfterFailure(t *testing.T) {
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		return endpointAgentJSONResponse(http.StatusInternalServerError, `{"error":"temporary failure"}`), nil
	})
	d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{})
	d.SetId("agent-123")
	if err := resourceEndpointAgentDelete(d, newEndpointAgentTestClient(transport, "12345")); err == nil {
		t.Fatal("expected API delete error")
	}
	if d.Id() != "agent-123" {
		t.Fatalf("expected ID to remain after failed delete, got %q", d.Id())
	}
}

func TestResourceEndpointAgentDelete_ProtectionRemovesStateWithoutAPIRequest(t *testing.T) {
	requests := 0
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests++
		return endpointAgentJSONResponse(http.StatusNoContent, ``), nil
	})
	d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{
		"deletion_protection": true,
	})
	d.SetId("agent-123")
	if err := resourceEndpointAgentDelete(d, newEndpointAgentTestClient(transport, "12345")); err != nil {
		t.Fatalf("protected delete should detach successfully, got %v", err)
	}
	if requests != 0 {
		t.Fatalf("deletion protection must not call the API, got %d requests", requests)
	}
	if d.Id() != "" {
		t.Fatalf("expected state ID to clear after protected delete, got %q", d.Id())
	}
}

func TestResourceEndpointAgentDelete_ProtectionDefaultsToFalse(t *testing.T) {
	resource := resourceEndpointAgent()
	if got := resource.Schema["deletion_protection"].Default; got != false {
		t.Fatalf("expected deletion protection default false, got %#v", got)
	}
}

func TestResourceEndpointAgentCreate_IsRejectedWithoutRequest(t *testing.T) {
	requests := 0
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests++
		return endpointAgentJSONResponse(http.StatusOK, `{}`), nil
	})
	d := schema.TestResourceDataRaw(t, resourceEndpointAgent().Schema, map[string]interface{}{"name": "New laptop"})
	err := resourceEndpointAgentCreate(d, newEndpointAgentTestClient(transport, "12345"))
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "import") {
		t.Fatalf("expected clear import-only create diagnostic, got %v", err)
	}
	if requests != 0 {
		t.Fatalf("create must not call the API, got %d requests", requests)
	}
}
