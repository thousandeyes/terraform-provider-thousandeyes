package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func resourceUser() *schema.Resource {
	resource := schema.Resource{
		Schema: schemas.UserSchema,
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "",
	}
	return &resource
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*administrative.UsersAPIService)(&apiClient.Common)

		req := api.GetUser(id)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*administrative.UsersAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes User %s", d.Id())
	update := buildUserStruct(d)

	req := api.UpdateUser(d.Id()).UserRequest(*update)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*administrative.UsersAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes User %s", d.Id())

	req := api.DeleteUser(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*administrative.UsersAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes User %s", d.Id())
	local := buildUserStruct(d)

	req := api.CreateUser().UserRequest(*local)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.Uid
	d.SetId(id)
	return resourceUserRead(d, m)
}

func buildUserStruct(d *schema.ResourceData) *administrative.UserRequest {
	return ResourceBuildStruct(d, &administrative.UserRequest{})
}
