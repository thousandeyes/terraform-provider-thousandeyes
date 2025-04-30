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
		Description: "This resource allows you to create and configure an agent-to-agent test. This test type evaluates the performance of the underlying network between two physical sites. For more information about agent-to-agent tests, see [Agent-to-Agent Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#agent-to-agent-test).",
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
	return resourceAccountGroupRead(d, m)
}

func buildAccountGroupStruct(d *schema.ResourceData) *administrative.AccountGroupRequest {
	return ResourceBuildStruct(d, &administrative.AccountGroupRequest{})
}
