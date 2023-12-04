package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceWebTransaction() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.WebTransaction{}, schemas, nil),
		Create: resourceWebTransactionCreate,
		Read:   resourceWebTransactionRead,
		Update: resourceWebTransactionUpdate,
		Delete: resourceWebTransactionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows users to create a transaction test. This test type is a scripted synthetic browser interaction that can traverse multiple pages and user actions. For more information, see [Transaction Tests](https://docs.thousandeyes.com/product-documentation/browser-synthetics/transaction-tests).",
	}
	return &resource
}

func resourceWebTransactionRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(client *thousandeyes.Client, id int64) (interface{}, error) {
		return client.GetWebTransaction(id)
	})
}

func resourceWebTransactionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	update := ResourceUpdate(d, &thousandeyes.WebTransaction{}).(*thousandeyes.WebTransaction)
	_, err := client.UpdateWebTransaction(id, *update)
	if err != nil {
		return err
	}
	return resourceWebTransactionRead(d, m)
}

func resourceWebTransactionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	if err := client.DeleteWebTransaction(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceWebTransactionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildWebTransactionStruct(d)
	remote, err := client.CreateWebTransaction(*local)
	if err != nil {
		return err
	}
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceWebTransactionRead(d, m)
}

func buildWebTransactionStruct(d *schema.ResourceData) *thousandeyes.WebTransaction {
	return ResourceBuildStruct(d, &thousandeyes.WebTransaction{}).(*thousandeyes.WebTransaction)
}
