package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
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
		Description: "This resource allows you to create and configure an agent-to-server test. This test type measures network performance as seen from ThousandEyes agent(s) towards a remote server. For more information about agent-to-server tests, see [Agent-to-Server Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#agent-to-agent-test).",
	}
	return &resource
}

func resourceAgentServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	remote, err := client.GetAgentServer(id)
	if err != nil {
		return err
	}
	err = ResourceRead(d, remote)
	if err != nil {
		return err
	}
	return nil
}

func resourceAgentServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
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
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
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
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceAgentServerRead(d, m)
}

func buildAgentServerStruct(d *schema.ResourceData) *thousandeyes.AgentServer {
	return ResourceBuildStruct(d, &thousandeyes.AgentServer{}).(*thousandeyes.AgentServer)
}
