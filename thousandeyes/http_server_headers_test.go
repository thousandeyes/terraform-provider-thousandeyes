package thousandeyes

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestTerraformHTTPServerCustomHeadersValue_preservesAllSections(t *testing.T) {
	root := map[string]string{"X-TE-Test": "upgrade-full"}
	all := map[string]string{"Authorization": "Bearer token"}
	domains := map[string]map[string]string{
		"example.com": {"X-Domain": "value"},
	}

	got := terraformHTTPServerCustomHeadersValue(&tests.TestCustomHeaders{
		Root:    &root,
		All:     &all,
		Domains: &domains,
	})

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

func TestTerraformHTTPServerCustomHeadersValue_nilReturnsEmpty(t *testing.T) {
	got := terraformHTTPServerCustomHeadersValue(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty slice for nil custom headers, got %#v", got)
	}
}

func TestTerraformHTTPServerCustomHeadersValue_allNilFieldsReturnsEmpty(t *testing.T) {
	got := terraformHTTPServerCustomHeadersValue(&tests.TestCustomHeaders{})
	if len(got) != 0 {
		t.Fatalf("expected empty slice for zero-value custom headers, got %#v", got)
	}
}

func TestTerraformHTTPServerOAuthValue_nilReturnsEmpty(t *testing.T) {
	got := terraformHTTPServerOAuthValue(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty slice for nil oauth, got %#v", got)
	}
}

func TestTerraformHTTPServerOAuthValue_zeroValueReturnsEmpty(t *testing.T) {
	got := terraformHTTPServerOAuthValue(tests.NewOAuth())
	if len(got) != 0 {
		t.Fatalf("expected empty slice for zero-value oauth, got %#v", got)
	}
}

func TestTerraformHTTPServerOAuthValue_preservesConfiguredFields(t *testing.T) {
	configID := "2660950"
	testURL := "https://auth.example.com/oauth/token"
	requestMethod := tests.REQUESTMETHOD_GET
	header := "Authorization: Basic test-client"
	username := "oauth-user"
	authType := tests.TESTAUTHTYPE_BASIC

	got := terraformHTTPServerOAuthValue(&tests.OAuth{
		ConfigId:      &configID,
		TestUrl:       &testURL,
		RequestMethod: &requestMethod,
		Headers:       &header,
		Username:      &username,
		AuthType:      &authType,
	})

	want := []interface{}{
		map[string]interface{}{
			"config_id":      configID,
			"test_url":       testURL,
			"request_method": string(requestMethod),
			"headers":        header,
			"username":       username,
			"auth_type":      string(authType),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected terraform oauth value: got %#v want %#v", got, want)
	}
}

func TestRawConfigCustomHeadersAllowsEmptyBlock(t *testing.T) {
	empty := &mockRawConfigReader{
		values: map[string]mockRawConfigValue{
			"custom_headers": {present: true, block: map[string]mockRawConfigValue{}},
		},
	}

	got, ok := rawConfigCustomHeaders(empty)
	if !ok {
		t.Fatalf("expected empty custom_headers block to be treated as configured")
	}
	if got != nil {
		t.Fatalf("expected nil custom headers for empty block, got %#v", got)
	}
}

func TestRawConfigCustomHeadersNotConfigured(t *testing.T) {
	reader := &mockRawConfigReader{values: map[string]mockRawConfigValue{}}

	_, ok := rawConfigCustomHeaders(reader)
	if ok {
		t.Fatalf("expected not-configured when custom_headers is absent from raw config")
	}
}

func TestRawConfigOAuthConfigured(t *testing.T) {
	reader := &mockRawConfigReader{
		values: map[string]mockRawConfigValue{
			"oauth": {present: true, block: map[string]mockRawConfigValue{}},
		},
	}

	if !rawConfigOAuthConfigured(reader) {
		t.Fatalf("expected oauth block to be treated as configured")
	}
}

func TestRawConfigOAuthNotConfigured(t *testing.T) {
	reader := &mockRawConfigReader{values: map[string]mockRawConfigValue{}}

	if rawConfigOAuthConfigured(reader) {
		t.Fatalf("expected oauth to be treated as absent when not configured")
	}
}

func TestRawConfigHeaderStringsSortsValues(t *testing.T) {
	reader := &mockRawConfigReader{
		values: map[string]mockRawConfigValue{
			"headers": {
				present: true,
				list: []string{
					"X-TE-Test: upgrade-full",
					"Content-Type: application/json",
				},
			},
		},
	}

	got, ok := rawConfigHeaderStrings(reader)
	if !ok {
		t.Fatalf("expected headers to be configured")
	}

	want := []string{
		"Content-Type: application/json",
		"X-TE-Test: upgrade-full",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected header order: got %#v want %#v", got, want)
	}
}

func TestRawConfigHeaderStringsTrimsWhitespace(t *testing.T) {
	reader := &mockRawConfigReader{
		values: map[string]mockRawConfigValue{
			"headers": {
				present: true,
				list:    []string{"  Content-Type: application/json  "},
			},
		},
	}

	got, ok := rawConfigHeaderStrings(reader)
	if !ok {
		t.Fatalf("expected headers to be configured")
	}
	if got[0] != "Content-Type: application/json" {
		t.Fatalf("expected trimmed header, got %q", got[0])
	}
}

func TestRawConfigHeaderStringsNotConfigured(t *testing.T) {
	reader := &mockRawConfigReader{values: map[string]mockRawConfigValue{}}

	_, ok := rawConfigHeaderStrings(reader)
	if ok {
		t.Fatalf("expected not-configured when headers is absent from raw config")
	}
}

func TestHttpHeaderSourceMode_headersConfigured(t *testing.T) {
	reader := &mockRawConfigReader{
		values: map[string]mockRawConfigValue{
			"headers": {present: true, list: []string{"X-Test: 1"}},
		},
	}
	if mode := httpHeaderSourceMode(reader); mode != httpHeaderSourceModeHeaders {
		t.Fatalf("expected %q, got %q", httpHeaderSourceModeHeaders, mode)
	}
}

func TestHttpHeaderSourceMode_customHeadersConfigured(t *testing.T) {
	reader := &mockRawConfigReader{
		values: map[string]mockRawConfigValue{
			"custom_headers": {present: true, block: map[string]mockRawConfigValue{}},
		},
	}
	if mode := httpHeaderSourceMode(reader); mode != httpHeaderSourceModeCustomHeaders {
		t.Fatalf("expected %q, got %q", httpHeaderSourceModeCustomHeaders, mode)
	}
}

func TestHttpHeaderSourceMode_neitherConfiguredDefaultsToHeaders(t *testing.T) {
	reader := &mockRawConfigReader{values: map[string]mockRawConfigValue{}}
	if mode := httpHeaderSourceMode(reader); mode != httpHeaderSourceModeHeaders {
		t.Fatalf("expected default %q, got %q", httpHeaderSourceModeHeaders, mode)
	}
}

func TestHttpHeaderSourceMode_customHeadersTakesPrecedence(t *testing.T) {
	reader := &mockRawConfigReader{
		values: map[string]mockRawConfigValue{
			"headers":        {present: true, list: []string{"X-Test: 1"}},
			"custom_headers": {present: true, block: map[string]mockRawConfigValue{}},
		},
	}
	if mode := httpHeaderSourceMode(reader); mode != httpHeaderSourceModeCustomHeaders {
		t.Fatalf("expected %q when both present, got %q", httpHeaderSourceModeCustomHeaders, mode)
	}
}

// --- mocks ---

type mockRawConfigReader struct {
	values map[string]mockRawConfigValue
}

type mockRawConfigValue struct {
	present bool
	list    []string
	block   map[string]mockRawConfigValue
}

func (m *mockRawConfigReader) GetRawConfigAt(path cty.Path) (cty.Value, diag.Diagnostics) {
	if len(path) != 1 {
		return cty.NullVal(cty.DynamicPseudoType), nil
	}

	attr, ok := path[0].(cty.GetAttrStep)
	if !ok {
		return cty.NullVal(cty.DynamicPseudoType), nil
	}

	raw, ok := m.values[attr.Name]
	if !ok || !raw.present {
		return cty.NullVal(cty.DynamicPseudoType), nil
	}

	return raw.toCty(), nil
}

func (m mockRawConfigValue) toCty() cty.Value {
	if m.block != nil {
		obj := make(map[string]cty.Value, len(m.block))
		for k, v := range m.block {
			obj[k] = v.toCty()
		}
		return cty.ListVal([]cty.Value{cty.ObjectVal(obj)})
	}

	if m.list != nil {
		values := make([]cty.Value, 0, len(m.list))
		for _, v := range m.list {
			values = append(values, cty.StringVal(v))
		}
		return cty.ListVal(values)
	}

	return cty.ObjectVal(map[string]cty.Value{})
}
