package thousandeyes

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func resourceDashboard() *schema.Resource {
	resource := schema.Resource{
		Schema: schemas.DashboardSchema,
		Create: resourceDashboardCreate,
		Read:   resourceDashboardRead,
		Update: resourceDashboardUpdate,
		Delete: resourceDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "This resource allows you to create and manage Dashboards. For more information, see [Dashboards](https://developer.cisco.com/docs/thousandeyes/list-dashboards/).",
	}
	return &resource
}

func resourceDashboardCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*dashboards.DashboardsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Dashboard %s", d.Get("title"))
	local := buildDashboardStruct(d)

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

	log.Printf("[INFO] Updating ThousandEyes Dashboard %s", d.Id())
	update := buildDashboardStruct(d)

	req := api.UpdateDashboard(d.Id()).Dashboard(*update)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()

	if err != nil {
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

func buildDashboardStruct(d *schema.ResourceData) *dashboards.Dashboard {
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

			// duration comes as int from schema, cast appropriately
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

	return dashboard
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

	return nil
}
