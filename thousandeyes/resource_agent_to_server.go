package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func resourceAgentToServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.AgentServer{}, schemas),
		Create: resourceAgentServerCreate,
		Read:   resourceAgentServerRead,
		Update: resourceAgentServerUpdate,
		Delete: resourceAgentServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	return &resource
}

func resourceAgentServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetAgentServer(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourceAgentServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.AgentServer{}).(*thousandeyes.AgentServer)
	_, err := client.UpdateAgentServer(id, *update)
	if err != nil {
		return err
	}
	return resourceAgentServerRead(d, m)
}

func resourceAgentServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteAgentServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAgentServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildAgentServerStruct(d)
	remote, err := client.CreateAgentServer(*local)
	if err != nil {
		return err
	}
	id := remote.TestID
	d.SetId(strconv.Itoa(id))
	return resourceAgentServerRead(d, m)
}

func buildAgentServerStruct(d *schema.ResourceData) *thousandeyes.AgentServer {
	return ResourceBuildStruct(d, &thousandeyes.AgentServer{}).(*thousandeyes.AgentServer)
}
