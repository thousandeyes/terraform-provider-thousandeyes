package thousandeyes

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func TestDataSourceDashboardFilterRead_Success(t *testing.T) {
	const (
		expectedName       = "Target Filter"
		expectedFilterID   = "filter-123"
		expectedAid        = "12345"
		expectedRequestURL = "/v7/dashboards/filters"
	)

	var gotSearchPattern string
	var gotAid string
	var gotPath string

	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.Path
		gotSearchPattern = r.URL.Query().Get("searchPattern")
		gotAid = r.URL.Query().Get("aid")
		return jsonResponse(http.StatusOK, `{
			"dashboardFilters": [
				{"context":[],"aid":"12345","id":"filter-999","name":"Other Filter"},
				{"context":[],"aid":"12345","id":"filter-123","name":"Target Filter"}
			]
		}`), nil
	})

	d := schema.TestResourceDataRaw(t, dataSourceThousandeyesDashboardFilter().Schema, map[string]interface{}{
		"name": expectedName,
	})

	apiClient := newDashboardFilterTestClient(transport, expectedAid)
	err := dataSourceDashboardFilterRead(d, apiClient)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gotPath != expectedRequestURL {
		t.Fatalf("expected request path %q, got %q", expectedRequestURL, gotPath)
	}
	if gotSearchPattern != expectedName {
		t.Fatalf("expected searchPattern query %q, got %q", expectedName, gotSearchPattern)
	}
	if gotAid != expectedAid {
		t.Fatalf("expected aid query %q, got %q", expectedAid, gotAid)
	}

	if d.Id() != expectedFilterID {
		t.Fatalf("expected state id %q, got %q", expectedFilterID, d.Id())
	}
	if gotName := d.Get("name").(string); gotName != expectedName {
		t.Fatalf("expected name %q, got %q", expectedName, gotName)
	}
	if gotFilterID := d.Get("filter_id").(string); gotFilterID != expectedFilterID {
		t.Fatalf("expected filter_id %q, got %q", expectedFilterID, gotFilterID)
	}
}

func TestDataSourceDashboardFilterRead_NotFound(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{
			"dashboardFilters": [
				{"context":[],"aid":"12345","id":"filter-111","name":"Some Other Filter"}
			]
		}`), nil
	})

	d := schema.TestResourceDataRaw(t, dataSourceThousandeyesDashboardFilter().Schema, map[string]interface{}{
		"name": "Missing Filter",
	})

	apiClient := newDashboardFilterTestClient(transport, "12345")
	err := dataSourceDashboardFilterRead(d, apiClient)
	if err == nil {
		t.Fatal("expected error when filter is not found, got nil")
	}
	if !strings.Contains(err.Error(), "unable to locate any filter with the name: Missing Filter") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func newDashboardFilterTestClient(transport http.RoundTripper, aid string) *client.APIClient {
	cfg := &client.Configuration{
		ServerURL: "https://api.thousandeyes.com/v7",
		HTTPClient: &http.Client{
			Transport: transport,
		},
		Context: context.WithValue(context.Background(), accountGroupIdKey, aid),
	}
	return client.NewAPIClient(cfg)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
}

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
