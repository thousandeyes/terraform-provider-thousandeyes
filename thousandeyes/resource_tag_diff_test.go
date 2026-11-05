package thousandeyes

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func TestResourceTagValidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(map[string]interface{})
		err    string
	}{
		{name: "dynamic tag"},
		{name: "custom key", change: func(c map[string]interface{}) {
			tagConfigFilter(c)["key"] = "department"
			tagConfigFilter(c)["scope"] = "custom"
		}},
		{name: "API validates keys", change: func(c map[string]interface{}) { tagConfigFilter(c)["key"] = "future-agent-field" }},
		{name: "static endpoint agent", change: func(c map[string]interface{}) { c["type"] = "static" }, err: "endpoint-agent tags require type = dynamic"},
		{name: "missing type", change: func(c map[string]interface{}) { delete(c, "type") }, err: "endpoint-agent tags require type = dynamic"},
		{name: "dynamic test", change: func(c map[string]interface{}) { c["object_type"] = "test" }, err: "dynamic tags require object_type = endpoint-agent"},
		{name: "missing object type", change: func(c map[string]interface{}) { delete(c, "object_type") }, err: "dynamic tags require object_type = endpoint-agent"},
		{name: "missing match type", change: func(c map[string]interface{}) { delete(c, "match_type") }, err: "dynamic tags require match_type"},
		{name: "missing filters", change: func(c map[string]interface{}) { delete(c, "filters") }, err: "dynamic tags require at least one filter"},
		{name: "static with match type", change: func(c map[string]interface{}) { c["object_type"] = "test"; c["type"] = "static"; delete(c, "filters") }, err: "only supported for dynamic endpoint-agent tags"},
		{name: "static with filters", change: func(c map[string]interface{}) {
			c["object_type"] = "test"
			c["type"] = "static"
			delete(c, "match_type")
		}, err: "only supported for dynamic endpoint-agent tags"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := tagDiffConfig()
			if tc.change != nil {
				tc.change(config)
			}
			_, err := tagConfigDiff(config, nil)
			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
	for _, objectType := range []string{"test", "dashboard", "endpoint-test", "v-agent"} {
		t.Run("existing "+objectType, func(t *testing.T) {
			config := map[string]interface{}{"object_type": objectType, "key": "office", "value": "london", "access_type": "all"}
			require.False(t, resourceTag().Validate(terraform.NewResourceConfigRaw(config)).HasError())
			_, err := tagConfigDiff(config, nil)
			require.NoError(t, err)
		})
	}
}

func TestResourceTagSchemaValidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(map[string]interface{})
	}{
		{"invalid object type", func(c map[string]interface{}) { c["object_type"] = "invalid" }},
		{"invalid type", func(c map[string]interface{}) { c["type"] = "unknown" }},
		{"invalid match type", func(c map[string]interface{}) { c["match_type"] = "unknown" }},
		{"empty filters", func(c map[string]interface{}) { c["filters"] = []interface{}{} }},
		{"missing key", func(c map[string]interface{}) { delete(tagConfigFilter(c), "key") }},
		{"empty key", func(c map[string]interface{}) { tagConfigFilter(c)["key"] = "" }},
		{"missing values", func(c map[string]interface{}) { delete(tagConfigFilter(c), "values") }},
		{"empty values", func(c map[string]interface{}) { tagConfigFilter(c)["values"] = []interface{}{} }},
		{"empty value", func(c map[string]interface{}) { tagConfigFilter(c)["values"] = []interface{}{""} }},
		{"missing mode", func(c map[string]interface{}) { delete(tagConfigFilter(c), "mode") }},
		{"invalid mode", func(c map[string]interface{}) { tagConfigFilter(c)["mode"] = "unknown" }},
		{"invalid scope", func(c map[string]interface{}) { tagConfigFilter(c)["scope"] = "unknown" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := tagDiffConfig()
			tc.change(config)
			diags := resourceTag().Validate(terraform.NewResourceConfigRaw(config))
			require.True(t, diags.HasError(), "%v", diags)
		})
	}
}

