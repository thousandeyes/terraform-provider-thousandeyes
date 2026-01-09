package schemas

import (
	"context"
	"reflect"
	"testing"
)

func testLegacyTestStateUpgradeStateDataV0() map[string]any {
	return map[string]any{
		"agents":               []interface{}{map[string]interface{}{"agent_id": 1}},
		"alert_rules":          []interface{}{map[string]interface{}{"rule_id": 1}},
		"bgp_monitors":         []interface{}{map[string]interface{}{"monitor_id": 1}},
		"shared_with_accounts": []interface{}{map[string]interface{}{"aid": 1}},
	}
}

func testLegacyTestStateUpgradeStateDataV1() map[string]any {
	return map[string]any{
		"agents":               []interface{}{"1"},
		"alert_rules":          []interface{}{"1"},
		"bgp_monitors":         interface{}(nil),
		"monitors":             []interface{}{"1"},
		"shared_with_accounts": []interface{}{"1"},
	}
}

func TestLegacyTestStateUpgradeUpgradeV0(t *testing.T) {
	expected := testLegacyTestStateUpgradeStateDataV1()
	actual, err := LegacyTestStateUpgrade(context.Background(), testLegacyTestStateUpgradeStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}
