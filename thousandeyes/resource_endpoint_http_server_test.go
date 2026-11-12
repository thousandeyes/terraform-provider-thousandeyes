package thousandeyes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointtests"
)

func TestEndpointHTTPServerLifecycle(t *testing.T) {
	for _, mode := range []string{"all-agents", "specific-agents", "agent-tags"} {
		t.Run(mode, func(t *testing.T) {
			config := map[string]interface{}{"test_name": "example", "url": "https://example.com", "agent_selector_type": mode, "password": "secret", "is_enabled": false}
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
				path := "/v7/endpoint/tests/scheduled-tests/http-server"
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
						if body["password"] != "secret" {
							t.Error("missing password")
						}
					}
					if r.Method == "PATCH" {
						if body["isEnabled"] != false {
							t.Error("lost false value")
						}
						if _, ok := body["password"]; ok {
							t.Error("create-only field in update")
						}
					}
				}
				if r.Method == "DELETE" {
					return endpointAgentJSONResponse(204, ""), nil
				}
				b, _ := json.Marshal(map[string]interface{}{"type": "http-server", "testId": "test", "testName": "example", "url": "https://example.com", "isEnabled": false, "agentSelectorConfig": selector})
				return endpointAgentJSONResponse(200, string(b)), nil
			}), "12345")
			d := schema.TestResourceDataRaw(t, resourceEndpointHTTPServer().Schema, config)
			if err := resourceEndpointHTTPServerCreate(d, client); err != nil {
				t.Fatal(err)
			}
			if d.Id() != "test" || d.Get("password") != "secret" {
				t.Fatal("lost identity or password")
			}
			if fmt.Sprint(calls) != "[POST PATCH GET]" {
				t.Fatal(calls)
			}
			imported := schema.TestResourceDataRaw(t, resourceEndpointHTTPServer().Schema, nil)
			imported.SetId("test")
			if err := resourceEndpointHTTPServerRead(imported, client); err != nil {
				t.Fatal(err)
			}
			if imported.Get("agent_selector_type") != mode {
				t.Fatal("selector not restored")
			}
			if err := resourceEndpointHTTPServerDelete(d, client); err != nil {
				t.Fatal(err)
			}
			if d.Id() != "" {
				t.Fatal("ID retained")
			}
		})
	}
}
func TestEndpointHTTPServerErrors(t *testing.T) {
	for _, status := range []int{404, 500} {
		client := newEndpointAgentTestClient(endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return endpointAgentJSONResponse(status, `{"error":"failure"}`), nil
		}), "")
		d := schema.TestResourceDataRaw(t, resourceEndpointHTTPServer().Schema, nil)
		d.SetId("test")
		err := resourceEndpointHTTPServerRead(d, client)
		if status == 404 {
			if err != nil || d.Id() != "" {
				t.Fatal(err)
			}
		} else if err == nil || d.Id() == "" {
			t.Fatal("API error swallowed")
		}
	}
}
func TestEndpointHTTPServerSchemaAndSelector(t *testing.T) {
	r := resourceEndpointHTTPServer()
	for _, k := range []string{"agents", "tag_ids", "max_machines", "password", "auth_type"} {
		if !r.Schema[k].ForceNew {
			t.Errorf("%s must replace", k)
		}
	}
	for _, k := range []string{"test_name", "url", "interval", "is_enabled", "protocol", "tcp_probe_mode"} {
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
func TestEndpointHTTPServerFlattenClearsInactiveSelectors(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceEndpointHTTPServer().Schema, map[string]interface{}{"agents": []interface{}{"stale"}, "tag_ids": []interface{}{"stale"}})
	var resp endpointtests.EndpointHttpServerTest
	if err := json.Unmarshal([]byte(`{"type":"http-server","testName":"refreshed","httpTargetTime":1500,"server":"example.com:443","createdDate":"2026-11-09T10:00:00Z","verifyCertificate":false,"_links":{"self":{"href":"https://example.com/test"}},"labels":[{"labelId":"ignored"}],"agentSelectorConfig":{"agentSelectorType":"all-agents","maxMachines":12}}`), &resp); err != nil {
		t.Fatal(err)
	}
	if err := flattenEndpointHTTPServer(d, &resp); err != nil {
		t.Fatal(err)
	}
	if d.Get("target_response_time") != 1500 {
		t.Fatal("response time target not restored from httpTargetTime")
	}
	if d.Get("server") != "example.com" {
		t.Fatal("server hostname not restored")
	}
	if d.Get("link") != "https://example.com/test" || !d.Get("labels").(*schema.Set).Contains("ignored") {
		t.Fatal("response metadata not restored")
	}
	if d.Get("test_name") != "refreshed" || d.Get("created_date") != "2026-11-09T10:00:00Z" || d.Get("verify_certificate") != false {
		t.Fatal("ordinary response fields not restored")
	}
	if d.Get("agents").(*schema.Set).Len() != 0 || d.Get("tag_ids").(*schema.Set).Len() != 0 || d.Get("max_machines") != 12 {
		t.Fatal("stale selector state")
	}
}
func TestAccEndpointHTTPServer(t *testing.T) {
	name := os.Getenv("TE_ENDPOINT_AGENT_NAME")
	for _, mode := range []string{"all-agents", "specific-agents", "agent-tags"} {
		t.Run(mode, func(t *testing.T) {
			extra := ""
			if mode == "specific-agents" {
				extra = "agents = [data.thousandeyes_endpoint_agent.example.agent_id]"
			}
			if mode == "agent-tags" {
				extra = "tag_ids = [thousandeyes_tag.endpoint.id]"
			}
			tagConfig := ""
			if mode == "agent-tags" {
				tagConfig = fmt.Sprintf(`resource "thousandeyes_tag" "endpoint" {
 key = "endpoint-http-acceptance"
 value = %q
 access_type = "all"
 object_type = "endpoint-agent"
 type = "dynamic"
 match_type = "and"
 filters {
 key = "hostname"
 mode = "in"
 values = [%q]
 }
}`, name, name)
			}
			config := func(testName string) string {
				return tagConfig + "\n" + fmt.Sprintf(`data "thousandeyes_endpoint_agent" "example" { agent_name = %q }
resource "thousandeyes_endpoint_http_server" "example" {
 test_name = %q
 url = "https://example.com"
 agent_selector_type = %q
 %s
}`, name, testName, mode, extra)
			}
			resource.Test(t, resource.TestCase{PreCheck: func() {
				if name == "" {
					t.Fatal("TE_ENDPOINT_AGENT_NAME must identify an existing Endpoint Agent in the test account")
				}
				testAccPreCheck(t)
			}, ProviderFactories: providerFactories, CheckDestroy: func(s *terraform.State) error {
				return testAccCheckResourceDestroy([]ResourceType{{ResourceName: "thousandeyes_endpoint_http_server", GetResource: func(id string) (interface{}, error) {
					req := endpointHTTPServerAPI(testClient).GetHttpServerEndpointScheduledTest(id)
					returnValue, _, err := SetAidFromContext(testClient.GetConfig().Context, req).Execute()
					return returnValue, err
				}}}, s)
			}, Steps: []resource.TestStep{
				{Config: config("Endpoint HTTP acceptance"), Check: resource.TestCheckResourceAttr("thousandeyes_endpoint_http_server.example", "agent_selector_type", mode)},
				{Config: config("Endpoint HTTP updated")},
				{ResourceName: "thousandeyes_endpoint_http_server.example", ImportState: true, ImportStateVerify: true},
			}})
		})
	}
}

func TestEndpointHTTPServerPlanValidation(t *testing.T) {
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
		cfg := map[string]interface{}{"test_name": "test", "url": "https://example.com", "agent_selector_type": tc.mode}
		if tc.agents != nil {
			cfg["agents"] = tc.agents
		}
		if tc.tags != nil {
			cfg["tag_ids"] = tc.tags
		}
		_, err := resourceEndpointHTTPServer().Diff(context.Background(), nil, terraform.NewResourceConfigRaw(cfg), nil)
		if (err == nil) != tc.valid {
			t.Errorf("%+v: %v", tc, err)
		}
	}
	r := resourceEndpointHTTPServer()
	base := map[string]interface{}{"test_name": "test", "url": "https://example.com", "agent_selector_type": "all-agents"}
	d := schema.TestResourceDataRaw(t, r.Schema, base)
	d.SetId("test")
	for _, tc := range []struct {
		key     string
		value   interface{}
		replace bool
	}{
		{"test_name", "renamed", false}, {"url", "https://other.example.com", false}, {"max_machines", 10, true}, {"password", "secret", true},
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

func TestEndpointHTTPServerAuthenticatedURLReplacement(t *testing.T) {
	for _, authType := range []string{"none", "basic", "ntlm"} {
		for _, tc := range []struct {
			name    string
			url     cty.Value
			changed bool
		}{
			{"changed-url", cty.StringVal("https://other.example.com/health"), true},
			{"unchanged-url", cty.StringVal("https://example.com"), false},
			{"equivalent-url", cty.StringVal("example.com"), false},
			{"unknown-url", cty.UnknownVal(cty.String), true},
		} {
			t.Run(authType+"/"+tc.name, func(t *testing.T) {
				r := resourceEndpointHTTPServer()
				d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
					"test_name": "test", "url": "https://example.com", "auth_type": authType,
				})
				d.SetId("test")
				values := map[string]cty.Value{}
				for name, ty := range r.CoreConfigSchema().ImpliedType().AttributeTypes() {
					values[name] = cty.NullVal(ty)
				}
				values["test_name"] = cty.StringVal("renamed")
				values["url"] = tc.url
				values["auth_type"] = cty.StringVal(authType)
				cfg := terraform.NewResourceConfigShimmed(cty.ObjectVal(values), r.CoreConfigSchema())
				diff, err := r.Diff(context.Background(), d.State(), cfg, nil)
				if err != nil {
					t.Fatal(err)
				}
				wantReplacement := authType != "none" && tc.changed
				if diff.RequiresNew() != wantReplacement {
					t.Fatalf("replacement = %v, want %v", diff.RequiresNew(), wantReplacement)
				}
			})
		}
	}
}

