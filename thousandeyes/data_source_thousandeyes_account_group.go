package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func dataSourceThousandeyesAccountGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesAccountGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The name of the account group.",
			},
			"aid": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "The unique ID for the account group.",
			},
		},
		Description: "This data source allows you to configure a ThousandEyes account group. For more information, see [What is an Account Group](https://docs.thousandeyes.com/product-documentation/user-management/account-groups/what-is-an-account-group).",
	}
}

func dataSourceThousandeyesAccountGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes account group")

	searchName := d.Get("name").(string)

	accounts, err := client.GetAccountGroups()
	if err != nil {
		return err
	}

	var found *thousandeyes.SharedWithAccount

	for _, account := range *accounts {
		if *account.AccountGroupName == searchName {
			found = &account
			break
		}
	}

	if found == nil {
		return fmt.Errorf("unable to locate any account group with the name: %s", searchName)
	}
	log.Printf("[INFO] ## Found AccountGroup ID: %d - name: %s", found.AID, *found.AccountGroupName)

	d.SetId(fmt.Sprint(found.AID))
	err = d.Set("name", found.AccountGroupName)
	if err != nil {
		return err
	}
	err = d.Set("aid", found.AID)
	if err != nil {
		return err
	}

	return nil
}
