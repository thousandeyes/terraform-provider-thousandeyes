package thousandeyes

import (
	"context"
	"sort"
	"strings"

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

	customHeaders := req.CustomHeaders
	if customHeaders == nil {
		customHeaders = &tests.TestCustomHeaders{}
	}

	root := make(map[string]string, len(merged))
	for _, header := range merged {
		name, value, ok := splitHTTPHeader(header)
		if !ok {
			continue
		}
		root[name] = value
	}
	if len(root) > 0 {
		customHeaders.Root = &root
		req.CustomHeaders = customHeaders
	}
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

	customHeaders := resp.CustomHeaders
	if customHeaders == nil {
		customHeaders = &tests.TestCustomHeaders{}
	}

	root := make(map[string]string, len(merged))
	for _, header := range merged {
		name, value, ok := splitHTTPHeader(header)
		if !ok {
			continue
		}
		root[name] = value
	}
	if len(root) > 0 {
		customHeaders.Root = &root
		resp.CustomHeaders = customHeaders
	}
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
	headers := diffHeaderStrings(d.Get("headers"))
	customRoot := diffCustomRootHeaders(d.Get("custom_headers"))

	customHeaders := &tests.TestCustomHeaders{}
	if len(customRoot) > 0 {
		customHeaders.Root = &customRoot
	}

	merged := mergeHTTPHeaderStrings(headers, customHeaders)
	root := make(map[string]interface{}, len(merged))
	for _, header := range merged {
		name, value, ok := splitHTTPHeader(header)
		if !ok {
			continue
		}
		root[name] = value
	}

	if err := d.SetNew("headers", stringSliceToInterfaceSlice(merged)); err != nil {
		return err
	}

	return d.SetNew("custom_headers", terraformHTTPServerCustomHeadersValue(&tests.TestCustomHeaders{
		Root: &customRoot,
	}, merged))
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

func diffCustomRootHeaders(v interface{}) map[string]string {
	var items []interface{}
	switch raw := v.(type) {
	case *schema.Set:
		items = raw.List()
	case []interface{}:
		items = raw
	default:
		return nil
	}

	if len(items) == 0 {
		return nil
	}

	first, ok := items[0].(map[string]interface{})
	if !ok {
		return nil
	}

	rootVal, ok := first["root"]
	if !ok || rootVal == nil {
		return nil
	}

	rootMap, ok := rootVal.(map[string]interface{})
	if !ok {
		return nil
	}

	out := make(map[string]string, len(rootMap))
	for k, v := range rootMap {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
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
