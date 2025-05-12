package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func resourceAccountGroup() *schema.Resource {
	resource := schema.Resource{
		Schema: schemas.AccountGroupSchema,
		Create: resourceAccountGroupCreate,
		Read:   resourceAccountGroupRead,
		Update: resourceAccountGroupUpdate,
		Delete: resourceAccountGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Account groups are used to divide an organization into different sections. These operations can be used to create, retrieve, update and delete account groups. For more information, see [Account Groups](https://developer.cisco.com/docs/thousandeyes/list-account-groups/).",
	}
	return &resource
}

func resourceAccountGroupRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*administrative.AccountGroupsAPIService)(&apiClient.Common)

		req := api.GetAccountGroup(id).Expand(administrative.AllowedExpandAccountGroupOptionsEnumValues)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceAccountGroupUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*administrative.AccountGroupsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Account Group %s", d.Id())
	update := buildAccountGroupStruct(d)

	req := api.UpdateAccountGroup(d.Id()).AccountGroupRequest(*update).Expand(administrative.AllowedExpandAccountGroupOptionsEnumValues)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}

	return resourceAccountGroupRead(d, m)
}

func resourceAccountGroupDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*administrative.AccountGroupsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Account Group %s", d.Id())

	req := api.DeleteAccountGroup(d.Id())

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAccountGroupCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*administrative.AccountGroupsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Account Group %s", d.Id())
	local := buildAccountGroupStruct(d)

	req := api.CreateAccountGroup().AccountGroupRequest(*local).Expand(administrative.AllowedExpandAccountGroupOptionsEnumValues)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.Aid
	d.SetId(id)

	// Using a response because account group isn't available right after creation
	return ResourceRead(context.Background(), d, buildAccountGroupReadStruct(resp, local.Agents))
}

func buildAccountGroupStruct(d *schema.ResourceData) *administrative.AccountGroupRequest {
	return ResourceBuildStruct(d, &administrative.AccountGroupRequest{})
}

// POST response object do not have agents, but GET response does
func buildAccountGroupReadStruct(in *administrative.CreatedAccountGroup, agents []string) *administrative.AccountGroupDetail {
	out := &administrative.AccountGroupDetail{
		Aid:                   in.Aid,
		AccountGroupName:      in.AccountGroupName,
		IsCurrentAccountGroup: in.IsCurrentAccountGroup,
		IsDefaultAccountGroup: in.IsDefaultAccountGroup,
		OrganizationName:      in.OrganizationName,
		OrgId:                 in.OrgId,
		Users:                 in.Users,
		Links:                 in.Links,
	}
	if len(agents) > 0 {
		agentsArr := make([]administrative.EnterpriseAgent, 0, len(agents))
		for _, v := range agents {
			agentsArr = append(agentsArr, administrative.EnterpriseAgent{AgentId: getPointer(v)})
		}
		out.Agents = agentsArr
	}
	return out
}
