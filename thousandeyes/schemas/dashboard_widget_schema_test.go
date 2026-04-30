package schemas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func TestDashboardWidgetSchemaColorGridColumnsIsComputed(t *testing.T) {
	colorGridSchema := DashboardWidgetSchema["color_grid_config"]
	require.NotNil(t, colorGridSchema)

	colorGridResource, ok := colorGridSchema.Elem.(*schema.Resource)
	require.True(t, ok)

	columnsSchema := colorGridResource.Schema["columns"]
	require.NotNil(t, columnsSchema)
	assert.True(t, columnsSchema.Optional)
	assert.True(t, columnsSchema.Computed)

	unitSchema := colorGridResource.Schema["unit"]
	require.NotNil(t, unitSchema)
	assert.True(t, unitSchema.Optional)
	assert.False(t, unitSchema.Computed)
}

func TestDashboardWidgetSchemaAlertListAlertTypesIsComputed(t *testing.T) {
	alertListSchema := DashboardWidgetSchema["alert_list_config"]
	require.NotNil(t, alertListSchema)

	alertListResource, ok := alertListSchema.Elem.(*schema.Resource)
	require.True(t, ok)

	alertTypesSchema := alertListResource.Schema["alert_types"]
	require.NotNil(t, alertTypesSchema)
	assert.True(t, alertTypesSchema.Optional)
	assert.True(t, alertTypesSchema.Computed)
}

func TestDashboardWidgetSchemaTestTableFilterKeyUsesTagIDOnly(t *testing.T) {
	testTableSchema := DashboardWidgetSchema["test_table_config"]
	require.NotNil(t, testTableSchema)

	testTableResource, ok := testTableSchema.Elem.(*schema.Resource)
	require.True(t, ok)

	filterSchema := testTableResource.Schema["filter"]
	filterResource, ok := filterSchema.Elem.(*schema.Resource)
	require.True(t, ok)

	filtersSchema := filterResource.Schema["filters"]
	filtersResource, ok := filtersSchema.Elem.(*schema.Resource)
	require.True(t, ok)

	keySchema := filtersResource.Schema["key"]
	require.NotNil(t, keySchema)

	_, errs := keySchema.ValidateFunc("Tag ID", "widgets.0.test_table_config.0.filter.0.filters.0.key")
	assert.Empty(t, errs)

	_, errs = keySchema.ValidateFunc("Label ID", "widgets.0.test_table_config.0.filter.0.filters.0.key")
	assert.NotEmpty(t, errs)

	// Keep a direct guard on the enum list in case the helper implementation changes.
	validate := validation.StringInSlice([]string{
		"Anything",
		"Test Name",
		"Target",
		"Test ID",
		"Test type",
		"Tag ID",
	}, false)
	_, errs = validate("Tag ID", "widgets.0.test_table_config.0.filter.0.filters.0.key")
	assert.Empty(t, errs)
	_, errs = validate("Label ID", "widgets.0.test_table_config.0.filter.0.filters.0.key")
	assert.NotEmpty(t, errs)
}

func TestDashboardWidgetSchemaRejectsDeprecatedCommonFilterProperties(t *testing.T) {
	propertySchema := getNestedSchemaProperty(t, DashboardWidgetSchema["filter"], "property")

	assertDashboardLabelValidation(t, propertySchema, "widgets.0.filter.0.property", "TEST", true)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.filter.0.property", "AGENT", true)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.filter.0.property", "ENDPOINT_LABEL", true)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.filter.0.property", "TEST_LABEL", false)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.filter.0.property", "AGENT_LABEL", false)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.filter.0.property", "ENDPOINT_TEST_LABEL", false)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.filter.0.property", "Test Labels", true)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.filter.0.property", "Agent Labels", true)
}

func TestDashboardWidgetSchemaRejectsDeprecatedNumberCardFilterProperties(t *testing.T) {
	propertySchema := getNestedSchemaProperty(t, &schema.Schema{
		Elem: &schema.Resource{Schema: NumberCardSchema},
	}, "filter", "property")

	assertDashboardLabelValidation(t, propertySchema, "widgets.0.number_card_config.0.filter.0.property", "TEST", true)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.number_card_config.0.filter.0.property", "AGENT", true)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.number_card_config.0.filter.0.property", "ENDPOINT_LABEL", true)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.number_card_config.0.filter.0.property", "TEST_LABEL", false)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.number_card_config.0.filter.0.property", "AGENT_LABEL", false)
	assertDashboardLabelValidation(t, propertySchema, "widgets.0.number_card_config.0.filter.0.property", "ENDPOINT_TEST_LABEL", false)
}

