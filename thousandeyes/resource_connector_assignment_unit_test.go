package thousandeyes

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
)

func TestExpandOperationIDsPreservesEmptyIDs(t *testing.T) {
	d := schema.TestResourceDataRaw(t, schemas.ConnectorAssignmentSchema, map[string]interface{}{
		"connector_id": "connector-id",
		"operation_ids": []interface{}{
			"operation-2",
			"",
			"operation-1",
		},
	})

	got := expandOperationIDs(d)
	want := []string{"", "operation-1", "operation-2"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected operation_ids expansion: got=%v want=%v", got, want)
	}
}
