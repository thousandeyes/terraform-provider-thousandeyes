package thousandeyes

import (
	"strconv"

	"github.com/william20111/go-thousandeyes"
)

func expandAgents(v interface{}) thousandeyes.Agents {
	var agents thousandeyes.Agents

	for _, er := range v.([]interface{}) {
		rer := er.(map[string]interface{})
		agent := &thousandeyes.Agent{
			AgentID: rer["agent_id"].(int),
		}
		agents = append(agents, *agent)
	}

	return agents
}

func expandAlertRules(v interface{}) thousandeyes.AlertRules {
	var alertRules thousandeyes.AlertRules

	for _, er := range v.([]interface{}) {
		rer := er.(map[string]interface{})
		alertRule := &thousandeyes.AlertRule{
			RuleID: rer["rule_id"].(int),
		}
		alertRules = append(alertRules, *alertRule)
	}

	return alertRules
}

func expandBGPMonitors(v interface{}) thousandeyes.BGPMonitors {
	var bgpMonitors thousandeyes.BGPMonitors

	for _, er := range v.([]interface{}) {
		rer := er.(map[string]interface{})
		bgpMonitor := &thousandeyes.BGPMonitor{
			MonitorID: rer["monitor_id"].(int),
		}
		bgpMonitors = append(bgpMonitors, *bgpMonitor)
	}

	return bgpMonitors
}

func expandDNSServers(v interface{}) []thousandeyes.Server {
	var dnsServers []thousandeyes.Server

	for _, er := range v.([]interface{}) {
		rer := er.(map[string]interface{})
		targetDNSServer := &thousandeyes.Server{
			ServerName: rer["server_name"].(string),
		}
		dnsServers = append(dnsServers, *targetDNSServer)
	}

	return dnsServers
}

func unpackSIPAuthData(i interface{}) thousandeyes.SIPAuthData {
	var m = i.(map[string]interface{})
	var sipAuthData = thousandeyes.SIPAuthData{}

	for k, v := range m {
		if k == "auth_user" {
			sipAuthData.AuthUser = v.(string)
		}
		if k == "password" {
			sipAuthData.Password = v.(string)
		}
		if k == "port" {
			port, err := strconv.Atoi(v.(string))
			if err == nil {
				sipAuthData.Port = port
			}
		}
		if k == "protocol" {
			sipAuthData.Protocol = v.(string)
		}
		if k == "sip_proxy" {
			sipAuthData.SipProxy = v.(string)
		}
		if k == "sip_registrar" {
			sipAuthData.SipRegistrar = v.(string)
		}
		if k == "user" {
			sipAuthData.User = v.(string)
		}
	}

	return sipAuthData
}
