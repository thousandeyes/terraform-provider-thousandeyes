package thousandeyes

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointagents"
)

func dataSourceThousandeyesEndpointAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEndpointAgentRead,
		Schema: map[string]*schema.Schema{
			"agent_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The exact name of the Endpoint Agent to look up.",
			},
			"agent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique ID of the Endpoint Agent.",
			},
		},
		Description: "Use this data source to look up an existing ThousandEyes Endpoint Agent by exact name.",
	}
}

func dataSourceEndpointAgentRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*endpointagents.EndpointAgentsAPIService)(&apiClient.Common)
	searchName := d.Get("agent_name").(string)

	log.Printf("[INFO] Reading ThousandEyes Endpoint Agent %q", searchName)

	var cursor string
	for {
		req := api.GetEndpointAgents().AgentName(searchName)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)
		if cursor != "" {
			req = req.Cursor(cursor)
		}

		resp, _, err := req.Execute()
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("received an empty response while looking up Endpoint Agent %q", searchName)
		}

		for _, agent := range resp.GetAgents() {
			if agent.GetName() != searchName {
				continue
			}

			agentID := agent.GetId()
			if agentID == "" {
				return fmt.Errorf("endpoint agent %q was returned without an ID", searchName)
			}

			d.SetId(agentID)
			if err := d.Set("agent_name", agent.GetName()); err != nil {
				return err
			}
			if err := d.Set("agent_id", agentID); err != nil {
				return err
			}
			return nil
		}

		links := resp.GetLinks()
		next, ok := links.GetNextOk()
		if !ok || next == nil || next.GetHref() == "" {
			break
		}

		nextURL, err := url.Parse(next.GetHref())
		if err != nil {
			return fmt.Errorf("unable to parse next Endpoint Agents page link: %w", err)
		}
		cursor = nextURL.Query().Get("cursor")
		if cursor == "" {
			return fmt.Errorf("next Endpoint Agents page link did not contain a cursor")
		}
	}

	return fmt.Errorf("unable to locate any Endpoint Agent with the name: %s", searchName)
}
