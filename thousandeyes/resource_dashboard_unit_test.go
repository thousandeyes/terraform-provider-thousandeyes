package thousandeyes

import (
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildDashboardStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, d *dashboards.Dashboard)
	}{
		{
			name: "basic fields",
			input: map[string]interface{}{
				"title":       "Test Dashboard",
				"description": "Test Description",
				"is_private":  true,
			},
			validate: func(t *testing.T, d *dashboards.Dashboard) {
				assert.Equal(t, "Test Dashboard", d.GetTitle())
				assert.Equal(t, "Test Description", d.GetDescription())
				assert.True(t, d.GetIsPrivate())
			},
		},
		{
			name: "with global_filter_id and is_global_override",
			input: map[string]interface{}{
				"title":              "Test Dashboard",
				"description":        "Test Description",
				"is_private":         false,
				"global_filter_id":   "filter-123",
				"is_global_override": true,
			},
			validate: func(t *testing.T, d *dashboards.Dashboard) {
				assert.Equal(t, "Test Dashboard", d.GetTitle())
				assert.Equal(t, "filter-123", d.GetGlobalFilterId())
				assert.True(t, d.GetIsGlobalOverride())
			},
		},
		{
			name: "with duration timespan",
			input: map[string]interface{}{
				"title":       "Test Dashboard",
				"description": "Test Description",
				"is_private":  false,
				"default_timespan": []interface{}{
					map[string]interface{}{
						"duration": 3600,
					},
				},
			},
			validate: func(t *testing.T, d *dashboards.Dashboard) {
				ts := d.GetDefaultTimespan()
				assert.Equal(t, int64(3600), ts.GetDuration())
			},
		},
		{
			name: "with time range timespan",
			input: map[string]interface{}{
				"title":       "Test Dashboard",
				"description": "Test Description",
				"is_private":  false,
				"default_timespan": []interface{}{
					map[string]interface{}{
						"start": "2026-01-01T00:00:00Z",
						"end":   "2026-02-01T00:00:00Z",
					},
				},
			},
			validate: func(t *testing.T, d *dashboards.Dashboard) {
				ts := d.GetDefaultTimespan()
				expectedStart, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
				expectedEnd, _ := time.Parse(time.RFC3339, "2026-02-01T00:00:00Z")
				assert.Equal(t, expectedStart, ts.GetStart())
				assert.Equal(t, expectedEnd, ts.GetEnd())
			},
		},
		{
			name: "with invalid start time logs warning and skips",
			input: map[string]interface{}{
				"title":       "Test Dashboard",
				"description": "Test Description",
				"is_private":  false,
				"default_timespan": []interface{}{
					map[string]interface{}{
						"start": "invalid-time-format",
						"end":   "2026-02-01T00:00:00Z",
					},
				},
			},
			validate: func(t *testing.T, d *dashboards.Dashboard) {
				ts := d.GetDefaultTimespan()
				// Invalid start should result in zero time
				assert.True(t, ts.GetStart().IsZero())
				// Valid end should still be set
				expectedEnd, _ := time.Parse(time.RFC3339, "2026-02-01T00:00:00Z")
				assert.Equal(t, expectedEnd, ts.GetEnd())
			},
		},
		{
			name: "with invalid end time logs warning and skips",
			input: map[string]interface{}{
				"title":       "Test Dashboard",
				"description": "Test Description",
				"is_private":  false,
				"default_timespan": []interface{}{
					map[string]interface{}{
						"start": "2026-01-01T00:00:00Z",
						"end":   "invalid-time-format",
					},
				},
			},
			validate: func(t *testing.T, d *dashboards.Dashboard) {
				ts := d.GetDefaultTimespan()
				// Valid start should be set
				expectedStart, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
				assert.Equal(t, expectedStart, ts.GetStart())
				// Invalid end should result in zero time
				assert.True(t, ts.GetEnd().IsZero())
			},
		},
		{
			name: "without default_timespan",
			input: map[string]interface{}{
				"title":       "Test Dashboard",
				"description": "Test Description",
				"is_private":  false,
			},
			validate: func(t *testing.T, d *dashboards.Dashboard) {
				// No timespan should be set
				ts := d.GetDefaultTimespan()
				assert.Equal(t, int64(0), ts.GetDuration())
				assert.True(t, ts.GetStart().IsZero())
				assert.True(t, ts.GetEnd().IsZero())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, schemas.DashboardSchema, tc.input)
			result, err := buildDashboardStruct(d)
			require.NoError(t, err)
			tc.validate(t, result)
		})
	}
}

func TestBuildDashboardStruct_explicitEmptyWidgetsSendsEmptySlice(t *testing.T) {
	d := schema.TestResourceDataRaw(t, schemas.DashboardSchema, map[string]interface{}{
		"title":   "T",
		"widgets": []interface{}{},
	})
	result, err := buildDashboardStruct(d)
	require.NoError(t, err)
	assert.Empty(t, result.GetWidgets())
}