func TestDashboardWidgetSchemaRejectsDeprecatedGroupByLabelProperties(t *testing.T) {
	testCases := []struct {
		name      string
		root      *schema.Schema
		keys      []string
		fieldPath string
	}{
		{
			name:      "timeseries group_by",
			root:      DashboardWidgetSchema["timeseries_config"],
			keys:      []string{"group_by"},
			fieldPath: "widgets.0.timeseries_config.0.group_by",
		},
		{
			name:      "stacked area group_by",
			root:      DashboardWidgetSchema["stacked_area_config"],
			keys:      []string{"group_by"},
			fieldPath: "widgets.0.stacked_area_config.0.group_by",
		},
		{
			name:      "pie chart group_by",
			root:      DashboardWidgetSchema["pie_chart_config"],
			keys:      []string{"group_by"},
			fieldPath: "widgets.0.pie_chart_config.0.group_by",
		},
		{
			name:      "box and whiskers group_by",
			root:      DashboardWidgetSchema["box_and_whiskers_config"],
			keys:      []string{"group_by"},
			fieldPath: "widgets.0.box_and_whiskers_config.0.group_by",
		},
		{
			name:      "geo map group_by",
			root:      DashboardWidgetSchema["geo_map_config"],
			keys:      []string{"group_by"},
			fieldPath: "widgets.0.geo_map_config.0.group_by",
		},
		{
			name:      "table row_group_by",
			root:      DashboardWidgetSchema["table_config"],
			keys:      []string{"row_group_by"},
			fieldPath: "widgets.0.table_config.0.row_group_by",
		},
		{
			name:      "table column_group_by",
			root:      DashboardWidgetSchema["table_config"],
			keys:      []string{"column_group_by"},
			fieldPath: "widgets.0.table_config.0.column_group_by",
		},
		{
			name:      "stacked bar axis_group_by",
			root:      DashboardWidgetSchema["stacked_bar_chart_config"],
			keys:      []string{"axis_group_by"},
			fieldPath: "widgets.0.stacked_bar_chart_config.0.axis_group_by",
		},
		{
			name:      "grouped bar group_by",
			root:      DashboardWidgetSchema["grouped_bar_chart_config"],
			keys:      []string{"group_by"},
			fieldPath: "widgets.0.grouped_bar_chart_config.0.group_by",
		},
		{
			name:      "grouped bar axis_group_by",
			root:      DashboardWidgetSchema["grouped_bar_chart_config"],
			keys:      []string{"axis_group_by"},
			fieldPath: "widgets.0.grouped_bar_chart_config.0.axis_group_by",
		},
		{
			name:      "multi metric table row_group_by",
			root:      DashboardWidgetSchema["multi_metric_table_config"],
			keys:      []string{"row_group_by"},
			fieldPath: "widgets.0.multi_metric_table_config.0.row_group_by",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			propertySchema := getNestedSchemaProperty(t, testCase.root, testCase.keys...)
			assertDashboardLabelValidation(t, propertySchema, testCase.fieldPath, "TEST", true)
			assertDashboardLabelValidation(t, propertySchema, testCase.fieldPath, "AGENT", true)
			assertDashboardLabelValidation(t, propertySchema, testCase.fieldPath, "TEST_LABEL", false)
			assertDashboardLabelValidation(t, propertySchema, testCase.fieldPath, "AGENT_LABEL", false)
			assertDashboardLabelValidation(t, propertySchema, testCase.fieldPath, "ENDPOINT_TEST_LABEL", false)
		})
	}
}

func getNestedSchemaProperty(t *testing.T, root *schema.Schema, keys ...string) *schema.Schema {
	t.Helper()

	current := root
	for _, key := range keys {
		resource, ok := current.Elem.(*schema.Resource)
		require.True(t, ok)

		current = resource.Schema[key]
		require.NotNil(t, current)
	}

	return current
}

func assertDashboardLabelValidation(t *testing.T, propertySchema *schema.Schema, path string, value string, valid bool) {
	t.Helper()

	require.NotNil(t, propertySchema.ValidateFunc)

	_, errs := propertySchema.ValidateFunc(value, path)
	if valid {
		assert.Empty(t, errs)
		return
	}

	require.NotEmpty(t, errs)
	assert.ErrorContains(t, errs[0], "deprecated label value")
}