func TestResourceTagUnknownConfiguration(t *testing.T) {
	for _, key := range []string{"object_type", "type", "match_type", "filters"} {
		t.Run(key, func(t *testing.T) {
			config := tagDiffConfig()
			config[key] = "74D93920-ED26-11E3-AC10-0800200C9A66" // Terraform's legacy unknown-value sentinel.
			cfg := terraform.NewResourceConfigRaw(config)
			raw := map[string]cty.Value{"type": cty.StringVal("dynamic"), key: cty.DynamicVal}
			state := &terraform.InstanceState{RawConfig: cty.ObjectVal(raw)}
			_, err := resourceTag().Diff(context.Background(), state, cfg, nil)
			require.NoError(t, err)
		})
	}
}

func TestResourceTagStaticTypeRemainsComputed(t *testing.T) {
	config := map[string]interface{}{"object_type": "test", "key": "office", "value": "london", "access_type": "all"}
	d := schema.TestResourceDataRaw(t, resourceTag().Schema, config)
	d.SetId("static-tag")
	require.NoError(t, d.Set("type", "static"))
	diff, err := tagConfigDiff(config, d.State())
	require.NoError(t, err)
	require.True(t, diff.Empty(), "existing tags with computed type must have no plan: %#v", diff)
}

func TestResourceTagRefreshHasStablePlan(t *testing.T) {
	config := tagDiffConfig()
	tagConfigFilter(config)["values"] = []interface{}{"Windows", "Linux"}
	config["filters"] = append(config["filters"].([]interface{}), map[string]interface{}{
		"key": "department", "values": []interface{}{"Finance"}, "mode": "not-in", "scope": "custom",
	})
	d := schema.TestResourceDataRaw(t, resourceTag().Schema, config)
	d.SetId("tag-123")
	apiClient := client.NewAPIClient(&client.Configuration{
		ServerURL: "https://api.thousandeyes.com/v7",
		Context:   context.WithValue(context.Background(), accountGroupIdKey, "12345"),
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, `{"id":"tag-123","key":"platform","value":"selected","accessType":"all","objectType":"endpoint-agent","type":"dynamic","matchType":"and","filters":[{"key":"department","values":["Finance"],"mode":"not-in","scope":"custom"},{"key":"platform","values":["Linux","Windows"],"mode":"in"}]}`), nil
		})},
	})
	require.NoError(t, resourceTagRead(d, apiClient))
	diff, err := tagConfigDiff(config, d.State())
	require.NoError(t, err)
	require.True(t, diff.Empty(), "refresh must not cause a plan: %#v", diff)

	config["match_type"] = "or"
	tagConfigFilter(config)["values"] = []interface{}{"Mac"}
	diff, err = tagConfigDiff(config, d.State())
	require.NoError(t, err)
	require.False(t, diff.RequiresNew(), "criteria updates must preserve the tag ID")
	require.NotNil(t, diff.Attributes["match_type"])
	foundFilters := false
	for name := range diff.Attributes {
		foundFilters = foundFilters || strings.HasPrefix(name, "filters.")
	}
	require.True(t, foundFilters)
}

func TestResourceTagInvalidApplyDoesNotCallAPI(t *testing.T) {
	for _, call := range []func(*schema.ResourceData, interface{}) error{resourceTagCreate, resourceTagUpdate} {
		d := schema.TestResourceDataRaw(t, resourceTag().Schema, map[string]interface{}{"object_type": "endpoint-agent", "type": "static"})
		require.ErrorContains(t, call(d, nil), "endpoint-agent tags require type = dynamic")
	}
}

func tagDiffConfig() map[string]interface{} {
	return map[string]interface{}{
		"key": "platform", "value": "selected", "access_type": "all",
		"object_type": "endpoint-agent", "type": "dynamic", "match_type": "and",
		"filters": []interface{}{map[string]interface{}{"key": "platform", "values": []interface{}{"Windows"}, "mode": "in"}},
	}
}

func tagConfigFilter(config map[string]interface{}) map[string]interface{} {
	return config["filters"].([]interface{})[0].(map[string]interface{})
}

func tagConfigDiff(config map[string]interface{}, state *terraform.InstanceState) (*terraform.InstanceDiff, error) {
	if state == nil {
		state = &terraform.InstanceState{}
	}
	// Mirror the raw configuration supplied by Terraform over the provider protocol.
	typeValue := cty.NullVal(cty.String)
	if v, ok := config["type"].(string); ok {
		typeValue = cty.StringVal(v)
	}
	state.RawConfig = cty.ObjectVal(map[string]cty.Value{"type": typeValue})
	return resourceTag().Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
}
