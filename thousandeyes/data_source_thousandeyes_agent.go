package thousandeyes

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
	"log"
)

func dataSourceThousandeyesAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesAgentRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agent_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceThousandeyesAgentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes agent")

	searchName := d.Get("name").(string)

	agents, err := client.GetAgents()
	if err != nil {
		return err
	}

	var found *thousandeyes.Agent

	for _, agent := range *agents {
		if agent.AgentName == searchName {
			found = &agent
			break
		}
	}

	if found == nil {
		return fmt.Errorf("unable to locate any agent with the name: %s", searchName)
	}

	d.SetId(string(found.AgentId))
	d.Set("name", found.AgentName)
	d.Set("agent_id", found.AgentId)

	return nil
}
