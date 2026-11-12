package thousandeyes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointtests"
)

func TestEndpointAgentToServerDynamicLifecycle(t *testing.T) {
	for _, mode := range []string{"all-agents", "specific-agents", "agent-tags"} {
		t.Run(mode, func(t *testing.T) {
			config := map[string]interface{}{"test_name": "example", "application": "webex", "agent_selector_type": mode, "is_enabled": false}
			selector := map[string]interface{}{"agentSelectorType": mode, "maxMachines": 25}
			if mode == "specific-agents" {
				config["agents"] = []interface{}{"agent"}
				selector["agents"] = []string{"agent"}
			}
			if mode == "agent-tags" {
				config["tag_ids"] = []interface{}{"tag"}
				selector["tagIds"] = []string{"tag"}
			}
			var calls []string
			client := newEndpointAgentTestClient(endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				calls = append(calls, r.Method)
				if r.URL.Query().Get("aid") != "12345" {
					t.Errorf("missing aid")
				}
				path := "/v7/endpoint/tests/dynamic-tests/agent-to-server"
				if r.Method != "POST" {
					path += "/test"
				}
				if r.URL.Path != path {
					t.Errorf("unexpected path %s", r.URL.Path)
				}
				if r.Method == "POST" || r.Method == "PATCH" {
					var body map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Fatal(err)
					}
					if r.Method == "POST" {
						if body["agentSelectorType"] != mode || body["maxMachines"] != float64(25) {
							t.Errorf("bad selector: %#v", body)
						}
						for key, want := range map[string]interface{}{"agents": selector["agents"], "tagIds": selector["tagIds"]} {
							if fmt.Sprint(body[key]) != fmt.Sprint(want) {
								t.Errorf("%s: got %v, want %v", key, body[key], want)
							}
						}
						if body["application"] != "webex" {
							t.Error("missing application")
						}
					}
					if r.Method == "PATCH" {
						if body["isEnabled"] != false {
							t.Error("lost false value")
						}
						if _, ok := body["agentSelectorType"]; ok {
							t.Error("create-only field in update")
						}
					}
				}
				if r.Method == "DELETE" {
					return endpointAgentJSONResponse(204, ""), nil
				}
				b, _ := json.Marshal(map[string]interface{}{"testId": "test", "testName": "example", "application": "webex", "isEnabled": false, "agentSelectorConfig": selector})
				return endpointAgentJSONResponse(200, string(b)), nil
			}), "12345")
			d := schema.TestResourceDataRaw(t, resourceEndpointAgentToServerDynamic().Schema, config)
			if err := resourceEndpointAgentToServerDynamicCreate(d, client); err != nil {
				t.Fatal(err)
			}
			if d.Id() != "test" {
				t.Fatal("lost identity")
			}
			if fmt.Sprint(calls) != "[POST PATCH GET]" {
				t.Fatal(calls)
			}
			imported := schema.TestResourceDataRaw(t, resourceEndpointAgentToServerDynamic().Schema, nil)
			imported.SetId("test")
			if err := resourceEndpointAgentToServerDynamicRead(imported, client); err != nil {
				t.Fatal(err)
			}
			if imported.Get("agent_selector_type") != mode {
				t.Fatal("selector not restored")
			}
			if err := resourceEndpointAgentToServerDynamicDelete(d, client); err != nil {
				t.Fatal(err)
			}
			if d.Id() != "" {
				t.Fatal("ID retained")
			}
		})
	}
}

func TestEndpointAgentToServerDynamicErrors(t *testing.T) {
	for _, status := range []int{404, 500} {
		client := newEndpointAgentTestClient(endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return endpointAgentJSONResponse(status, `{"error":"failure"}`), nil
		}), "")
		d := schema.TestResourceDataRaw(t, resourceEndpointAgentToServerDynamic().Schema, nil)
		d.SetId("test")
		err := resourceEndpointAgentToServerDynamicRead(d, client)
		if status == 404 {
			if err != nil || d.Id() != "" {
				t.Fatal(err)
			}
		} else if err == nil || d.Id() == "" {
			t.Fatal("API error swallowed")
		}
	}
}

func TestEndpointAgentToServerDynamicSchemaAndSelector(t *testing.T) {
	r := resourceEndpointAgentToServerDynamic()
	for _, k := range []string{"agents", "tag_ids", "max_machines", "ip_version", "is_prioritized"} {
		if !r.Schema[k].ForceNew {
			t.Errorf("%s must replace", k)
		}
	}
	for _, k := range []string{"test_name", "application", "interval", "is_enabled", "protocol", "tcp_probe_mode"} {
		if r.Schema[k].ForceNew {
			t.Errorf("%s must update", k)
		}
	}
	for _, tc := range []struct {
		mode         string
		agents, tags int
		valid        bool
	}{{"all-agents", 0, 0, true}, {"specific-agents", 1, 0, true}, {"agent-tags", 0, 1, true}, {"all-agents", 1, 0, false}, {"specific-agents", 0, 0, false}, {"agent-tags", 0, 0, false}, {"agent-tags", 1, 1, false}} {
		if got := endpointTestSelectorValid(tc.mode, tc.agents, tc.tags); got != tc.valid {
			t.Errorf("%+v: %v", tc, got)
		}
	}
}

