package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func dataSourceThousandeyesAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesAgentRead,

		Schema: map[string]*schema.Schema{
			"agent_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agent_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"agent_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceThousandeyesAgentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes agent")

	searchName := d.Get("agent_name").(string)

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
	log.Printf("[INFO] ## Found Agent agent_id: %d - name: %s", found.AgentID, found.AgentName)

	d.SetId(fmt.Sprint(found.AgentID))
	d.Set("agent_name", found.AgentName)
	d.Set("agent_id", found.AgentID)

	return nil
}
