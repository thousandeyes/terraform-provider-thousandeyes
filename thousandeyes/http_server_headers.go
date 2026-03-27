package thousandeyes

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func syncHTTPServerHeaders(req *tests.HttpServerTestRequest) {
	if req == nil {
		return
	}

	merged := mergeHTTPHeaderStrings(req.Headers, req.CustomHeaders)
	if len(merged) == 0 {
		return
	}

	req.Headers = merged
	req.CustomHeaders = canonicalHTTPServerCustomHeaders(req.CustomHeaders, merged)
}

func syncHTTPServerResponseHeaders(resp *tests.HttpServerTestResponse) {
	if resp == nil {
		return
	}

	merged := mergeHTTPHeaderStrings(resp.Headers, resp.CustomHeaders)
	if len(merged) == 0 {
		return
	}

	resp.Headers = merged
	resp.CustomHeaders = canonicalHTTPServerCustomHeaders(resp.CustomHeaders, merged)
}

func mergeHTTPHeaderStrings(headers []string, customHeaders *tests.TestCustomHeaders) []string {
	seen := make(map[string]struct{})
	merged := make([]string, 0, len(headers))

	add := func(header string) {
		header = strings.TrimSpace(header)
		if header == "" {
			return
		}
		if _, ok := seen[header]; ok {
			return
		}
		seen[header] = struct{}{}
		merged = append(merged, header)
	}

	for _, header := range headers {
		add(header)
	}

	if customHeaders != nil && customHeaders.Root != nil {
		for name, value := range *customHeaders.Root {
			add(name + ": " + value)
		}
	}

	sort.Strings(merged)
	return merged
}

func splitHTTPHeader(header string) (string, string, bool) {
	parts := strings.SplitN(header, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	name := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if name == "" {
		return "", "", false
	}

	return name, value, true
}

func normalizeHTTPServerHeadersDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	headers, headersConfigured := rawConfigHeaderStrings(d)
	customHeaders, customHeadersConfigured := rawConfigCustomHeaders(d)

	merged := mergeHTTPHeaderStrings(headers, customHeaders)

	var headersValue interface{} = stringSliceToInterfaceSlice(merged)
	if !headersConfigured {
		before, _ := d.GetChange("headers")
		headersValue = before
	}
	if err := d.SetNew("headers", headersValue); err != nil {
		return err
	}

	var customHeadersValue interface{} = terraformHTTPServerCustomHeadersValue(customHeaders, merged)
	if !customHeadersConfigured {
		before, _ := d.GetChange("custom_headers")
		customHeadersValue = before
	}

	return d.SetNew("custom_headers", customHeadersValue)
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

func diffHeaderStrings(v interface{}) []string {
	switch headers := v.(type) {
	case *schema.Set:
		out := make([]string, 0, headers.Len())
		for _, item := range headers.List() {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, s)
			}
		}
		return out
	case []interface{}:
		out := make([]string, 0, len(headers))
		for _, item := range headers {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

func stringSliceToInterfaceSlice(v []string) []interface{} {
	out := make([]interface{}, 0, len(v))
	for _, item := range v {
		out = append(out, item)
	}
	return out
}

func terraformHTTPServerCustomHeadersValue(customHeaders *tests.TestCustomHeaders, mergedHeaders []string) []interface{} {
	root := make(map[string]interface{}, len(mergedHeaders))
	for _, header := range mergedHeaders {
		name, value, ok := splitHTTPHeader(header)
		if !ok {
			continue
		}
		root[name] = value
	}

	var all interface{} = map[string]interface{}{}
	var domains interface{} = map[string]interface{}{}
	if customHeaders != nil {
		all = normalizeStringMap(customHeaders.All)
		domains = normalizeNestedStringMap(customHeaders.Domains)
	}

	if len(root) == 0 && len(all.(map[string]interface{})) == 0 && len(domains.(map[string]interface{})) == 0 {
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

func canonicalHTTPServerCustomHeaders(customHeaders *tests.TestCustomHeaders, mergedHeaders []string) *tests.TestCustomHeaders {
	if len(mergedHeaders) == 0 {
		return customHeaders
	}

	if customHeaders == nil {
		customHeaders = &tests.TestCustomHeaders{}
	}

	root := make(map[string]string, len(mergedHeaders))
	for _, header := range mergedHeaders {
		name, value, ok := splitHTTPHeader(header)
		if !ok {
			continue
		}
		root[name] = value
	}

	if len(root) > 0 {
		customHeaders.Root = &root
	}
	return customHeaders
}

func interfaceMapToStringMap(v interface{}) map[string]string {
	rootMap, ok := v.(map[string]interface{})
	if !ok || rootMap == nil {
		return nil
	}

	out := make(map[string]string, len(rootMap))
	for k, raw := range rootMap {
		if s, ok := raw.(string); ok {
			out[k] = s
		}
	}
	return out
}

func interfaceNestedMapToStringMap(v interface{}) map[string]map[string]string {
	rawMap, ok := v.(map[string]interface{})
	if !ok || rawMap == nil {
		return nil
	}

	out := make(map[string]map[string]string, len(rawMap))
	for k, raw := range rawMap {
		if nested := interfaceMapToStringMap(raw); len(nested) > 0 {
			out[k] = nested
		}
	}
	return out
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
