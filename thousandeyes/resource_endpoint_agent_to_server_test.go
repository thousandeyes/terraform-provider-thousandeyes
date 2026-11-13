package thousandeyes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointtests"
)

const endpointAgentToServerPath = "/v7/endpoint/tests/scheduled-tests/agent-to-server"

func TestAccThousandEyesEndpointAgentToServer(t *testing.T) {
	selectorCases := []struct {
		name          string
		createConfig  string
		updateConfig  string
		selector      string
		selectorValue string
	}{
		{
			name:          "all_agents",
			createConfig:  "all_agents.tf",
			updateConfig:  "all_agents_update.tf",
			selector:      "agent_selector_type",
			selectorValue: "all-agents",
		},
		{
			name:          "specific_agents",
			createConfig:  "specific_agents.tf",
			updateConfig:  "specific_agents_update.tf",
			selector:      "agent_selector_type",
			selectorValue: "specific-agents",
		},
		{
			name:          "agent_tags",
			createConfig:  "agent_tags.tf",
			updateConfig:  "agent_tags_update.tf",
			selector:      "agent_selector_type",
			selectorValue: "agent-tags",
		},
	}

	for _, tc := range selectorCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceName := "thousandeyes_endpoint_agent_to_server.test"
			createChecks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "test_name", "Endpoint Agent-to-Server acceptance test"),
				resource.TestCheckResourceAttr(resourceName, "server", "www.thousandeyes.com"),
				resource.TestCheckResourceAttr(resourceName, tc.selector, tc.selectorValue),
			}
			if tc.selectorValue == endpointAgentSelectorAgentTags {
				createChecks = append(createChecks,
					resource.TestCheckResourceAttr("thousandeyes_tag.endpoint_agents", "object_type", "endpoint-agent"),
					resource.TestCheckResourceAttr("thousandeyes_tag.endpoint_agents", "type", "dynamic"),
					resource.TestCheckResourceAttr("thousandeyes_tag.endpoint_agents", "filters.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.0", "thousandeyes_tag.endpoint_agents", "id"),
				)
			}
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheck(t)
					if tc.selectorValue == endpointAgentSelectorSpecificAgents {
						endpointAgentID := os.Getenv("TE_ENDPOINT_AGENT_ID")
						if endpointAgentID == "" {
							t.Fatal("TE_ENDPOINT_AGENT_ID must be set for specific-agent acceptance coverage")
						}
						t.Setenv("TF_VAR_endpoint_agent_id", endpointAgentID)
					}
				},
				ProviderFactories: providerFactories,
				CheckDestroy:      testAccCheckEndpointAgentToServerDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccEndpointAgentToServerConfig(t, tc.createConfig),
						Check:  resource.ComposeTestCheckFunc(createChecks...),
					},
					{
						Config: testAccEndpointAgentToServerConfig(t, tc.updateConfig),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "test_name", "Endpoint Agent-to-Server acceptance test updated"),
							resource.TestCheckResourceAttr(resourceName, tc.selector, tc.selectorValue),
						),
					},
				},
			})
		})
	}
}

func testAccEndpointAgentToServerConfig(t *testing.T, filename string) string {
	t.Helper()
	path := filepath.Join("acceptance_resources", "endpoint_agent_to_server", filename)
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read acceptance config %s: %v", path, err)
	}
	return string(contents)
}

func testAccCheckEndpointAgentToServerDestroy(s *terraform.State) error {
	resources := []ResourceType{{
		ResourceName: "thousandeyes_endpoint_agent_to_server",
		GetResource: func(id string) (interface{}, error) {
			api := (*endpointtests.AgentToServerEndpointScheduledTestsAPIService)(&testClient.Common)
			req := api.GetAgentToServerEndpointScheduledTest(id)
			req = SetAidFromContext(testClient.GetConfig().Context, req)
			resp, _, err := req.Execute()
			return resp, err
		},
	}, {
		ResourceName: "thousandeyes_tag",
		GetResource: func(id string) (interface{}, error) {
			return getTag(id)
		},
	}}
	return testAccCheckResourceDestroy(resources, s)
}

