package thousandeyes

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

// testFilterSet builds a *schema.Set of filter blocks for use in builder tests.
// Accepts plain maps with []interface{} values — auto-wraps them into *schema.Set.
func testFilterSet(filters ...map[string]interface{}) *schema.Set {
	filterResource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"property": {Type: schema.TypeString, Required: true},
			"values":   {Type: schema.TypeSet, Required: true, Elem: &schema.Schema{Type: schema.TypeString}},
		},
	}
	items := make([]interface{}, len(filters))
	for i, f := range filters {
		converted := make(map[string]interface{}, len(f))
		for k, v := range f {
			if k == "values" {
				if vals, ok := v.([]interface{}); ok {
					converted[k] = schema.NewSet(schema.HashString, vals)
				} else {
					converted[k] = v
				}
			} else {
				converted[k] = v
			}
		}
		items[i] = converted
	}
	return schema.NewSet(schema.HashResource(filterResource), items)
}

func TestSetCommonBuilderFields_IgnoresIsEmbedded(t *testing.T) {
	widget := dashboards.NewApiGeoMapWidget("Map")

	setCommonBuilderFields(widget, map[string]interface{}{
		"is_embedded": true,
	})

	isEmbedded, ok := widget.GetIsEmbeddedOk()
	assert.False(t, ok)
	assert.Nil(t, isEmbedded)
}

func TestSetCommonBuilderFields_PreservesFalseBooleans(t *testing.T) {
	widget := dashboards.NewApiGeoMapWidget("Map")

	setCommonBuilderFields(widget, map[string]interface{}{
		"should_exclude_alert_suppression_windows": false,
	})

	shouldExclude, ok := widget.GetShouldExcludeAlertSuppressionWindowsOk()
	assert.True(t, ok)
	if assert.NotNil(t, shouldExclude) {
		assert.False(t, *shouldExclude)
	}
}

func TestSetCommonBuilderFields_PreservesTrueBooleans(t *testing.T) {
	widget := dashboards.NewApiGeoMapWidget("Map")

	setCommonBuilderFields(widget, map[string]interface{}{
		"should_exclude_alert_suppression_windows": true,
	})

	shouldExclude, ok := widget.GetShouldExcludeAlertSuppressionWindowsOk()
	assert.True(t, ok)
	if assert.NotNil(t, shouldExclude) {
		assert.True(t, *shouldExclude)
	}
}
