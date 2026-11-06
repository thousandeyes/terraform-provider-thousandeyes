package thousandeyes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func TestResourceTagAppearanceDefaults(t *testing.T) {
	for _, tc := range []struct {
		name  string
		color string
		icon  string
	}{
		{name: "omit both"},
		{name: "configure color", color: "#112233"},
		{name: "configure icon", icon: "STAR"},
		{name: "configure both", color: "#112233", icon: "STAR"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := map[string]interface{}{
				"key": "office", "value": "london", "object_type": "test", "access_type": "all",
			}
			if tc.color != "" {
				config["color"] = tc.color
			}
			if tc.icon != "" {
				config["icon"] = tc.icon
			}
			var remote map[string]interface{}
			apiClient := client.NewAPIClient(&client.Configuration{
				ServerURL: "https://api.thousandeyes.com/v7",
				Context:   context.WithValue(context.Background(), accountGroupIdKey, "12345"),
				HTTPClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					require.Equal(t, "12345", r.URL.Query().Get("aid"))
					if r.Method == http.MethodPost || r.Method == http.MethodPut {
						var request map[string]interface{}
						require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
						if r.Method == http.MethodPost {
							_, hasColor := request["color"]
							_, hasIcon := request["icon"]
							require.Equal(t, tc.color != "", hasColor)
							require.Equal(t, tc.icon != "", hasIcon)
						} else {
							require.Equal(t, "#445566", request["color"])
							require.Equal(t, "LABEL", request["icon"])
						}
						remote = request
						remote["id"] = "tag-123"
						remote["type"] = "static"
						if _, ok := remote["color"]; !ok {
							remote["color"] = "#A7EB10"
						}
						if _, ok := remote["icon"]; !ok {
							remote["icon"] = "LABEL"
						}
					} else {
						require.Equal(t, http.MethodGet, r.Method)
					}
					body, err := json.Marshal(remote)
					require.NoError(t, err)
					return jsonResponse(http.StatusOK, string(body)), nil
				})},
			})
			r := resourceTag()
			d := schema.TestResourceDataRaw(t, r.Schema, config)
			require.NoError(t, resourceTagCreate(d, apiClient))
			require.Equal(t, remote["color"], d.Get("color"))
			require.Equal(t, remote["icon"], d.Get("icon"))
			diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), nil)
			require.NoError(t, err)
			require.True(t, diff.Empty(), "API defaults must not cause a plan: %#v", diff)

			// Refresh and import must also adopt appearance returned by the API.
			imported := schema.TestResourceDataRaw(t, r.Schema, nil)
			imported.SetId(d.Id())
			states, err := r.Importer.StateContext(context.Background(), imported, apiClient)
			require.NoError(t, err)
			require.Len(t, states, 1)
			require.NoError(t, resourceTagRead(states[0], apiClient))
			diff, err = r.Diff(context.Background(), states[0].State(), terraform.NewResourceConfigRaw(config), nil)
			require.NoError(t, err)
			require.True(t, diff.Empty(), "imported appearance must not cause a plan: %#v", diff)

			// Users can override API defaults in a subsequent configuration.
			config["color"] = "#445566"
			config["icon"] = "LABEL"
			diff, err = r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), nil)
			require.NoError(t, err)
			require.False(t, diff.RequiresNew())
			updated, err := schema.InternalMap(r.Schema).Data(d.State(), diff)
			require.NoError(t, err)
			require.NoError(t, resourceTagUpdate(updated, apiClient))
			require.Equal(t, "#445566", updated.Get("color"))
			require.Equal(t, "LABEL", updated.Get("icon"))
			diff, err = r.Diff(context.Background(), updated.State(), terraform.NewResourceConfigRaw(config), nil)
			require.NoError(t, err)
			require.True(t, diff.Empty())
		})
	}
}