func endpointAgentToServerResourceForTest(t *testing.T) *schema.Resource {
	t.Helper()

	resource, ok := Provider().ResourcesMap["thousandeyes_endpoint_agent_to_server"]
	if !ok || resource == nil {
		t.Fatal("thousandeyes_endpoint_agent_to_server is not registered")
	}
	return resource
}

func TestEndpointAgentToServerResourceSchemaSupportsOnlyCurrentSelectors(t *testing.T) {
	resource := endpointAgentToServerResourceForTest(t)

	for _, key := range []string{"agent_selector_type", "agents", "tag_ids", "max_machines"} {
		if _, ok := resource.Schema[key]; !ok {
			t.Errorf("expected schema field %q", key)
		}
	}
	for _, key := range []string{"agent_labels", "endpoint_agent_labels", "labels"} {
		if _, ok := resource.Schema[key]; ok {
			t.Errorf("schema must not expose excluded field %q", key)
		}
	}

	validateSelector := resource.Schema["agent_selector_type"].ValidateFunc
	if validateSelector == nil {
		t.Fatal("agent_selector_type must validate supported selector values")
	}
	for _, value := range []string{"all-agents", "specific-agents", "agent-tags"} {
		if warnings, errs := validateSelector(value, "agent_selector_type"); len(errs) != 0 || len(warnings) != 0 {
			t.Errorf("selector %q should be accepted, warnings=%v errors=%v", value, warnings, errs)
		}
	}
	for _, value := range []string{"agent-labels", "unknown"} {
		if _, errs := validateSelector(value, "agent_selector_type"); len(errs) == 0 {
			t.Errorf("selector %q should be rejected", value)
		}
	}

	for _, key := range []string{"agent_selector_type", "agents", "tag_ids", "max_machines", "ip_version", "is_prioritized"} {
		if !resource.Schema[key].ForceNew {
			t.Errorf("create-only field %q must require replacement", key)
		}
	}
}

