package thousandeyes

import (
	"context"
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

	id := *resp.DashboardId
	d.SetId(id)
	return resourceDashboardRead(d, m)
}

func resourceDashboardRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*dashboards.DashboardsAPIService)(&apiClient.Common)

		log.Printf("[INFO] Reading ThousandEyes Dashboard %s", d.Id())
		req := api.GetDashboard(id)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		if err != nil {
			return nil, err
		}

		return buildDashboardStructFromApiResponse(resp), err
	})
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
	return ResourceBuildStruct(d, &dashboards.Dashboard{})
}

// Separate Resource fields API responses by mapping dashboards.ApiDashboard to DashboardResourceData
func buildDashboardStructFromApiResponse(resp *dashboards.ApiDashboard) *DashboardResourceData {
	ds := &DashboardResourceData{
		Aid:                   resp.GetAid(),
		DashboardCreatedBy:    resp.GetDashboardCreatedBy(),
		DashboardModifiedBy:   resp.GetDashboardModifiedBy(),
		DashboardModifiedDate: resp.GetDashboardModifiedDate(),
		Description:           resp.GetDescription(),
		GlobalFilterId:        resp.GetGlobalFilterId(),
		Title:                 resp.GetTitle(),
		IsBuiltIn:             resp.GetIsBuiltIn(),
		IsDefaultForAccount:   resp.GetIsDefaultForAccount(),
		IsDefaultForUser:      resp.GetIsDefaultForUser(),
		IsGlobalOverride:      resp.GetIsGlobalOverride(),
		IsMigratedReport:      resp.GetIsMigratedReport(),
		IsPrivate:             resp.GetIsPrivate(),
	}

	if timespan, ok := resp.GetDefaultTimespanOk(); ok && timespan != nil {
		ds.DefaultTimespan = &DefaultTimespanStruct{
			Duration: timespan.GetDuration(),
			Start:    timespan.GetStart(),
			End:      timespan.GetEnd(),
		}
	}

	return ds
}

type DashboardResourceData struct {
	Aid                   string
	DashboardCreatedBy    string
	DashboardModifiedBy   string
	DashboardModifiedDate time.Time
	Description           string
	GlobalFilterId        string
	Title                 string
	IsBuiltIn             bool
	IsDefaultForAccount   bool
	IsDefaultForUser      bool
	IsGlobalOverride      bool
	IsMigratedReport      bool
	IsPrivate             bool
	DefaultTimespan       *DefaultTimespanStruct
}

type DefaultTimespanStruct struct {
	Duration int64
	Start    time.Time
	End      time.Time
}
