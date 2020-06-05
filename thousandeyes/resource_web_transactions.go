package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func resourceWebTransaction() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.WebTransaction{}, schemas),
		Create: resourceWebTransactionCreate,
		Read:   resourceWebTransactionRead,
		Update: resourceWebTransactionUpdate,
		Delete: resourceWebTransactionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	return &resource
}

func resourceWebTransactionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetWebTransaction(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourceWebTransactionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
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
	id, _ := strconv.Atoi(d.Id())
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
	id := remote.TestID
	d.SetId(strconv.Itoa(id))
	return resourceWebTransactionRead(d, m)
}

func buildWebTransactionStruct(d *schema.ResourceData) *thousandeyes.WebTransaction {
	return ResourceBuildStruct(d, &thousandeyes.WebTransaction{}).(*thousandeyes.WebTransaction)
}
