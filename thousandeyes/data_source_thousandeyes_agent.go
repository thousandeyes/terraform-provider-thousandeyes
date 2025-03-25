package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/agents"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique ID of the agent.",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The array of IP Addresses entries for the agent.",
			},
		},
		Description: "This data source allows you to define a ThousandEyes agent. For more information, see [Global Vantage Points](https://docs.thousandeyes.com/product-documentation/global-vantage-points).",
	}
}

func dataSourceThousandeyesAgentRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*agents.CloudAndEnterpriseAgentsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Reading Thousandeyes agent")

	searchName := d.Get("agent_name").(string)

	req := api.GetAgents()
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	var found *agents.AgentResponse

	for _, agent := range resp.GetAgents() {
		if *agent.AgentResponse.AgentName == searchName {
			found = agent.AgentResponse
			break
		}
	}

	if found == nil {
		return fmt.Errorf("unable to locate any agent with the name: %s", searchName)
	}
	log.Printf("[INFO] ## Found Agent agent_id: %d - name: %s", found.AgentId, *found.AgentName)

	d.SetId(*found.AgentId)
	err = d.Set("agent_name", found.AgentName)
	if err != nil {
		return err
	}
	err = d.Set("agent_id", found.AgentId)
	if err != nil {
		return err
	}
	err = d.Set("ip_addresses", found.IpAddresses)
	if err != nil {
		return err
	}

	return nil
}
