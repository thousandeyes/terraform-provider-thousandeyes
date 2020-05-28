package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func resourcePageLoad() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.PageLoad{}, schemas),
		Create: resourcePageLoadCreate,
		Read:   resourcePageLoadRead,
		Update: resourcePageLoadUpdate,
		Delete: resourcePageLoadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	return &resource
}

func resourcePageLoadRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetPageLoad(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourcePageLoadUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.PageLoad{}).(*thousandeyes.PageLoad)
	_, err := client.UpdatePageLoad(id, *update)
	if err != nil {
		return err
	}
	return resourcePageLoadRead(d, m)
}

func resourcePageLoadDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeletePageLoad(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourcePageLoadCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildPageLoadStruct(d)
	remote, err := client.CreatePageLoad(*local)
	if err != nil {
		return err
	}
	id := remote.TestID
	d.SetId(strconv.Itoa(id))
	return resourcePageLoadRead(d, m)
}

func buildPageLoadStruct(d *schema.ResourceData) *thousandeyes.PageLoad {
	return ResourceBuildStruct(d, &thousandeyes.PageLoad{}).(*thousandeyes.PageLoad)
}