func TestEndpointHTTPServerSelectorReadFixup(t *testing.T) {
	name := "agent_selector_config"
	value := map[string]interface{}{"agent_selector_type": "all-agents"}
	if _, err := FixReadValues(context.Background(), nil, value, &name); err != nil {
		t.Fatal(err)
	}
	if name != "" {
		t.Fatal("nested selector must be mapped separately")
	}
}

func TestEndpointHTTPServerSelectorWithUnknownIDs(t *testing.T) {
	for _, tc := range []struct{ mode, field string }{{"specific-agents", "agents"}, {"agent-tags", "tag_ids"}} {
		t.Run(tc.mode, func(t *testing.T) {
			r := resourceEndpointHTTPServer()
			values := map[string]cty.Value{}
			for name, ty := range r.CoreConfigSchema().ImpliedType().AttributeTypes() {
				values[name] = cty.NullVal(ty)
			}
			values["test_name"] = cty.StringVal("test")
			values["url"] = cty.StringVal("https://example.com")
			values["agent_selector_type"] = cty.StringVal(tc.mode)
			values[tc.field] = cty.SetVal([]cty.Value{cty.UnknownVal(cty.String)})
			cfg := terraform.NewResourceConfigShimmed(cty.ObjectVal(values), r.CoreConfigSchema())
			if _, err := r.Diff(context.Background(), nil, cfg, nil); err != nil {
				t.Fatalf("selector with a computed ID should be deferred: %v", err)
			}
		})
	}
}

