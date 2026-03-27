package thousandeyes

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestSyncHTTPServerHeaders_mergesHeadersAndCustomRoot(t *testing.T) {
	root := map[string]string{
		"X-TE-Test": "upgrade-full",
	}

	req := &tests.HttpServerTestRequest{
		Headers: []string{
			"Accept: application/json",
			"Content-Type: application/json",
		},
		CustomHeaders: &tests.TestCustomHeaders{
			Root: &root,
		},
	}

	syncHTTPServerHeaders(req)

	wantHeaders := []string{
		"Accept: application/json",
		"Content-Type: application/json",
		"X-TE-Test: upgrade-full",
	}

	if !reflect.DeepEqual(req.Headers, wantHeaders) {
		t.Fatalf("unexpected headers: got %#v want %#v", req.Headers, wantHeaders)
	}

	wantRoot := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"X-TE-Test":    "upgrade-full",
	}

	if req.CustomHeaders == nil || req.CustomHeaders.Root == nil {
		t.Fatalf("expected custom_headers.root to be populated")
	}

	if !reflect.DeepEqual(*req.CustomHeaders.Root, wantRoot) {
		t.Fatalf("unexpected custom_headers.root: got %#v want %#v", *req.CustomHeaders.Root, wantRoot)
	}
}

func TestDiffCustomHeaders_preservesAllSections(t *testing.T) {
	v := []interface{}{
		map[string]interface{}{
			"root": map[string]interface{}{
				"X-TE-Test": "upgrade-full",
			},
			"all": map[string]interface{}{
				"Authorization": "Bearer token",
			},
			"domains": map[string]interface{}{
				"example.com": map[string]interface{}{
					"X-Domain": "value",
				},
			},
		},
	}

	got := diffCustomHeaders(v)
	if got == nil || got.Root == nil || got.All == nil || got.Domains == nil {
		t.Fatalf("expected all custom header sections to be preserved, got %#v", got)
	}

	if (*got.Root)["X-TE-Test"] != "upgrade-full" {
		t.Fatalf("unexpected root headers: %#v", *got.Root)
	}
	if (*got.All)["Authorization"] != "Bearer token" {
		t.Fatalf("unexpected all headers: %#v", *got.All)
	}
	if (*got.Domains)["example.com"]["X-Domain"] != "value" {
		t.Fatalf("unexpected domain headers: %#v", *got.Domains)
	}
}

func TestTerraformHTTPServerCustomHeadersValue_preservesAllSections(t *testing.T) {
	all := map[string]string{"Authorization": "Bearer token"}
	domains := map[string]map[string]string{
		"example.com": {"X-Domain": "value"},
	}

	got := terraformHTTPServerCustomHeadersValue(&tests.TestCustomHeaders{
		All:     &all,
		Domains: &domains,
	}, []string{"X-TE-Test: upgrade-full"})

	want := []interface{}{
		map[string]interface{}{
			"root": map[string]interface{}{
				"X-TE-Test": "upgrade-full",
			},
			"all": map[string]interface{}{
				"Authorization": "Bearer token",
			},
			"domains": map[string]interface{}{
				"example.com": map[string]interface{}{
					"X-Domain": "value",
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected terraform custom_headers value: got %#v want %#v", got, want)
	}
}

func TestDiffHeaderStringsFromSet(t *testing.T) {
	set := schema.NewSet(schema.HashString, []interface{}{
		"X-TE-Test: upgrade-full",
		"Content-Type: application/json",
	})

	got := diffHeaderStrings(set)
	want := []string{
		"X-TE-Test: upgrade-full",
		"Content-Type: application/json",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected diff headers: got %#v want %#v", got, want)
	}
}
