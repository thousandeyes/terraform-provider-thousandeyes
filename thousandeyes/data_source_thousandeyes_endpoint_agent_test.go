package thousandeyes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func TestDataSourceEndpointAgentRead_Success(t *testing.T) {
	const (
		agentName = "Alice's laptop"
		agentID   = "endpoint-agent-123"
		aid       = "12345"
	)

	var requestCount int
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		requestCount++
		if r.URL.Path != "/v7/endpoint/agents" {
			t.Errorf("expected path /v7/endpoint/agents, got %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("agentName"); got != agentName {
			t.Errorf("expected agentName query %q, got %q", agentName, got)
		}
		if got := r.URL.Query().Get("aid"); got != aid {
			t.Errorf("expected aid query %q, got %q", aid, got)
		}
		return endpointAgentJSONResponse(http.StatusOK, `{"agents":[{"id":"other-id","name":"Other laptop"},{"id":"endpoint-agent-123","name":"Alice's laptop"}]}`), nil
	})

	d := schema.TestResourceDataRaw(t, dataSourceThousandeyesEndpointAgent().Schema, map[string]interface{}{
		"agent_name": agentName,
	})
	err := dataSourceEndpointAgentRead(d, newEndpointAgentTestClient(transport, aid))
	if err != nil {
		t.Fatalf("expected successful lookup, got %v", err)
	}
	if requestCount != 1 {
		t.Fatalf("expected one API request, got %d", requestCount)
	}
	if d.Id() != agentID {
		t.Fatalf("expected Terraform ID %q, got %q", agentID, d.Id())
	}
	if got := d.Get("agent_id").(string); got != agentID {
		t.Fatalf("expected agent_id %q, got %q", agentID, got)
	}
	if got := d.Get("agent_name").(string); got != agentName {
		t.Fatalf("expected agent_name %q, got %q", agentName, got)
	}
}

func TestDataSourceEndpointAgentRead_PaginatesUntilExactMatch(t *testing.T) {
	const agentName = "Alice's laptop"
	var requests int
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests++
		if got := r.URL.Query().Get("agentName"); got != agentName {
			t.Errorf("request %d: expected agentName %q, got %q", requests, agentName, got)
		}
		switch requests {
		case 1:
			if got := r.URL.Query().Get("cursor"); got != "" {
				t.Errorf("first request should not have a cursor, got %q", got)
			}
			return endpointAgentJSONResponse(http.StatusOK, `{"agents":[{"id":"first-page-id","name":"Different laptop"}],"_links":{"next":{"href":"https://api.thousandeyes.com/v7/endpoint/agents?agentName=Alice%27s+laptop&cursor=page-two"}}}`), nil
		case 2:
			if got := r.URL.Query().Get("cursor"); got != "page-two" {
				t.Errorf("expected second-page cursor %q, got %q", "page-two", got)
			}
			return endpointAgentJSONResponse(http.StatusOK, `{"agents":[{"id":"endpoint-agent-456","name":"Alice's laptop"}]}`), nil
		default:
			t.Errorf("unexpected request %d", requests)
			return endpointAgentJSONResponse(http.StatusInternalServerError, `{}`), nil
		}
	})

	d := schema.TestResourceDataRaw(t, dataSourceThousandeyesEndpointAgent().Schema, map[string]interface{}{
		"agent_name": agentName,
	})
	if err := dataSourceEndpointAgentRead(d, newEndpointAgentTestClient(transport, "12345")); err != nil {
		t.Fatalf("expected successful lookup on second page, got %v", err)
	}
	if requests != 2 {
		t.Fatalf("expected two API requests, got %d", requests)
	}
	if d.Id() != "endpoint-agent-456" || d.Get("agent_id").(string) != "endpoint-agent-456" {
		t.Fatalf("expected second-page agent ID in state, got ID %q and agent_id %q", d.Id(), d.Get("agent_id").(string))
	}
}

