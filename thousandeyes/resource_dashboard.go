package thousandeyes

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func resourceDashboard() *schema.Resource {
	resource := schema.Resource{
		Schema:        schemas.DashboardSchema,
		Create:        resourceDashboardCreate,
		Read:          resourceDashboardRead,
		Update:        resourceDashboardUpdate,
		Delete:        resourceDashboardDelete,
		CustomizeDiff: resourceDashboardCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "This resource allows you to create and manage Dashboards. For more information, see [Dashboards](https://developer.cisco.com/docs/thousandeyes/list-dashboards/).",
	}
	return &resource
}

func resourceDashboardCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Optional+Computed TypeList: Terraform merges prior state into the plan when no widget
	// blocks are present in config, suppressing any diff. Force an empty planned value so
	// Terraform detects the removal and calls Update.
	//
	// Guard on len(oldList) > 0 to avoid a spurious diff on new resources where both
	// config and state are already empty.
	if configWidgetsNullOrEmpty(d.GetRawConfig()) {
		old, _ := d.GetChange("widgets")
		oldList, _ := old.([]interface{})
		if len(oldList) > 0 {
			if err := d.SetNew("widgets", []interface{}{}); err != nil {
				return err
			}
		}
	}

	raw := d.Get("widgets")
	if raw == nil {
		return nil
	}
	widgets, ok := raw.([]interface{})
	if !ok {
		return fmt.Errorf("widgets: expected a list of objects, got %T", raw)
	}

	for i, w := range widgets {
		if w == nil {
			return fmt.Errorf("widgets[%d]: must not be null", i)
		}
		widget, ok := w.(map[string]interface{})
		if !ok {
			return fmt.Errorf("widgets[%d]: expected object, got %T", i, w)
		}
		typeVal, ok := widget["type"]
		if !ok || typeVal == nil {
			return fmt.Errorf("widgets[%d]: type is required", i)
		}
		widgetType, ok := typeVal.(string)
		if !ok || widgetType == "" {
			return fmt.Errorf("widgets[%d]: type must be a non-empty string", i)
		}

		if widgetType == WidgetTypeStackedArea {
			rawCfg := widget["stacked_area_config"]
			stackedAreaConfig, ok := rawCfg.([]interface{})
			if !ok {
				return fmt.Errorf("widgets[%d]: stacked_area_config must be a list for widget type %q", i, WidgetTypeStackedArea)
			}
			if len(stackedAreaConfig) == 0 {
				return fmt.Errorf("widgets[%d]: stacked_area_config is required for widget type '%s'", i, WidgetTypeStackedArea)
			}
			config, ok := stackedAreaConfig[0].(map[string]interface{})
			if !ok || config == nil {
				return fmt.Errorf("widgets.%d: stacked_area_config.0 must be an object", i)
			}
			if groupBy, ok := config["group_by"].(string); !ok || groupBy == "" {
				return fmt.Errorf("widgets[%d].stacked_area_config.group_by is required for widget type '%s'", i, WidgetTypeStackedArea)
			}
		}

		if widgetType == WidgetTypePieChart {
			rawCfg := widget["pie_chart_config"]
			pieChartConfig, ok := rawCfg.([]interface{})
			if !ok {
				return fmt.Errorf("widgets[%d]: pie_chart_config must be a list for widget type %q", i, WidgetTypePieChart)
			}
			if len(pieChartConfig) == 0 {
				return fmt.Errorf("widgets[%d]: pie_chart_config is required for widget type '%s'", i, WidgetTypePieChart)
			}
			config, ok := pieChartConfig[0].(map[string]interface{})
			if !ok || config == nil {
				return fmt.Errorf("widgets.%d: pie_chart_config.0 must be an object", i)
			}
			if groupBy, ok := config["group_by"].(string); !ok || groupBy == "" {
				return fmt.Errorf("widgets[%d].pie_chart_config.group_by is required for widget type '%s'", i, WidgetTypePieChart)
			}
		}
	}

	return nil
}