func TestEndpointAgentToServerValidateSelectorInputs(t *testing.T) {
	stringSet := func(values []string) *schema.Set {
		items := make([]interface{}, 0, len(values))
		for _, value := range values {
			items = append(items, value)
		}
		return schema.NewSet(schema.HashString, items)
	}
	cases := []struct {
		name     string
		selector string
		agents   []string
		tags     []string
		wantErr  bool
	}{
		{name: "all agents", selector: "all-agents"},
		{name: "specific agents", selector: "specific-agents", agents: []string{"agent-1"}},
		{name: "agent tags", selector: "agent-tags", tags: []string{"tag-1"}},
		{name: "all agents rejects agent IDs", selector: "all-agents", agents: []string{"agent-1"}, wantErr: true},
		{name: "all agents rejects tag IDs", selector: "all-agents", tags: []string{"tag-1"}, wantErr: true},
		{name: "specific agents requires IDs", selector: "specific-agents", wantErr: true},
		{name: "specific agents rejects tag IDs", selector: "specific-agents", agents: []string{"agent-1"}, tags: []string{"tag-1"}, wantErr: true},
		{name: "agent tags requires IDs", selector: "agent-tags", wantErr: true},
		{name: "agent tags rejects agent IDs", selector: "agent-tags", agents: []string{"agent-1"}, tags: []string{"tag-1"}, wantErr: true},
		{name: "unknown selector", selector: "agent-labels", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateEndpointAgentSelector(tc.selector, stringSet(tc.agents), stringSet(tc.tags))
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error=%t, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestEndpointAgentToServerCreateSelectorPayloadAndState(t *testing.T) {
	cases := []struct {
		name           string
		selector       string
		config         map[string]interface{}
		responseConfig string
		wantAgents     []string
		wantTagIDs     []string
		wantMax        int
	}{
		{
			name:     "all agents omits selector lists and uses API max default",
			selector: "all-agents",
			config: map[string]interface{}{
				"agent_selector_type": "all-agents",
			},
			responseConfig: `{"agentSelectorType":"all-agents","maxMachines":25}`,
			wantMax:        25,
		},
		{
			name:     "specific agents sends and restores agent IDs",
			selector: "specific-agents",
			config: map[string]interface{}{
				"agent_selector_type": "specific-agents",
				"agents":              []interface{}{"agent-1", "agent-2"},
				"max_machines":        4,
			},
			responseConfig: `{"agentSelectorType":"specific-agents","agents":["agent-1","agent-2"],"maxMachines":4}`,
			wantAgents:     []string{"agent-1", "agent-2"},
			wantMax:        4,
		},
		{
			name:     "agent tags sends and restores tag IDs",
			selector: "agent-tags",
			config: map[string]interface{}{
				"agent_selector_type": "agent-tags",
				"tag_ids":             []interface{}{"tag-1"},
				"max_machines":        2,
			},
			responseConfig: `{"agentSelectorType":"agent-tags","tagIds":["tag-1"],"maxMachines":2}`,
			wantTagIDs:     []string{"tag-1"},
			wantMax:        2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resource := endpointAgentToServerResourceForTest(t)
			config := map[string]interface{}{
				"test_name": "scheduled endpoint test",
				"server":    "example.com",
				"interval":  60,
				"protocol":  "icmp",
			}
			for key, value := range tc.config {
				config[key] = value
			}
			d := schema.TestResourceDataRaw(t, resource.Schema, config)

			transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				if got := r.URL.Path; got != endpointAgentToServerPath {
					if got != endpointAgentToServerPath+"/test-123" {
						return endpointAgentJSONResponse(http.StatusNotFound, `{"error":"unexpected path"}`), nil
					}
				}
				if got := r.URL.Query().Get("aid"); got != "aid-123" {
					t.Errorf("expected aid=aid-123, got %q", got)
				}

				switch r.Method {
				case http.MethodPost:
					if r.URL.Path != endpointAgentToServerPath {
						return endpointAgentJSONResponse(http.StatusNotFound, `{}`), nil
					}
					var payload map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						t.Errorf("decode create payload: %v", err)
					}
					if got := payload["agentSelectorType"]; got != tc.selector {
						t.Errorf("expected agentSelectorType %q, got %#v", tc.selector, got)
					}
					if got := payload["agents"]; !sameJSONSlice(got, tc.wantAgents) {
						t.Errorf("expected agents %v, got %#v", tc.wantAgents, got)
					}
					if got := payload["tagIds"]; !sameJSONSlice(got, tc.wantTagIDs) {
						t.Errorf("expected tagIds %v, got %#v", tc.wantTagIDs, got)
					}
					if tc.config["max_machines"] == nil {
						if _, ok := payload["maxMachines"]; ok {
							t.Errorf("omitted max_machines should use the API default, payload included maxMachines=%v", payload["maxMachines"])
						}
					} else if got := payload["maxMachines"]; got != float64(tc.wantMax) {
						t.Errorf("expected maxMachines %d, got %#v", tc.wantMax, got)
					}
					if _, ok := payload["endpointAgentLabels"]; ok {
						t.Errorf("create payload must not contain deprecated endpointAgentLabels")
					}
					if _, ok := payload["labels"]; ok {
						t.Errorf("create payload must not contain excluded Endpoint Labels")
					}
					return endpointAgentJSONResponse(http.StatusCreated, `{"testId":"test-123","type":"agent-to-server"}`), nil
				case http.MethodGet:
					if r.URL.Path != endpointAgentToServerPath+"/test-123" {
						return endpointAgentJSONResponse(http.StatusNotFound, `{}`), nil
					}
					body := fmt.Sprintf(`{"testId":"test-123","type":"agent-to-server","testName":"scheduled endpoint test","server":"example.com","interval":60,"protocol":"icmp","agentSelectorConfig":%s}`, tc.responseConfig)
					return endpointAgentJSONResponse(http.StatusOK, body), nil
				default:
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
					return endpointAgentJSONResponse(http.StatusMethodNotAllowed, `{}`), nil
				}
			})

			if err := resource.Create(d, newEndpointAgentTestClient(transport, "aid-123")); err != nil {
				t.Fatalf("create failed: %v", err)
			}
			if got := d.Id(); got != "test-123" {
				t.Fatalf("expected resource ID test-123, got %q", got)
			}
			if got := d.Get("agent_selector_type"); got != tc.selector {
				t.Errorf("expected selector state %q, got %#v", tc.selector, got)
			}
			if got := resourceTestStringSet(d, "agents"); !equalStrings(got, tc.wantAgents) {
				t.Errorf("expected agents state %v, got %v", tc.wantAgents, got)
			}
			if got := resourceTestStringSet(d, "tag_ids"); !equalStrings(got, tc.wantTagIDs) {
				t.Errorf("expected tag_ids state %v, got %v", tc.wantTagIDs, got)
			}
			if got := d.Get("max_machines"); got != tc.wantMax {
				t.Errorf("expected max_machines state %d, got %#v", tc.wantMax, got)
			}
		})
	}
}

