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
