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

	assertDashboardFilterPropertyValidation(t, propertySchema, "TEST", true)
	assertDashboardFilterPropertyValidation(t, propertySchema, "AGENT", true)
	assertDashboardFilterPropertyValidation(t, propertySchema, "ENDPOINT_LABEL", true)
	assertDashboardFilterPropertyValidation(t, propertySchema, "TEST_LABEL", false)
	assertDashboardFilterPropertyValidation(t, propertySchema, "AGENT_LABEL", false)
	assertDashboardFilterPropertyValidation(t, propertySchema, "ENDPOINT_TEST_LABEL", false)
	assertDashboardFilterPropertyValidation(t, propertySchema, "Test Labels", true)
	assertDashboardFilterPropertyValidation(t, propertySchema, "Agent Labels", true)
}

func TestDashboardWidgetSchemaRejectsDeprecatedNumberCardFilterProperties(t *testing.T) {
	propertySchema := getNestedSchemaProperty(t, &schema.Schema{
		Elem: &schema.Resource{Schema: NumberCardSchema},
	}, "filter", "property")

	assertDashboardFilterPropertyValidation(t, propertySchema, "TEST", true)
	assertDashboardFilterPropertyValidation(t, propertySchema, "AGENT", true)
	assertDashboardFilterPropertyValidation(t, propertySchema, "ENDPOINT_LABEL", true)
	assertDashboardFilterPropertyValidation(t, propertySchema, "TEST_LABEL", false)
	assertDashboardFilterPropertyValidation(t, propertySchema, "AGENT_LABEL", false)
	assertDashboardFilterPropertyValidation(t, propertySchema, "ENDPOINT_TEST_LABEL", false)
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

func assertDashboardFilterPropertyValidation(t *testing.T, propertySchema *schema.Schema, value string, valid bool) {
	t.Helper()

	require.NotNil(t, propertySchema.ValidateFunc)

	_, errs := propertySchema.ValidateFunc(value, "widgets.0.filter.0.property")
	if valid {
		assert.Empty(t, errs)
		return
	}

	require.NotEmpty(t, errs)
	assert.ErrorContains(t, errs[0], "deprecated label filter property")
}