func TestEndpointAgentToServerReadAfterImportRestoresSelectorState(t *testing.T) {
	resource := endpointAgentToServerResourceForTest(t)
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet || r.URL.Path != endpointAgentToServerPath+"/imported-test" {
			t.Errorf("expected GET for imported scheduled test, got %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("aid"); got != "aid-123" {
			t.Errorf("expected aid=aid-123, got %q", got)
		}
		return endpointAgentJSONResponse(http.StatusOK, `{"testId":"imported-test","type":"agent-to-server","testName":"imported","server":"2001:db8::1","_links":{"self":{"href":"https://api.example.com/tests/imported-test"}},"labels":[{"labelId":"legacy-label"}],"agentSelectorConfig":{"agentSelectorType":"agent-tags","tagIds":["tag-7"],"maxMachines":25}}`), nil
	})

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	imported, err := resource.Importer.StateContext(context.Background(), d, nil)
	if err != nil {
		t.Fatalf("import state failed: %v", err)
	}
	if len(imported) != 1 {
		t.Fatalf("expected one imported state, got %d", len(imported))
	}
	imported[0].SetId("imported-test")
	if err := resource.Read(imported[0], newEndpointAgentTestClient(transport, "aid-123")); err != nil {
		t.Fatalf("read after import failed: %v", err)
	}
	if got := imported[0].Get("agent_selector_type"); got != "agent-tags" {
		t.Errorf("expected imported selector agent-tags, got %#v", got)
	}
	if got := imported[0].Get("server"); got != "2001:db8::1" {
		t.Errorf("expected IPv6 server to be preserved, got %#v", got)
	}
	if got := imported[0].Get("link"); got != "https://api.example.com/tests/imported-test" {
		t.Errorf("expected link from SDK response, got %#v", got)
	}
	if got := resourceTestStringSet(imported[0], "tag_ids"); !equalStrings(got, []string{"tag-7"}) {
		t.Errorf("expected imported tag_ids [tag-7], got %v", got)
	}
	if got := imported[0].Get("max_machines"); got != 25 {
		t.Errorf("expected imported max_machines 25, got %#v", got)
	}
}

