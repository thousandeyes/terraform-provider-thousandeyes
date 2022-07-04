package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceVoiceCall() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.VoiceCall{}, schemas),
		Create: resourceVoiceCallCreate,
		Read:   resourceVoiceCallRead,
		Update: resourceVoiceCallUpdate,
		Delete: resourceVoiceCallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	return &resource
}

func resourceVoiceCallRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetVoiceCall(id)
	if err != nil {
		return err
	}
	err = ResourceRead(d, remote)
	if err != nil {
		return err
	}
	return nil
}

func resourceVoiceCallUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.VoiceCall{}).(*thousandeyes.VoiceCall)
	_, err := client.UpdateVoiceCall(id, *update)
	if err != nil {
		return err
	}
	return resourceVoiceCallRead(d, m)
}

func resourceVoiceCallDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteVoiceCall(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceVoiceCallCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildVoiceCallStruct(d)
	remote, err := client.CreateVoiceCall(*local)
	if err != nil {
		return err
	}

	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceVoiceCallRead(d, m)
}

func buildVoiceCallStruct(d *schema.ResourceData) *thousandeyes.VoiceCall {
	return ResourceBuildStruct(d, &thousandeyes.VoiceCall{}).(*thousandeyes.VoiceCall)
}
