package schemas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardWidgetSchemaFixedTimespanIsComputed(t *testing.T) {
	fixedTimespanSchema := DashboardWidgetSchema["fixed_timespan"]
	require.NotNil(t, fixedTimespanSchema)
	assert.True(t, fixedTimespanSchema.Optional)
	assert.True(t, fixedTimespanSchema.Computed)

	fixedTimespanResource, ok := fixedTimespanSchema.Elem.(*schema.Resource)
	require.True(t, ok)
	assert.True(t, fixedTimespanResource.Schema["value"].Computed)
	assert.True(t, fixedTimespanResource.Schema["unit"].Computed)
}

func TestDashboardWidgetSchemaSuppressesTestTableAlertsDataSourceDrift(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DashboardSchema, map[string]interface{}{
		"title": "test",
		"widgets": []interface{}{
			map[string]interface{}{
				"type":        "Test Table",
				"data_source": "ALERTS",
			},
		},
	})

	suppress := DashboardWidgetSchema["data_source"].DiffSuppressFunc
	require.NotNil(t, suppress)
	assert.True(t, suppress("widgets.0.data_source", "", "ALERTS", d))
	assert.True(t, suppress("widgets.0.data_source", "ALERTS", "", d))
}

func TestDashboardWidgetSchemaDoesNotSuppressOtherDataSourceDiffs(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DashboardSchema, map[string]interface{}{
		"title": "test",
		"widgets": []interface{}{
			map[string]interface{}{
				"type":        "Table",
				"data_source": "ALERTS",
			},
		},
	})

	suppress := DashboardWidgetSchema["data_source"].DiffSuppressFunc
	require.NotNil(t, suppress)
	assert.False(t, suppress("widgets.0.data_source", "", "ALERTS", d))
	assert.False(t, suppress("widgets.0.data_source", "", "CLOUD_NATIVE_MONITORING", d))
}