func TestEndpointAgentToServerDynamicPlanValidation(t *testing.T) {
	for _, tc := range []struct {
		mode         string
		agents, tags []interface{}
		valid        bool
	}{
		{mode: "all-agents", valid: true},
		{mode: "specific-agents", agents: []interface{}{"a"}, valid: true},
		{mode: "agent-tags", tags: []interface{}{"t"}, valid: true},
		{mode: "specific-agents", valid: false},
		{mode: "all-agents", agents: []interface{}{"a"}, valid: false},
		{mode: "agent-tags", agents: []interface{}{"a"}, tags: []interface{}{"t"}, valid: false},
	} {
		cfg := map[string]interface{}{"test_name": "test", "application": "webex", "agent_selector_type": tc.mode}
		if tc.agents != nil {
			cfg["agents"] = tc.agents
		}
		if tc.tags != nil {
			cfg["tag_ids"] = tc.tags
		}
		_, err := resourceEndpointAgentToServerDynamic().Diff(context.Background(), nil, terraform.NewResourceConfigRaw(cfg), nil)
		if (err == nil) != tc.valid {
			t.Errorf("%+v: %v", tc, err)
		}
	}
	r := resourceEndpointAgentToServerDynamic()
	base := map[string]interface{}{"test_name": "test", "application": "webex", "agent_selector_type": "all-agents"}
	d := schema.TestResourceDataRaw(t, r.Schema, base)
	d.SetId("test")
	for _, tc := range []struct {
		key     string
		value   interface{}
		replace bool
	}{
		{"test_name", "renamed", false}, {"application", "zoom", false}, {"max_machines", 10, true}, {"is_prioritized", true, true},
	} {
		cfg := map[string]interface{}{}
		for k, v := range base {
			cfg[k] = v
		}
		cfg[tc.key] = tc.value
		diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(cfg), nil)
		if err != nil {
			t.Fatal(err)
		}
		if diff.RequiresNew() != tc.replace {
			t.Errorf("%s replacement = %v", tc.key, diff.RequiresNew())
		}
	}
}

func TestEndpointAgentToServerDynamicSelectorWithUnknownIDs(t *testing.T) {
	for _, tc := range []struct{ mode, field string }{{"specific-agents", "agents"}, {"agent-tags", "tag_ids"}} {
		t.Run(tc.mode, func(t *testing.T) {
			r := resourceEndpointAgentToServerDynamic()
			values := map[string]cty.Value{}
			for name, ty := range r.CoreConfigSchema().ImpliedType().AttributeTypes() {
				values[name] = cty.NullVal(ty)
			}
			values["test_name"] = cty.StringVal("test")
			values["application"] = cty.StringVal("webex")
			values["agent_selector_type"] = cty.StringVal(tc.mode)
			values[tc.field] = cty.SetVal([]cty.Value{cty.UnknownVal(cty.String)})
			cfg := terraform.NewResourceConfigShimmed(cty.ObjectVal(values), r.CoreConfigSchema())
			if _, err := r.Diff(context.Background(), nil, cfg, nil); err != nil {
				t.Fatalf("selector with a computed ID should be deferred: %v", err)
			}
		})
	}
}

func TestEndpointAgentToServerDynamicRefresh(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceEndpointAgentToServerDynamic().Schema, map[string]interface{}{"agents": []interface{}{"old"}, "tag_ids": []interface{}{"old"}})
	var resp endpointtests.DynamicTest
	if err := json.Unmarshal([]byte(`{"testId":"test","application":"zoom","testName":"updated","isEnabled":false,"hasPing":false,"hasTraceroute":true,"networkMeasurements":false,"interval":60,"protocol":"icmp","tcpProbeMode":"auto","_links":{"self":{"href":"https://example.com/test"}},"labels":[],"agentSelectorConfig":{"agentSelectorType":"all-agents","maxMachines":25}}`), &resp); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := flattenEndpointAgentToServerDynamic(d, &resp); err != nil {
			t.Fatal(err)
		}
		for k, want := range map[string]interface{}{"application": "zoom", "test_name": "updated", "is_enabled": false, "has_ping": false, "has_traceroute": true, "network_measurements": false, "max_machines": 25, "link": "https://example.com/test"} {
			if got := d.Get(k); got != want {
				t.Errorf("%s = %v, want %v", k, got, want)
			}
		}
		if d.Get("agents").(*schema.Set).Len() != 0 || d.Get("tag_ids").(*schema.Set).Len() != 0 {
			t.Fatal("stale selector values")
		}
	}
}
func TestEndpointAgentToServerDynamicCreateMissingID(t *testing.T) {
	c := newEndpointAgentTestClient(endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) { return endpointAgentJSONResponse(201, `{}`), nil }), "")
	d := schema.TestResourceDataRaw(t, resourceEndpointAgentToServerDynamic().Schema, map[string]interface{}{"test_name": "minimal", "application": "webex"})
	if err := resourceEndpointAgentToServerDynamicCreate(d, c); err == nil || d.Id() != "" {
		t.Fatal("missing ID must fail creation")
	}
}
func TestEndpointAgentToServerDynamicMinimumAndReplacement(t *testing.T) {
	r := resourceEndpointAgentToServerDynamic()
	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{"test_name": "minimal", "application": "webex"})
	diff, err := r.Diff(context.Background(), nil, cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	data, err := schema.InternalMap(r.Schema).Data(nil, diff)
	if err != nil {
		t.Fatal(err)
	}
	data.SetId("test")
	state := data.State()
	if diff, err := r.Diff(context.Background(), state, cfg, nil); err != nil || (diff != nil && !diff.Empty()) {
		t.Fatalf("minimum config drift: %v %v", diff, err)
	}
	for _, tc := range []struct {
		key     string
		value   interface{}
		replace bool
	}{{"application", "zoom", false}, {"interval", 120, false}, {"is_enabled", false, false}, {"max_machines", 12, true}, {"is_prioritized", true, true}, {"ip_version", "V6_ONLY", true}} {
		c := map[string]interface{}{"test_name": "minimal", "application": "webex"}
		c[tc.key] = tc.value
		diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(c), nil)
		if err != nil {
			t.Fatal(err)
		}
		if diff.RequiresNew() != tc.replace {
			t.Errorf("%s replacement = %v", tc.key, diff.RequiresNew())
		}
	}
}

