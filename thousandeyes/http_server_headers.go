package thousandeyes

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func normalizeHTTPServerHeadersDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	headers, headersConfigured := rawConfigHeaderStrings(d)
	customHeaders, customHeadersConfigured := rawConfigCustomHeaders(d)

	if headersConfigured && customHeadersConfigured {
		return fmt.Errorf("`headers` and `custom_headers` are mutually exclusive; configure only one")
	}

	if headersConfigured {
		if err := d.SetNew(httpHeaderSourceModeField, httpHeaderSourceModeHeaders); err != nil {
			return err
		}
		if err := d.SetNew("headers", stringSliceToInterfaceSlice(headers)); err != nil {
			return err
		}
		before, _ := d.GetChange("custom_headers")
		return d.SetNew("custom_headers", before)
	}

	if customHeadersConfigured {
		if err := d.SetNew(httpHeaderSourceModeField, httpHeaderSourceModeCustomHeaders); err != nil {
			return err
		}
		before, _ := d.GetChange("headers")
		if err := d.SetNew("headers", before); err != nil {
			return err
		}
		return d.SetNew("custom_headers", terraformHTTPServerCustomHeadersValue(customHeaders))
	}

	beforeMode, _ := d.GetChange(httpHeaderSourceModeField)
	if err := d.SetNew(httpHeaderSourceModeField, beforeMode); err != nil {
		return err
	}
	beforeHeaders, _ := d.GetChange("headers")
	if err := d.SetNew("headers", beforeHeaders); err != nil {
		return err
	}
	beforeCustomHeaders, _ := d.GetChange("custom_headers")
	return d.SetNew("custom_headers", beforeCustomHeaders)
}

func rawConfigHeaderStrings(d rawConfigReader) ([]string, bool) {
	raw, diags := d.GetRawConfigAt(cty.Path{cty.GetAttrStep{Name: "headers"}})
	if diags.HasError() || !raw.IsKnown() || raw.IsNull() {
		return nil, false
	}

	var out []string
	it := raw.ElementIterator()
	for it.Next() {
		_, v := it.Element()
		if !v.IsKnown() || v.IsNull() {
			continue
		}
		out = append(out, strings.TrimSpace(v.AsString()))
	}
	sort.Strings(out)
	return out, true
}

func rawConfigCustomHeaders(d rawConfigReader) (*tests.TestCustomHeaders, bool) {
	raw, diags := d.GetRawConfigAt(cty.Path{cty.GetAttrStep{Name: "custom_headers"}})
	if diags.HasError() || !raw.IsKnown() || raw.IsNull() || raw.LengthInt() == 0 {
		return nil, false
	}

	it := raw.ElementIterator()
	if !it.Next() {
		return nil, false
	}
	_, first := it.Element()
	if !first.IsKnown() || first.IsNull() {
		return nil, false
	}

	customHeaders := &tests.TestCustomHeaders{}
	if root := ctyObjectToStringMap(first, "root"); len(root) > 0 {
		customHeaders.Root = &root
	}
	if all := ctyObjectToStringMap(first, "all"); len(all) > 0 {
		customHeaders.All = &all
	}
	if domains := ctyObjectToNestedStringMap(first, "domains"); len(domains) > 0 {
		customHeaders.Domains = &domains
	}

	if customHeaders.Root == nil && customHeaders.All == nil && customHeaders.Domains == nil {
		return nil, true
	}
	return customHeaders, true
}

func stringSliceToInterfaceSlice(v []string) []interface{} {
	out := make([]interface{}, 0, len(v))
	for _, item := range v {
		out = append(out, item)
	}
	return out
}

func terraformHTTPServerCustomHeadersValue(customHeaders *tests.TestCustomHeaders) []interface{} {
	if customHeaders == nil {
		return []interface{}{}
	}

	root := normalizeStringMap(customHeaders.Root)
	all := normalizeStringMap(customHeaders.All)
	domains := normalizeNestedStringMap(customHeaders.Domains)

	if len(root) == 0 && len(all) == 0 && len(domains) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"root":    root,
			"domains": domains,
			"all":     all,
		},
	}
}

func ctyObjectToStringMap(v cty.Value, attr string) map[string]string {
	if !v.Type().HasAttribute(attr) {
		return nil
	}
	field := v.GetAttr(attr)
	if !field.IsKnown() || field.IsNull() {
		return nil
	}

	out := make(map[string]string, field.LengthInt())
	it := field.ElementIterator()
	for it.Next() {
		k, val := it.Element()
		if !val.IsKnown() || val.IsNull() {
			continue
		}
		out[k.AsString()] = val.AsString()
	}
	return out
}

func ctyObjectToNestedStringMap(v cty.Value, attr string) map[string]map[string]string {
	if !v.Type().HasAttribute(attr) {
		return nil
	}
	field := v.GetAttr(attr)
	if !field.IsKnown() || field.IsNull() {
		return nil
	}

	out := make(map[string]map[string]string, field.LengthInt())
	it := field.ElementIterator()
	for it.Next() {
		k, val := it.Element()
		if nested := ctyMapToStringMap(val); len(nested) > 0 {
			out[k.AsString()] = nested
		}
	}
	return out
}

func ctyMapToStringMap(v cty.Value) map[string]string {
	if !v.IsKnown() || v.IsNull() {
		return nil
	}

	out := make(map[string]string, v.LengthInt())
	it := v.ElementIterator()
	for it.Next() {
		k, val := it.Element()
		if !val.IsKnown() || val.IsNull() {
			continue
		}
		out[k.AsString()] = val.AsString()
	}
	return out
}

type rawConfigReader interface {
	GetRawConfigAt(cty.Path) (cty.Value, diag.Diagnostics)
}

func httpHeaderSourceMode(d rawConfigReader) string {
	if _, ok := rawConfigCustomHeaders(d); ok {
		return httpHeaderSourceModeCustomHeaders
	}
	if _, ok := rawConfigHeaderStrings(d); ok {
		return httpHeaderSourceModeHeaders
	}

	if resourceData, ok := d.(*schema.ResourceData); ok {
		if mode, ok := resourceData.Get(httpHeaderSourceModeField).(string); ok && mode != "" {
			return mode
		}
	}

	return httpHeaderSourceModeHeaders
}