func TestEndpointHTTPServerCreateUsesNetworkMeasurements(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceEndpointHTTPServer().Schema, map[string]interface{}{
		"test_name": "test", "url": "https://example.com", "network_measurements": false,
	})
	req := buildEndpointHTTPServerCreate(d)
	if req.HasPing != nil || req.HasTraceroute != nil {
		t.Fatal("deprecated probe flags must not be sent")
	}
	if req.NetworkMeasurements == nil || *req.NetworkMeasurements {
		t.Fatal("explicit false network_measurements must be preserved")
	}
}

func TestEndpointHTTPServerIPVersionProtocolValidation(t *testing.T) {
	for _, tc := range []struct {
		protocol, ipVersion string
		valid               bool
	}{
		{"icmp", "V6_ONLY", true}, {"icmp", "V6_PREFER", true}, {"tcp", "V4_ONLY", true}, {"tcp", "V6_ONLY", false},
	} {
		cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
			"test_name": "test", "url": "https://example.com", "protocol": tc.protocol, "ip_version": tc.ipVersion,
		})
		_, err := resourceEndpointHTTPServer().Diff(context.Background(), nil, cfg, nil)
		if (err == nil) != tc.valid {
			t.Errorf("protocol=%s ip_version=%s: %v", tc.protocol, tc.ipVersion, err)
		}
	}
}