// configWidgetsNullOrEmpty reports whether widgets is absent (null) or explicitly empty in
// configuration. This covers both ways a user can express "no widget blocks":
//   - Omitting all widget blocks → Terraform sets the attribute to null for Optional+Computed lists
//   - Writing zero blocks → Terraform produces an empty tuple/list
func configWidgetsNullOrEmpty(cfg cty.Value) bool {
	w, ok := widgetsFromRawConfig(cfg)
	if !ok {
		return false
	}
	if w.IsNull() {
		return true
	}
	return w.IsKnown() &&
		(w.Type().IsListType() || w.Type().IsTupleType()) &&
		w.LengthInt() == 0
}

// widgetsFromRawConfig extracts the widgets cty.Value from raw config. Returns false when
// raw config is unavailable or not a proper object (e.g. some test helpers).
func widgetsFromRawConfig(cfg cty.Value) (cty.Value, bool) {
	if cfg == cty.NilVal || cfg.IsNull() || !cfg.IsKnown() {
		return cty.NilVal, false
	}
	ty := cfg.Type()
	if !ty.IsObjectType() || !ty.HasAttribute("widgets") {
		return cty.NilVal, false
	}
	return cfg.GetAttr("widgets"), true
}

func resourceDashboardCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*dashboards.DashboardsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Dashboard %s", d.Get("title"))
	local, err := buildDashboardStruct(d)
	if err != nil {
		return err
	}

	req := api.CreateDashboard().Dashboard(*local)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := resp.GetDashboardId()
	d.SetId(id)
	return resourceDashboardRead(d, m)
}

func resourceDashboardRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*dashboards.DashboardsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Reading ThousandEyes Dashboard %s", d.Id())
	req := api.GetDashboard(d.Id())
	ctx := apiClient.GetConfig().Context
	req = SetAidFromContext(ctx, req)

	resp, httpResp, err := req.Execute()
	if err != nil {

		// if the dashboard doesn't exist, return nil and remove from state
		if httpResp != nil && httpResp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	return resourceDataApiDashboardMapper(d, *resp)
}

func resourceDashboardUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*dashboards.DashboardsAPIService)(&apiClient.Common)
	ctx := apiClient.GetConfig().Context

	log.Printf("[INFO] Updating ThousandEyes Dashboard %s", d.Id())
	update, err := buildDashboardStruct(d)
	if err != nil {
		return err
	}

	getReq := api.GetDashboard(d.Id())
	getReq = SetAidFromContext(ctx, getReq)
	existing, _, err := getReq.Execute()
	if err != nil {
		return fmt.Errorf("fetching current dashboard for merge: %w", err)
	}

	merged := mergeUnmanagedWidgets(update.GetWidgets(), existing.GetWidgets())
	update.SetWidgets(merged)

	req := api.UpdateDashboard(d.Id()).Dashboard(*update)
	req = SetAidFromContext(ctx, req)

	_, _, err = req.Execute()
	if err != nil {
		log.Printf("[ERROR] API error updating ThousandEyes Dashboard %s: %v", d.Id(), err)
		return err
	}

	return resourceDashboardRead(d, m)
}

func resourceDashboardDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*dashboards.DashboardsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Dashboard %s", d.Id())

	req := api.DeleteDashboard(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func buildDashboardStruct(d *schema.ResourceData) (*dashboards.Dashboard, error) {
	dashboard := &dashboards.Dashboard{}

	if v, ok := d.GetOk("title"); ok {
		dashboard.SetTitle(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		dashboard.SetDescription(v.(string))
	}
	if v, ok := d.GetOk("is_private"); ok {
		dashboard.SetIsPrivate(v.(bool))
	}
	if v, ok := d.GetOk("global_filter_id"); ok {
		dashboard.SetGlobalFilterId(v.(string))
	}
	if v, ok := d.GetOk("is_global_override"); ok {
		dashboard.SetIsGlobalOverride(v.(bool))
	}

	if v, ok := d.GetOk("default_timespan"); ok {
		timespanList := v.([]interface{})
		if len(timespanList) > 0 {
			ts := timespanList[0].(map[string]interface{})
			t := dashboards.DefaultTimespan{}

			if dur, ok := ts["duration"].(int); ok && dur != 0 {
				t.SetDuration(int64(dur))
			}
			if startStr, ok := ts["start"].(string); ok && startStr != "" {
				if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
					t.SetStart(startTime)
				} else {
					log.Printf("[WARN] Invalid start time: %s", startStr)
				}
			}
			if endStr, ok := ts["end"].(string); ok && endStr != "" {
				if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
					t.SetEnd(endTime)
				} else {
					log.Printf("[WARN] Invalid end time: %s", endStr)
				}
			}
			dashboard.SetDefaultTimespan(t)
		}
	}

	// Always send an explicit widget slice — populated or empty — so the API
	// performs a true replacement. For Optional+Computed TypeList attributes,
	// d.Get returns the prior-state value when no blocks are written in config.
	// We read the raw config instead to get the authoritative intended value.
	var widgetList []interface{}
	if configWidgetsNullOrEmpty(d.GetRawConfig()) {
		widgetList = []interface{}{}
	} else {
		widgetList, _ = d.Get("widgets").([]interface{})
	}
	widgets, err := BuildWidgets(widgetList)
	if err != nil {
		return nil, err
	}
	if widgets == nil {
		widgets = []dashboards.ApiWidget{}
	}
	dashboard.SetWidgets(widgets)

	return dashboard, nil
}

func resourceDataApiDashboardMapper(d *schema.ResourceData, dashboard dashboards.ApiDashboard) error {
	if err := d.Set("aid", dashboard.GetAid()); err != nil {
		return err
	}
	if err := d.Set("description", dashboard.GetDescription()); err != nil {
		return err
	}
	if err := d.Set("title", dashboard.GetTitle()); err != nil {
		return err
	}
	if err := d.Set("dashboard_created_by", dashboard.GetDashboardCreatedBy()); err != nil {
		return err
	}
	if err := d.Set("dashboard_modified_by", dashboard.GetDashboardModifiedBy()); err != nil {
		return err
	}
	if modDate := dashboard.GetDashboardModifiedDate(); !modDate.IsZero() {
		if err := d.Set("dashboard_modified_date", modDate.Format(time.RFC3339)); err != nil {
			return err
		}
	}

	if err := d.Set("global_filter_id", dashboard.GetGlobalFilterId()); err != nil {
		return err
	}
	if err := d.Set("is_global_override", dashboard.GetIsGlobalOverride()); err != nil {
		return err
	}
	if err := d.Set("is_migrated_report", dashboard.GetIsMigratedReport()); err != nil {
		return err
	}
	if err := d.Set("is_private", dashboard.GetIsPrivate()); err != nil {
		return err
	}

	if timespan, ok := dashboard.GetDefaultTimespanOk(); ok && timespan != nil {
		t := map[string]any{}

		hasStart := !timespan.GetStart().IsZero()
		hasEnd := !timespan.GetEnd().IsZero()

		// Only set duration if we're not in time range mode (start/end)
		if !hasStart && !hasEnd {
			if dur, ok := timespan.GetDurationOk(); ok && dur != nil && *dur != 0 {
				t["duration"] = *dur
			}
		}

		if hasStart {
			t["start"] = timespan.GetStart().Format(time.RFC3339)
		}

		if hasEnd {
			t["end"] = timespan.GetEnd().Format(time.RFC3339)
		}

		if len(t) > 0 {
			if err := d.Set("default_timespan", []interface{}{t}); err != nil {
				return err
			}
		}
	} else {
		err := d.Set("default_timespan", nil)
		if err != nil {
			return err
		}
	}

	// Handle widgets
	if widgets := dashboard.GetWidgets(); len(widgets) > 0 {
		mappedWidgets, err := MapWidgets(widgets)
		if err != nil {
			return err
		}
		if err := d.Set("widgets", mappedWidgets); err != nil {
			return err
		}
	} else {
		if err := d.Set("widgets", []interface{}{}); err != nil {
			return err
		}
	}

	return nil
}
