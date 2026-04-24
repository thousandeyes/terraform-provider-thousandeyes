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