func TestEndpointAgentToServerDynamicMinimumLifecycle(t *testing.T) {
	var calls []string
	c := newEndpointAgentTestClient(endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls = append(calls, r.Method)
		if r.Method == http.MethodPost {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			for key, want := range map[string]interface{}{"testName": "minimal", "application": "webex", "interval": float64(60), "protocol": "icmp", "tcpProbeMode": "auto", "isPrioritized": false, "maxMachines": float64(25)} {
				if body[key] != want {
					t.Errorf("%s = %v, want %v", key, body[key], want)
				}
			}
			for _, key := range []string{"endpointAgentLabels", "agents", "tagIds", "ipVersion", "hasPathTraceInSession", "isEnabled"} {
				if _, exists := body[key]; exists {
					t.Errorf("unexpected input %s", key)
				}
			}
		}
		return endpointAgentJSONResponse(200, `{"testId":"test","testName":"minimal","application":"webex","agentSelectorConfig":{"agentSelectorType":"all-agents"}}`), nil
	}), "")
	d := schema.TestResourceDataRaw(t, resourceEndpointAgentToServerDynamic().Schema, map[string]interface{}{"test_name": "minimal", "application": "webex"})
	if err := resourceEndpointAgentToServerDynamicCreate(d, c); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(calls) != "[POST GET]" {
		t.Fatal(calls)
	}
	if d.Get("agents").(*schema.Set).Len() != 0 || d.Get("max_machines") != 25 {
		t.Fatal("unexpected default selector")
	}
}

func TestEndpointAgentToServerDynamicMutationErrors(t *testing.T) {
	for _, method := range []string{"POST", "PATCH", "DELETE"} {
		for _, status := range []int{404, 500} {
			t.Run(fmt.Sprintf("%s/%d", method, status), func(t *testing.T) {
				c := newEndpointAgentTestClient(endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
					return endpointAgentJSONResponse(status, `{"error":"failure"}`), nil
				}), "")
				d := schema.TestResourceDataRaw(t, resourceEndpointAgentToServerDynamic().Schema, map[string]interface{}{"test_name": "test", "application": "webex"})
				d.SetId("test")
				var err error
				switch method {
				case "POST":
					d.SetId("")
					err = resourceEndpointAgentToServerDynamicCreate(d, c)
				case "PATCH":
					err = resourceEndpointAgentToServerDynamicUpdate(d, c)
				case "DELETE":
					err = resourceEndpointAgentToServerDynamicDelete(d, c)
				}
				if method == "DELETE" && status == 404 {
					if err != nil || d.Id() != "" {
						t.Fatalf("delete 404: %v", err)
					}
				} else if err == nil {
					t.Fatal("API error swallowed")
				} else if method != "POST" && d.Id() != "test" {
					t.Fatal("lost ID after failed mutation")
				}
			})
		}
	}
}

func TestEndpointAgentToServerDynamicInvalidSettings(t *testing.T) {
	for key, value := range map[string]interface{}{"test_name": "", "application": "invalid", "max_machines": 0, "interval": 0, "protocol": "ast-autodetect", "ip_version": "unknown", "tcp_probe_mode": "invalid", "agent_selector_type": "agent-labels"} {
		t.Run(key, func(t *testing.T) {
			s := resourceEndpointAgentToServerDynamic().Schema[key]
			_, errors := s.ValidateFunc(value, key)
			if len(errors) == 0 {
				t.Fatalf("accepted %s=%v", key, value)
			}
		})
	}
}
