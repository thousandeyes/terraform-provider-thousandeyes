package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func resourceRole() *schema.Resource {
	resource := schema.Resource{
		Schema: schemas.RoleSchema,
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Create, retrieve and update roles for the current user. For more information, see [Roles](https://developer.cisco.com/docs/thousandeyes/list-roles/).",
	}
	return &resource
}

func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*administrative.RolesAPIService)(&apiClient.Common)

		req := api.GetRole(id)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*administrative.RolesAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Role %s", d.Id())
	update := buildRoleStruct(d)

	req := api.UpdateRole(d.Id()).RoleRequestBody(*update)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}

	return resourceRoleRead(d, m)
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*administrative.RolesAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Role %s", d.Id())

	req := api.DeleteRole(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*administrative.RolesAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Role %s", d.Id())
	local := buildRoleStruct(d)

	req := api.CreateRole().RoleRequestBody(*local)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.RoleId
	d.SetId(id)
	return resourceRoleRead(d, m)
}

func buildRoleStruct(d *schema.ResourceData) *administrative.RoleRequestBody {
	return ResourceBuildStruct(d, &administrative.RoleRequestBody{})
}
