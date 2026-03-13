package thousandeyes

import (
	"context"
	"log"

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

		return buildDashboardFromReadResponse(dashboards.NewDashboard(), resp), err
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

// Map to dashboards.Dashboard as this only has fields that exist in the DashboardSchema
func buildDashboardFromReadResponse(dashboard *dashboards.Dashboard, resp *dashboards.ApiDashboard) *dashboards.Dashboard {
	dashboard.SetDashboardId(resp.GetDashboardId())
	dashboard.SetDescription(resp.GetDescription())
	dashboard.SetTitle(resp.GetTitle())
	dashboard.SetCreatedBy(resp.GetDashboardCreatedBy())
	dashboard.SetModifiedBy(resp.GetDashboardModifiedBy())
	dashboard.SetModifiedDate(resp.GetDashboardModifiedDate())
	dashboard.SetGlobalFilterId(resp.GetGlobalFilterId())
	dashboard.SetIsBuiltIn(resp.GetIsBuiltIn())
	dashboard.SetIsDefaultForAccount(resp.GetIsDefaultForAccount())
	dashboard.SetIsDefaultForUser(resp.GetIsDefaultForUser())
	dashboard.SetIsGlobalOverride(resp.GetIsGlobalOverride())
	dashboard.SetIsMigratedReport(resp.GetIsMigratedReport())
	dashboard.SetIsPrivate(resp.GetIsPrivate())
	dashboard.SetDefaultTimespan(resp.GetDefaultTimespan())
	return dashboard
}