func TestEndpointAgentToServerCreateAppliesUpdateOnlyFields(t *testing.T) {
	cases := []struct {
		name          string
		config        map[string]interface{}
		patchField    string
		wantPatch     interface{}
		responseField string
		responseValue interface{}
	}{
		{
			name: "disabled on create",
			config: map[string]interface{}{
				"agent_selector_type": "all-agents",
				"test_name":           "scheduled endpoint test",
				"server":              "example.com",
				"is_enabled":          false,
			},
			patchField:    "isEnabled",
			wantPatch:     false,
			responseField: "isEnabled",
			responseValue: false,
		},
		{
			name: "TCP probe mode on create",
			config: map[string]interface{}{
				"agent_selector_type": "all-agents",
				"test_name":           "scheduled endpoint test",
				"server":              "example.com",
				"protocol":            "tcp",
				"tcp_probe_mode":      "sack",
			},
			patchField:    "tcpProbeMode",
			wantPatch:     "sack",
			responseField: "tcpProbeMode",
			responseValue: "sack",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			patches := 0
			transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				if got := r.URL.Query().Get("aid"); got != "aid-123" {
					t.Errorf("expected aid=aid-123, got %q", got)
				}
				switch r.Method {
				case http.MethodPost:
					if r.URL.Path != endpointAgentToServerPath {
						t.Errorf("expected create path %s, got %s", endpointAgentToServerPath, r.URL.Path)
					}
					var payload map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						t.Errorf("decode create payload: %v", err)
					}
					for _, field := range []string{"isEnabled", "tcpProbeMode"} {
						if _, ok := payload[field]; ok {
							t.Errorf("update-only field %q must not be sent in create payload", field)
						}
					}
					return endpointAgentJSONResponse(http.StatusCreated, `{"testId":"test-123","type":"agent-to-server"}`), nil
				case http.MethodPatch:
					if r.URL.Path != endpointAgentToServerPath+"/test-123" {
						t.Errorf("expected update path for created test, got %s", r.URL.Path)
					}
					patches++
					var payload map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						t.Errorf("decode update payload: %v", err)
					}
					if len(payload) != 1 || payload[tc.patchField] != tc.wantPatch {
						t.Errorf("expected update payload {%s: %#v}, got %#v", tc.patchField, tc.wantPatch, payload)
					}
					return endpointAgentJSONResponse(http.StatusOK, `{"testId":"test-123","type":"agent-to-server"}`), nil
				case http.MethodGet:
					if r.URL.Path != endpointAgentToServerPath+"/test-123" {
						t.Errorf("expected read path for created test, got %s", r.URL.Path)
					}
					response := map[string]interface{}{
						"testId":              "test-123",
						"type":                "agent-to-server",
						"testName":            "scheduled endpoint test",
						"server":              "example.com",
						"agentSelectorConfig": map[string]interface{}{"agentSelectorType": "all-agents", "maxMachines": 25},
						tc.responseField:      tc.responseValue,
					}
					body, err := json.Marshal(response)
					if err != nil {
						t.Fatalf("encode read response: %v", err)
					}
					return endpointAgentJSONResponse(http.StatusOK, string(body)), nil
				default:
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
					return endpointAgentJSONResponse(http.StatusMethodNotAllowed, `{}`), nil
				}
			})

			d := schema.TestResourceDataRaw(t, endpointAgentToServerResourceForTest(t).Schema, tc.config)
			if err := endpointAgentToServerResourceForTest(t).Create(d, newEndpointAgentTestClient(transport, "aid-123")); err != nil {
				t.Fatalf("create failed: %v", err)
			}
			if patches != 1 {
				t.Fatalf("expected one post-create update, got %d", patches)
			}
			if tc.responseField == "isEnabled" {
				if got := d.Get("is_enabled"); got != tc.responseValue {
					t.Errorf("expected refreshed is_enabled %v, got %#v", tc.responseValue, got)
				}
			} else if got := d.Get("tcp_probe_mode"); got != tc.responseValue {
				t.Errorf("expected refreshed tcp_probe_mode %v, got %#v", tc.responseValue, got)
			}
		})
	}
}

