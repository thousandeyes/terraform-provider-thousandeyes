package thousandeyes

import "github.com/william20111/go-thousandeyes"

func expandAgents(v interface{}) thousandeyes.Agents {
	var agents thousandeyes.Agents

	for _, er := range v.([]interface{}) {
		rer := er.(map[string]interface{})
		agent := &thousandeyes.Agent{
			AgentId: rer["agent_id"].(int),
		}
		agents = append(agents, *agent)
	}

	return agents
}
