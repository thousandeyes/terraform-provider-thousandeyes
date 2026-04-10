package thousandeyes

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

// dashboardDiff runs the SDK schema Diff (not TF_ACC / terraform CLI).
func dashboardDiff(t *testing.T, state *terraform.InstanceState, cfg *terraform.ResourceConfig) *terraform.InstanceDiff {
	t.Helper()
	diff, err := resourceDashboard().Diff(context.Background(), state, cfg, nil)
	require.NoError(t, err)
	return diff
}

func assertNoDefaultTimespanAttributeDiff(t *testing.T, diff *terraform.InstanceDiff) {
	t.Helper()
	if diff == nil {
		return
	}
	for attr := range diff.Attributes {
		if strings.HasPrefix(attr, "default_timespan") {
			t.Fatalf("unexpected plan diff on %q: %#v", attr, diff.Attributes[attr])
		}
	}
}

func dashboardResourceDataAfterRead(t *testing.T, api *dashboards.ApiDashboard) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, schemas.DashboardSchema, map[string]interface{}{
		"title":              api.GetTitle(),
		"description":        api.GetDescription(),
		"is_private":         api.GetIsPrivate(),
		"is_global_override": api.GetIsGlobalOverride(),
	})
	d.SetId("diff-test-dashboard-id")
	require.NoError(t, resourceDataApiDashboardMapper(d, *api))
	return d
}

// TestDashboardDefaultTimespan_omitFromConfigNoPlanDiff: CP-4085 — omit block in config while
// state holds API default timespan; plan must not try to change default_timespan.
func TestDashboardDefaultTimespan_omitFromConfigNoPlanDiff(t *testing.T) {
	api := dashboards.NewApiDashboard()
	api.SetTitle("Diff Test Dashboard")
	api.SetDescription("diff test")
	api.SetIsPrivate(false)
	api.SetIsGlobalOverride(false)
	ts := dashboards.NewDefaultTimespan()
	ts.SetDuration(7200)
	api.SetDefaultTimespan(*ts)

	d := dashboardResourceDataAfterRead(t, api)
	state := d.State()
	require.Equal(t, "7200", state.Attributes["default_timespan.0.duration"])

	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
		"title":              "Diff Test Dashboard",
		"description":        "diff test",
		"is_private":         false,
		"is_global_override": false,
	})

	diff := dashboardDiff(t, state, cfg)
	assertNoDefaultTimespanAttributeDiff(t, diff)
}

// TestDashboardDefaultTimespan_explicitDurationMatchesStateNoPlanDiff: user sets duration in
// config matching refreshed state — no drift on default_timespan.
func TestDashboardDefaultTimespan_explicitDurationMatchesStateNoPlanDiff(t *testing.T) {
	const seconds = 21600

	api := dashboards.NewApiDashboard()
	api.SetTitle("Explicit Timespan Dashboard")
	api.SetDescription("diff test")
	api.SetIsPrivate(false)
	api.SetIsGlobalOverride(false)
	ts := dashboards.NewDefaultTimespan()
	ts.SetDuration(seconds)
	api.SetDefaultTimespan(*ts)

	d := dashboardResourceDataAfterRead(t, api)
	state := d.State()
	require.Equal(t, "21600", state.Attributes["default_timespan.0.duration"])

	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
		"title":              "Explicit Timespan Dashboard",
		"description":        "diff test",
		"is_private":         false,
		"is_global_override": false,
		"default_timespan": []interface{}{
			map[string]interface{}{"duration": seconds},
		},
	})

	diff := dashboardDiff(t, state, cfg)
	assertNoDefaultTimespanAttributeDiff(t, diff)
}

// TestDashboardDefaultTimespan_removeBlockFromConfigKeepsStateNoPlanDiff documents current
// behavior: after an explicit timespan was applied, removing the block from config does not
// produce a plan diff on default_timespan (Optional+Computed; build path uses GetOk, no
// raw-config clear). The API value remains reflected in state.
func TestDashboardDefaultTimespan_removeBlockFromConfigKeepsStateNoPlanDiff(t *testing.T) {
	const seconds = 21600

	api := dashboards.NewApiDashboard()
	api.SetTitle("Remove Block Dashboard")
	api.SetDescription("diff test")
	api.SetIsPrivate(false)
	api.SetIsGlobalOverride(false)
	ts := dashboards.NewDefaultTimespan()
	ts.SetDuration(seconds)
	api.SetDefaultTimespan(*ts)

	d := dashboardResourceDataAfterRead(t, api)
	state := d.State()

	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
		"title":              "Remove Block Dashboard",
		"description":        "diff test",
		"is_private":         false,
		"is_global_override": false,
	})

	diff := dashboardDiff(t, state, cfg)
	assertNoDefaultTimespanAttributeDiff(t, diff)
}
