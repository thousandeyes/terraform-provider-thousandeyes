package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func resourceGroupLabel() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.GroupLabel{}, schemas),
		Create: resourceGroupLabelCreate,
		Read:   resourceGroupLabelRead,
		Update: resourceGroupLabelUpdate,
		Delete: resourceGroupLabelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	resource.Schema["type"] = schemas["type-label"]
	resource.Schema["agents"] = schemas["agents-label"]
	return &resource
}

func resourceGroupLabelRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Label %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetGroupLabel(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourceGroupLabelUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Label %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.GroupLabel{}).(*thousandeyes.GroupLabel)
	_, err := client.UpdateGroupLabel(id, *update)
	if err != nil {
		return err
	}
	return resourceGroupLabelRead(d, m)
}

func resourceGroupLabelDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Label %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteGroupLabel(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceGroupLabelCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Label %s", d.Id())
	local := buildGroupLabelStruct(d)
	remote, err := client.CreateGroupLabel(*local)
	if err != nil {
		return err
	}
	id := remote.GroupID
	d.SetId(strconv.Itoa(id))
	return resourceGroupLabelRead(d, m)
}

func buildGroupLabelStruct(d *schema.ResourceData) *thousandeyes.GroupLabel {
	return ResourceBuildStruct(d, &thousandeyes.GroupLabel{}).(*thousandeyes.GroupLabel)
}
