package schemas

import (
	"context"
	"testing"
)

func TestLegacyTestStateUpgradePreservesBgpFields(t *testing.T) {
	rawState := map[string]any{
		"test_name":          "My HTTP Server Test",
		"interval":           300,
		"url":                "https://example.com",
		"bgp_measurements":   false,
		"use_public_bgp":     false,
		"mtu_measurements":   true,
		"num_path_traces":    3,
		"agents": []interface{}{
			map[string]interface{}{
				"agent_id":   1185,
				"agent_name": "São Paulo, Brazil",
			},
		},
	}

	result, err := LegacyTestStateUpgrade(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("LegacyTestStateUpgrade returned error: %v", err)
	}

	checks := []struct {
		field    string
		expected any
	}{
		{"bgp_measurements", false},
		{"use_public_bgp", false},
		{"mtu_measurements", true},
		{"num_path_traces", 3},
	}

	for _, c := range checks {
		val, ok := result[c.field]
		if !ok {
			t.Errorf("field %q was deleted during state upgrade; it should be preserved", c.field)
			continue
		}
		if val != c.expected {
			t.Errorf("field %q = %v, want %v", c.field, val, c.expected)
		}
	}
}

func TestLegacyTestStateUpgradeWorksWithoutBgpFields(t *testing.T) {
	rawState := map[string]any{
		"test_name": "Minimal Test",
		"interval":  300,
		"url":       "https://example.com",
		"agents": []interface{}{
			map[string]interface{}{
				"agent_id":   1185,
				"agent_name": "São Paulo, Brazil",
			},
		},
	}

	result, err := LegacyTestStateUpgrade(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("LegacyTestStateUpgrade returned error: %v", err)
	}

	for _, field := range []string{"bgp_measurements", "use_public_bgp", "mtu_measurements", "num_path_traces"} {
		if _, ok := result[field]; ok {
			t.Errorf("field %q should not appear in state when it was not present in v2 state", field)
		}
	}
}