func TestEndpointAgentToServerUpdateSendsOnlyMutableFields(t *testing.T) {
	resource := endpointAgentToServerResourceForTest(t)
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method == http.MethodPatch && r.URL.Path == endpointAgentToServerPath+"/test-123" {
			var payload map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode update payload: %v", err)
			}
			if got := payload["testName"]; got != "renamed" {
				t.Errorf("expected update testName renamed, got %#v", got)
			}
			for _, key := range []string{"agentSelectorType", "agents", "tagIds", "maxMachines", "ipVersion", "isPrioritized"} {
				if _, ok := payload[key]; ok {
					t.Errorf("create-only field %q must not be sent in update payload", key)
				}
			}
			return endpointAgentJSONResponse(http.StatusOK, `{"testId":"test-123","type":"agent-to-server","testName":"renamed","server":"example.com","agentSelectorConfig":{"agentSelectorType":"specific-agents","agents":["agent-1"],"maxMachines":25}}`), nil
		}
		if r.Method == http.MethodGet && r.URL.Path == endpointAgentToServerPath+"/test-123" {
			return endpointAgentJSONResponse(http.StatusOK, `{"testId":"test-123","type":"agent-to-server","testName":"renamed","server":"example.com","agentSelectorConfig":{"agentSelectorType":"specific-agents","agents":["agent-1"],"maxMachines":25}}`), nil
		}
		t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		return endpointAgentJSONResponse(http.StatusMethodNotAllowed, `{}`), nil
	})

	d := resource.Data(&terraform.InstanceState{
		ID: "test-123",
		Attributes: map[string]string{
			"test_name":           "old-name",
			"server":              "example.com",
			"interval":            "60",
			"protocol":            "icmp",
			"agent_selector_type": "specific-agents",
			"agents.#":            "1",
			"agents.0":            "agent-1",
			"max_machines":        "25",
		},
	})
	if err := d.Set("test_name", "renamed"); err != nil {
		t.Fatalf("set updated test_name: %v", err)
	}
	if err := resource.Update(d, newEndpointAgentTestClient(transport, "aid-123")); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if got := d.Get("test_name"); got != "renamed" {
		t.Errorf("expected updated test_name state, got %#v", got)
	}
}

func TestEndpointAgentToServerDeleteClearsState(t *testing.T) {
	resource := endpointAgentToServerResourceForTest(t)
	transport := endpointAgentRoundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodDelete || r.URL.Path != endpointAgentToServerPath+"/test-123" {
			t.Errorf("expected DELETE for scheduled test, got %s %s", r.Method, r.URL.Path)
		}
		return endpointAgentJSONResponse(http.StatusNoContent, ""), nil
	})

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("test-123")
	if err := resource.Delete(d, newEndpointAgentTestClient(transport, "aid-123")); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if d.Id() != "" {
		t.Errorf("expected delete to clear state ID, got %q", d.Id())
	}
}

func sameJSONSlice(got interface{}, want []string) bool {
	values, ok := got.([]interface{})
	if !ok {
		return len(want) == 0 && got == nil
	}
	if len(values) != len(want) {
		return false
	}
	wantValues := make(map[string]int, len(want))
	for _, value := range want {
		wantValues[value]++
	}
	for _, value := range values {
		stringValue, ok := value.(string)
		if !ok || wantValues[stringValue] == 0 {
			return false
		}
		wantValues[stringValue]--
	}
	return true
}

func resourceTestStringSet(d *schema.ResourceData, key string) []string {
	set, ok := d.Get(key).(*schema.Set)
	if !ok || set == nil {
		return nil
	}
	values := make([]string, 0, set.Len())
	for _, value := range set.List() {
		values = append(values, fmt.Sprint(value))
	}
	return values
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	wantValues := make(map[string]int, len(want))
	for _, value := range want {
		wantValues[value]++
	}
	for _, value := range got {
		if wantValues[value] == 0 {
			return false
		}
		wantValues[value]--
	}
	return true
}
