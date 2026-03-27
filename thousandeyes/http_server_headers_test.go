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