func TestDataSourceEndpointAgentRead_NotFoundWhenNoExactNameMatch(t *testing.T) {
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		return endpointAgentJSONResponse(http.StatusOK, `{"agents":[{"id":"partial-id","name":"Alice's laptop (backup)"}]}`), nil
	})
	d := schema.TestResourceDataRaw(t, dataSourceThousandeyesEndpointAgent().Schema, map[string]interface{}{
		"agent_name": "Alice's laptop",
	})
	err := dataSourceEndpointAgentRead(d, newEndpointAgentTestClient(transport, "12345"))
	if err == nil {
		t.Fatal("expected an error when no exact name match exists")
	}
	if !strings.Contains(err.Error(), "Alice's laptop") {
		t.Fatalf("expected no-match error to name the requested agent, got %v", err)
	}
}

func TestDataSourceEndpointAgentRead_PropagatesAPIError(t *testing.T) {
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		return endpointAgentJSONResponse(http.StatusInternalServerError, `{"error":"temporary failure"}`), nil
	})
	d := schema.TestResourceDataRaw(t, dataSourceThousandeyesEndpointAgent().Schema, map[string]interface{}{
		"agent_name": "Alice's laptop",
	})
	err := dataSourceEndpointAgentRead(d, newEndpointAgentTestClient(transport, "12345"))
	if err == nil {
		t.Fatal("expected API error to be returned")
	}
}

func TestDataSourceEndpointAgentRead_RefreshesDriftedID(t *testing.T) {
	const agentName = "Alice's laptop"
	responses := []string{
		`{"agents":[{"id":"old-agent-id","name":"Alice's laptop"}]}`,
		`{"agents":[{"id":"new-agent-id","name":"Alice's laptop"}]}`,
	}
	var reads int
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		if reads >= len(responses) {
			return nil, fmt.Errorf("unexpected extra API read %d", reads+1)
		}
		body := responses[reads]
		reads++
		return endpointAgentJSONResponse(http.StatusOK, body), nil
	})
	d := schema.TestResourceDataRaw(t, dataSourceThousandeyesEndpointAgent().Schema, map[string]interface{}{
		"agent_name": agentName,
	})
	apiClient := newEndpointAgentTestClient(transport, "12345")

	if err := dataSourceEndpointAgentRead(d, apiClient); err != nil {
		t.Fatalf("first read failed: %v", err)
	}
	if d.Id() != "old-agent-id" || d.Get("agent_id").(string) != "old-agent-id" {
		t.Fatalf("expected initial ID old-agent-id, got Terraform ID %q and agent_id %q", d.Id(), d.Get("agent_id").(string))
	}

	if err := dataSourceEndpointAgentRead(d, apiClient); err != nil {
		t.Fatalf("second read failed: %v", err)
	}
	if d.Id() != "new-agent-id" {
		t.Fatalf("expected refreshed Terraform ID new-agent-id, got %q", d.Id())
	}
	if got := d.Get("agent_id").(string); got != "new-agent-id" {
		t.Fatalf("expected refreshed agent_id new-agent-id, got %q", got)
	}
	if got := d.Get("agent_name").(string); got != agentName {
		t.Fatalf("expected configured name %q after refresh, got %q", agentName, got)
	}
}

func newEndpointAgentTestClient(transport http.RoundTripper, aid string) *client.APIClient {
	cfg := &client.Configuration{
		ServerURL: "https://api.example.test/v7",
		HTTPClient: &http.Client{
			Transport: transport,
		},
		Context: context.WithValue(context.Background(), accountGroupIdKey, aid),
	}
	httpClient := cfg.HTTPClient
	apiClient := client.NewAPIClient(cfg)
	// Bypass the SDK retry wrapper: these mocks test error handling, not retries.
	apiClient.GetConfig().HTTPClient = httpClient
	return apiClient
}

func endpointAgentJSONResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
}

type endpointAgentRoundTripFunc func(r *http.Request) (*http.Response, error)

func (f endpointAgentRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
