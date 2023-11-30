package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceAgentToAgent() *schema.Resource {
	agentToAgentSchemasOverride := map[string]*schema.Schema{
		"port": {
			Type:         schema.TypeInt,
			Description:  "The target port.",
			ValidateFunc: validation.IntBetween(1, 65535),
			Default:      49153,
			Optional:     true,
		},
	}

	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.AgentAgent{}, schemas, agentToAgentSchemasOverride),
		Create: resourceAgentAgentCreate,
		Read:   resourceAgentAgentRead,
		Update: resourceAgentAgentUpdate,
		Delete: resourceAgentAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create and configure an agent-to-agent test. This test type evaluates the performance of the underlying network between two physical sites. For more information about agent-to-agent tests, see [Agent-to-Agent Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#agent-to-agent-test).",
	}
	resource.Schema["protocol"] = schemas["protocol-agent_to_agent"]
	return &resource
}

func resourceAgentAgentRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(client *thousandeyes.Client, id int64) (interface{}, error) {
		return client.GetAgentAgent(id)
	})
}

func resourceAgentAgentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	update := ResourceUpdate(d, &thousandeyes.AgentAgent{}).(*thousandeyes.AgentAgent)
	_, err := client.UpdateAgentAgent(id, *update)
	if err != nil {
		return err
	}
	return resourceAgentAgentRead(d, m)
}

func resourceAgentAgentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	if err := client.DeleteAgentAgent(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAgentAgentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildAgentAgentStruct(d)
	remote, err := client.CreateAgentAgent(*local)
	if err != nil {
		return err
	}
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceAgentAgentRead(d, m)
}

func buildAgentAgentStruct(d *schema.ResourceData) *thousandeyes.AgentAgent {
	return ResourceBuildStruct(d, &thousandeyes.AgentAgent{}).(*thousandeyes.AgentAgent)
}