func TestBuildDashboardStruct_noWidgetsSendsEmptySlice(t *testing.T) {
	d := schema.TestResourceDataRaw(t, schemas.DashboardSchema, map[string]interface{}{
		"title": "T",
	})
	result, err := buildDashboardStruct(d)
	require.NoError(t, err)
	assert.Empty(t, result.GetWidgets())
}

func TestBuildDashboardStruct_withWidgets(t *testing.T) {
	d := schema.TestResourceDataRaw(t, schemas.DashboardSchema, map[string]interface{}{
		"title": "T",
		"widgets": []interface{}{
			map[string]interface{}{
				"type":  "List",
				"title": "My List",
			},
		},
	})
	result, err := buildDashboardStruct(d)
	require.NoError(t, err)
	require.Len(t, result.GetWidgets(), 1)
	assert.NotNil(t, result.GetWidgets()[0].ApiListWidget)
}

// TestConfigWidgetsNullOrEmpty tests configWidgetsNullOrEmpty which covers the case where
// Terraform represents absent Optional+Computed block lists as null in raw config.
func TestConfigWidgetsNullOrEmpty(t *testing.T) {
	widgetsListType := cty.List(cty.DynamicPseudoType)

	tests := []struct {
		name string
		cfg  cty.Value
		want bool
	}{
		{
			name: "nil value",
			cfg:  cty.NilVal,
			want: false,
		},
		{
			name: "null value",
			cfg:  cty.NullVal(cty.EmptyObject),
			want: false,
		},
		{
			name: "unknown value",
			cfg:  cty.DynamicVal,
			want: false,
		},
		{
			name: "object without widgets attribute",
			cfg:  cty.EmptyObjectVal,
			want: false,
		},
		{
			name: "widgets is null — the key case: Optional+Computed absent blocks",
			cfg: cty.ObjectVal(map[string]cty.Value{
				"widgets": cty.NullVal(widgetsListType),
			}),
			want: true, // null = user omitted all blocks
		},
		{
			name: "widgets is explicitly empty list",
			cfg: cty.ObjectVal(map[string]cty.Value{
				"widgets": cty.ListValEmpty(cty.DynamicPseudoType),
			}),
			want: true, // empty list = user wrote zero blocks
		},
		{
			name: "widgets has one element",
			cfg: cty.ObjectVal(map[string]cty.Value{
				"widgets": cty.ListVal([]cty.Value{cty.StringVal("a")}),
			}),
			want: false,
		},
		{
			name: "widgets attribute is a string (wrong type)",
			cfg: cty.ObjectVal(map[string]cty.Value{
				"widgets": cty.StringVal("not-a-list"),
			}),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, configWidgetsNullOrEmpty(tc.cfg))
		})
	}
}

// TestCustomizeDiff_explicitEmptyWidgets_helperIntegration verifies the two halves of the
// "remove all widgets" fix work together:
//  1. configWidgetsNullOrEmpty detects both null (absent blocks) and empty list in raw config.
//  2. buildDashboardStruct produces a payload with an explicit empty widget slice.
func TestCustomizeDiff_explicitEmptyWidgets_helperIntegration(t *testing.T) {
	widgetsListType := cty.List(cty.DynamicPseudoType)

	// Case A: Terraform uses null for absent Optional+Computed block attributes — the
	// most common case when a user removes all widget blocks from their HCL.
	nullCfg := cty.ObjectVal(map[string]cty.Value{
		"widgets": cty.NullVal(widgetsListType),
	})
	assert.True(t, configWidgetsNullOrEmpty(nullCfg),
		"null widgets (absent blocks) should trigger SetNew in CustomizeDiff")

	// Case B: Terraform produces an empty tuple when it can confirm zero blocks.
	emptyCfg := cty.ObjectVal(map[string]cty.Value{
		"widgets": cty.ListValEmpty(cty.DynamicPseudoType),
	})
	assert.True(t, configWidgetsNullOrEmpty(emptyCfg),
		"empty list widgets should trigger SetNew in CustomizeDiff")

	// Part 2: by the time Update runs, d.Get("widgets") returns [] because CustomizeDiff
	// already called SetNew. Confirm buildDashboardStruct sends an empty slice to the API.
	d := schema.TestResourceDataRaw(t, schemas.DashboardSchema, map[string]interface{}{
		"title":   "T",
		"widgets": []interface{}{},
	})
	result, err := buildDashboardStruct(d)
	require.NoError(t, err)
	assert.NotNil(t, result.Widgets, "widgets field should be set (not omitted) on the payload")
	assert.Empty(t, result.GetWidgets(), "widgets slice should be empty so the API clears all widgets")
}

