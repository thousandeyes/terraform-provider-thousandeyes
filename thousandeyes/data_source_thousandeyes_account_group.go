package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func dataSourceThousandeyesAccountGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesAccountGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account group.",
			},
			"aid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique ID for the account group.",
			},
		},
		Description: "This data source allows you to configure a ThousandEyes account group. For more information, see [What is an Account Group](https://docs.thousandeyes.com/product-documentation/user-management/account-groups/what-is-an-account-group).",
	}
}

func dataSourceThousandeyesAccountGroupRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*administrative.AccountGroupsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Reading Thousandeyes account group")

	searchName := d.Get("name").(string)

	resp, _, err := api.GetAccountGroups().Execute()
	if err != nil {
		return err
	}

	var found *administrative.AccountGroupInfo

	for _, account := range resp.GetAccountGroups() {
		if *account.AccountGroupName == searchName {
			found = &account
			break
		}
	}

	if found == nil {
		return fmt.Errorf("unable to locate any account group with the name: %s", searchName)
	}
	log.Printf("[INFO] ## Found AccountGroup ID: %s - name: %s", *found.Aid, *found.AccountGroupName)

	d.SetId(*found.Aid)
	err = d.Set("name", found.AccountGroupName)
	if err != nil {
		return err
	}
	err = d.Set("aid", found.Aid)
	if err != nil {
		return err
	}

	return nil
}