func TestEndpointHTTPServerSchemelessURL(t *testing.T) {
	r := resourceEndpointHTTPServer()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"test_name": "test", "url": "example.com:8443/health"})
	if got := buildEndpointHTTPServerCreate(d).Url; got != "https://example.com:8443/health" {
		t.Fatal(got)
	}
	if got := buildEndpointHTTPServerUpdate(d).GetUrl(); got != "https://example.com:8443/health" {
		t.Fatal(got)
	}
	if !r.Schema["url"].DiffSuppressFunc("url", "https://example.com:8443/health", "example.com:8443/health", d) {
		t.Fatal("equivalent URL forms should not drift")
	}
	if r.Schema["url"].DiffSuppressFunc("url", "http://example.com", "example.com", d) {
		t.Fatal("HTTP and default HTTPS are different")
	}
}

func TestAccEndpointHTTPServerMinimum(t *testing.T) {
	const address = "thousandeyes_endpoint_http_server.minimal"
	const minimal = `resource "thousandeyes_endpoint_http_server" "minimal" {
 test_name = "Endpoint HTTP minimum acceptance"
 url = "example.com"
}`
	nondefault := strings.Replace(minimal, "\n}", `
 interval = 300
 http_time_limit = 7000
 target_response_time = 2000
 network_measurements = false
 verify_certificate = false
}`, 1)
	empty := strings.Replace(minimal, "\n}", `
 agents = []
 tag_ids = []
 username = ""
 password = ""
}`, 1)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckResourceDestroy([]ResourceType{{ResourceName: "thousandeyes_endpoint_http_server", GetResource: func(id string) (interface{}, error) {
				req := endpointHTTPServerAPI(testClient).GetHttpServerEndpointScheduledTest(id)
				resp, _, err := SetAidFromContext(testClient.GetConfig().Context, req).Execute()
				return resp, err
			}}}, s)
		},
		Steps: []resource.TestStep{
			{Config: minimal, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "agent_selector_type", "all-agents"),
				resource.TestCheckResourceAttr(address, "max_machines", "25"),
				resource.TestCheckResourceAttr(address, "interval", "60"),
			)},
			{ResourceName: address, ImportState: true, ImportStateVerify: true},
			{Config: nondefault},
			{Config: minimal, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "http_time_limit", "5000"),
				resource.TestCheckResourceAttr(address, "target_response_time", "1000"),
				resource.TestCheckResourceAttr(address, "network_measurements", "true"),
				resource.TestCheckResourceAttr(address, "verify_certificate", "true"),
			)},
			{Config: empty},
			{Config: minimal},
		},
	})
}

func TestEndpointHTTPServerReadPreservesIPv6Server(t *testing.T) {
	for _, tc := range []struct{ server, want string }{
		{"2001:db8::1", "2001:db8::1"},
		{"[2001:db8::1]:443", "2001:db8::1"},
	} {
		t.Run(tc.server, func(t *testing.T) {
			client := newEndpointAgentTestClient(endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				response := fmt.Sprintf(`{"type":"http-server","testId":"ipv6-test","server":%q,"url":"https://[2001:db8::1]","agentSelectorConfig":{"agentSelectorType":"all-agents"}}`, tc.server)
				return endpointAgentJSONResponse(200, response), nil
			}), "12345")
			d := schema.TestResourceDataRaw(t, resourceEndpointHTTPServer().Schema, nil)
			d.SetId("ipv6-test")
			// The first read populates imported state; the second refreshes it.
			for i := 0; i < 2; i++ {
				if err := resourceEndpointHTTPServerRead(d, client); err != nil {
					t.Fatal(err)
				}
				if got := d.Get("server"); got != tc.want {
					t.Fatalf("read %d: got %q, want %q", i, got, tc.want)
				}
			}
		})
	}
}