func TestResourceDataApiDashboardMapper(t *testing.T) {
	tests := []struct {
		name     string
		input    dashboards.ApiDashboard
		validate func(t *testing.T, d *schema.ResourceData)
	}{
		{
			name: "basic fields",
			input: func() dashboards.ApiDashboard {
				d := dashboards.NewApiDashboard()
				d.SetAid("123456")
				d.SetTitle("Test Dashboard")
				d.SetDescription("Test Description")
				d.SetIsPrivate(true)
				d.SetDashboardCreatedBy("user1")
				d.SetDashboardModifiedBy("user2")
				return *d
			}(),
			validate: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "123456", d.Get("aid"))
				assert.Equal(t, "Test Dashboard", d.Get("title"))
				assert.Equal(t, "Test Description", d.Get("description"))
				assert.True(t, d.Get("is_private").(bool))
				assert.Equal(t, "user1", d.Get("dashboard_created_by"))
				assert.Equal(t, "user2", d.Get("dashboard_modified_by"))
			},
		},
		{
			name: "with global_filter_id and is_global_override",
			input: func() dashboards.ApiDashboard {
				d := dashboards.NewApiDashboard()
				d.SetTitle("Test Dashboard")
				d.SetDescription("Test Description")
				d.SetGlobalFilterId("filter-123")
				d.SetIsGlobalOverride(true)
				d.SetIsMigratedReport(false)
				return *d
			}(),
			validate: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "filter-123", d.Get("global_filter_id"))
				assert.True(t, d.Get("is_global_override").(bool))
				assert.False(t, d.Get("is_migrated_report").(bool))
			},
		},
		{
			name: "with duration timespan",
			input: func() dashboards.ApiDashboard {
				d := dashboards.NewApiDashboard()
				d.SetTitle("Test Dashboard")
				d.SetDescription("Test Description")
				ts := dashboards.NewDefaultTimespan()
				ts.SetDuration(3600)
				d.SetDefaultTimespan(*ts)
				return *d
			}(),
			validate: func(t *testing.T, d *schema.ResourceData) {
				timespan := d.Get("default_timespan").([]interface{})
				assert.Len(t, timespan, 1)
				ts := timespan[0].(map[string]interface{})
				// Schema stores as int, SDK returns int64, so compare as int
				assert.Equal(t, 3600, ts["duration"])
			},
		},
		{
			name: "with time range timespan",
			input: func() dashboards.ApiDashboard {
				d := dashboards.NewApiDashboard()
				d.SetTitle("Test Dashboard")
				d.SetDescription("Test Description")
				ts := dashboards.NewDefaultTimespan()
				startTime, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2026-02-01T00:00:00Z")
				ts.SetStart(startTime)
				ts.SetEnd(endTime)
				d.SetDefaultTimespan(*ts)
				return *d
			}(),
			validate: func(t *testing.T, d *schema.ResourceData) {
				timespan := d.Get("default_timespan").([]interface{})
				assert.Len(t, timespan, 1)
				ts := timespan[0].(map[string]interface{})
				assert.Equal(t, "2026-01-01T00:00:00Z", ts["start"])
				assert.Equal(t, "2026-02-01T00:00:00Z", ts["end"])
			},
		},
		{
			name: "with modified date",
			input: func() dashboards.ApiDashboard {
				d := dashboards.NewApiDashboard()
				d.SetTitle("Test Dashboard")
				d.SetDescription("Test Description")
				modDate, _ := time.Parse(time.RFC3339, "2026-03-15T10:30:00Z")
				d.SetDashboardModifiedDate(modDate)
				return *d
			}(),
			validate: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "2026-03-15T10:30:00Z", d.Get("dashboard_modified_date"))
			},
		},
		{
			name: "without default_timespan sets empty list",
			input: func() dashboards.ApiDashboard {
				d := dashboards.NewApiDashboard()
				d.SetTitle("Test Dashboard")
				d.SetDescription("Test Description")
				// No default_timespan set
				return *d
			}(),
			validate: func(t *testing.T, d *schema.ResourceData) {
				timespan := d.Get("default_timespan").([]interface{})
				assert.Empty(t, timespan)
			},
		},
		{
			name: "with zero modified date does not set field",
			input: func() dashboards.ApiDashboard {
				d := dashboards.NewApiDashboard()
				d.SetTitle("Test Dashboard")
				d.SetDescription("Test Description")
				// Zero time - should not set dashboard_modified_date
				return *d
			}(),
			validate: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "", d.Get("dashboard_modified_date"))
			},
		},
		{
			name: "with empty timespan map does not set default_timespan",
			input: func() dashboards.ApiDashboard {
				d := dashboards.NewApiDashboard()
				d.SetTitle("Test Dashboard")
				d.SetDescription("Test Description")
				// Empty timespan with zero duration
				ts := dashboards.NewDefaultTimespan()
				ts.SetDuration(0)
				d.SetDefaultTimespan(*ts)
				return *d
			}(),
			validate: func(t *testing.T, d *schema.ResourceData) {
				// Empty timespan map should result in empty list
				timespan := d.Get("default_timespan").([]interface{})
				assert.Empty(t, timespan)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, schemas.DashboardSchema, map[string]interface{}{})
			err := resourceDataApiDashboardMapper(d, tc.input)
			assert.NoError(t, err)
			tc.validate(t, d)
		})
	}
}
