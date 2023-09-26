package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func dataSourceThousandeyesAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesAgentRead,

		Schema: map[string]*schema.Schema{
			"agent_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the agent.",
			},
			"agent_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique ID of the agent.",
			},
		},
		Description: "This data source allows you to define a ThousandEyes agent. For more information, see [Global Vantage Points](https://docs.thousandeyes.com/product-documentation/global-vantage-points).",
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
		if *agent.AgentName == searchName {
			found = &agent
			break
		}
	}

	if found == nil {
		return fmt.Errorf("unable to locate any agent with the name: %s", searchName)
	}
	log.Printf("[INFO] ## Found Agent agent_id: %d - name: %s", found.AgentID, *found.AgentName)

	d.SetId(fmt.Sprint(found.AgentID))
	err = d.Set("agent_name", found.AgentName)
	if err != nil {
		return err
	}
	err = d.Set("agent_id", found.AgentID)
	if err != nil {
		return err
	}

	return nil
}
