package thousandeyes

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

func TestRawConfigCustomHeaders_preservesAllSections(t *testing.T) {
	d := mockRawConfigReader{
		values: map[string]cty.Value{
			"custom_headers": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"root": cty.MapVal(map[string]cty.Value{
						"X-TE-Test": cty.StringVal("upgrade-full"),
					}),
					"all": cty.MapVal(map[string]cty.Value{
						"Authorization": cty.StringVal("Bearer token"),
					}),
					"domains": cty.MapVal(map[string]cty.Value{
						"example.com": cty.MapVal(map[string]cty.Value{
							"X-Domain": cty.StringVal("value"),
						}),
					}),
				}),
			}),
		},
	}

	got, configured := rawConfigCustomHeaders(d)
	if !configured {
		t.Fatalf("expected custom_headers to be configured")
	}
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

type mockRawConfigReader struct {
	values map[string]cty.Value
}

func (m mockRawConfigReader) GetRawConfigAt(path cty.Path) (cty.Value, diag.Diagnostics) {
	if len(path) != 1 {
		return cty.NilVal, nil
	}

	step, ok := path[0].(cty.GetAttrStep)
	if !ok {
		return cty.NilVal, nil
	}

	if v, exists := m.values[step.Name]; exists {
		return v, nil
	}
	return cty.NilVal, nil
}
