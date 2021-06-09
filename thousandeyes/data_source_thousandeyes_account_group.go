package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func dataSourceThousandeyesAccountGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesAccountGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
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
		if account.AccountGroupName == searchName {
			found = &account
			break
		}
	}

	if found == nil {
		return fmt.Errorf("unable to locate any account group with the name: %s", searchName)
	}
	log.Printf("[INFO] ## Found AccountGroup ID: %d - name: %s", found.AID, found.AccountGroupName)

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
