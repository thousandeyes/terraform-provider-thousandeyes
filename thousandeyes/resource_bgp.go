package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func resourceBGP() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.BGP{}, schemas),
		Create: resourceBGPCreate,
		Read:   resourceBGPRead,
		Update: resourceBGPUpdate,
		Delete: resourceBGPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	return &resource
}

func resourceBGPRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetBGP(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourceBGPUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.BGP{}).(*thousandeyes.BGP)
	_, err := client.UpdateBGP(id, *update)
	if err != nil {
		return err
	}
	return resourceBGPRead(d, m)
}

func resourceBGPDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteBGP(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceBGPCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildBGPStruct(d)
	remote, err := client.CreateBGP(*local)
	if err != nil {
		return err
	}
	id := remote.TestID
	d.SetId(strconv.Itoa(id))
	return resourceBGPRead(d, m)
}

func buildBGPStruct(d *schema.ResourceData) *thousandeyes.BGP {
	return ResourceBuildStruct(d, &thousandeyes.BGP{}).(*thousandeyes.BGP)
}
