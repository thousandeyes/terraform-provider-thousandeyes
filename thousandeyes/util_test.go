package thousandeyes

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func getReferenceData(schemaData map[string]*schema.Schema, attrs map[string]string) *schema.ResourceData {
	referenceResource := schema.Resource{
		Schema: schemaData,
	}
	referenceState := terraform.InstanceState{
		ID:         "testState",
		Attributes: attrs,
	}

	return referenceResource.Data(&referenceState)
}

func TestResourceBuildStruct(t *testing.T) {
	prefix := "8.19.2.2/19"
	cmpStruct := thousandeyes.BGP{
		Prefix: thousandeyes.String(prefix),
	}
	newStruct := thousandeyes.BGP{}
	attrs := map[string]string{
		"prefix": "8.19.2.2/19",
	}
	d := getReferenceData(schemas, attrs)
	ResourceBuildStruct(d, &newStruct)

	if *newStruct.Prefix != *cmpStruct.Prefix {
		t.Error("Building resource did not assign struct field correctly.")
	}

}

func TestResourceRead(t *testing.T) {
	prefix := "8.19.2.2/19"
	attrs := map[string]string{}
	d := getReferenceData(schemas, attrs)
	remoteResource := thousandeyes.BGP{
		Prefix: thousandeyes.String(prefix),
	}
	err := ResourceRead(d, &remoteResource)
	if err != nil {
		t.Errorf("Setting resource data returned error: %+v", err.Error())
	}
	if d.Get("prefix") != *remoteResource.Prefix {
		t.Errorf("Reading resource did not assign resource data correctly.\nStruct is %+v\nResource is %+v", remoteResource, d.State().Attributes)
	}
}

func TestResourceReadValue(t *testing.T) {
	testStruct := thousandeyes.AccountGroupRole{
		RoleName:                 thousandeyes.String("TestRole"),
		RoleID:                   thousandeyes.Int64(1),
		HasManagementPermissions: thousandeyes.Bool(true),
		Builtin:                  thousandeyes.Bool(false),
		Permissions: &[]thousandeyes.Permission{
			{
				IsManagementPermission: thousandeyes.Bool(false),
				Label:                  thousandeyes.String("foo"),
				PermissionID:           thousandeyes.Int64(27),
			},
			{
				IsManagementPermission: thousandeyes.Bool(true),
				Label:                  thousandeyes.String("bar"),
				PermissionID:           thousandeyes.Int64(28),
			},
		},
	}
	resultMap := map[string]interface{}{
		"role_name":                  thousandeyes.String("TestRole"),
		"role_id":                    thousandeyes.Int64(1),
		"has_management_permissions": thousandeyes.Bool(true),
		"builtin":                    thousandeyes.Bool(false),
		"permissions": []map[string]interface{}{
			{
				"is_management_permission": thousandeyes.Bool(false),
				"label":                    thousandeyes.String("foo"),
				"permission_id":            thousandeyes.Int64(27),
			},
			{
				"is_management_permission": thousandeyes.Bool(true),
				"label":                    thousandeyes.String("bar"),
				"permission_id":            thousandeyes.Int64(28),
			},
		},
	}

	result, err := ReadValue(&testStruct)
	if err != nil {
		t.Errorf("Error running ReadValue: %+v", err)
	}

	// We need to inform Go of the type of the "permissions" field in the result,
	// in order for reflect.DeepEqual to identify them as equivalent.
	// Since it's a list of maps we can't do that in one step.
	typedResult := result.(map[string]interface{})
	permissions := typedResult["permissions"].([]interface{})
	var typedPermissions []map[string]interface{}
	for _, v := range permissions {
		typedPermissions = append(typedPermissions, v.(map[string]interface{}))
	}
	typedResult["permissions"] = typedPermissions

	if reflect.DeepEqual(typedResult, resultMap) != true {
		t.Errorf("Struct processing and test map do not match: \n\n%+v\n\n%+v\n", result, resultMap)
	}

}

func TestFixReadValues(t *testing.T) {
	var output interface{}
	var err error

	// non-map, non-list
	normalInput := 4
	normalTarget := 4
	output, err = FixReadValues(normalInput, "normal")
	if err != nil {
		t.Errorf("normal input returned error: %s", err.Error())
	}
	if output.(int) != 4 {
		t.Errorf("Returned wrong value for int input. Received  %#v, expected %#v", output, normalTarget)
	}

	// agents
	agentsInput := []interface{}{
		map[string]interface{}{
			"agent_name": "foo",
			"agent_id":   1,
		},
		map[string]interface{}{
			"agent_name": "bar",
			"agent_id":   2,
		},
	}
	agentsTarget := []interface{}{
		map[string]interface{}{
			"agent_id": 1,
		},
		map[string]interface{}{
			"agent_id": 2,
		},
	}
	output, err = FixReadValues(agentsInput, "agents")
	if err != nil {
		t.Errorf("agents input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, agentsTarget) != true {
		t.Errorf("Values not stripped correctly from agents input: Received %#v Expected %#v", output, agentsTarget)
	}

	// alert_rules
	alertRulesInput := []interface{}{
		map[string]interface{}{
			"rule_name": thousandeyes.String("foo"),
			"rule_id":   thousandeyes.Int(1),
			"default":   thousandeyes.Bool(false),
		},
		map[string]interface{}{
			"rule_name": thousandeyes.String("bar"),
			"rule_id":   thousandeyes.Int(2),
			"default":   thousandeyes.Bool(false),
		},
		map[string]interface{}{
			"rule_name": thousandeyes.String("bar"),
			"rule_id":   thousandeyes.Int(3),
			"default":   thousandeyes.Bool(true),
		},
	}
	alertRulesTarget := []interface{}{
		map[string]interface{}{
			"rule_id": thousandeyes.Int(1),
		},
		map[string]interface{}{
			"rule_id": thousandeyes.Int(2),
		},
		map[string]interface{}{
			"rule_id": thousandeyes.Int(3),
		},
	}
	output, err = FixReadValues(alertRulesInput, "alert_rules")
	if err != nil {
		t.Errorf("alert_rules input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, alertRulesTarget) != true {
		t.Errorf("Values not stripped correctly from alert_rules input: Received %#v Expected %#v", output, alertRulesTarget)
	}

	// bgp_monitors
	monitorsInput := []interface{}{
		map[string]interface{}{
			"monitor_name": thousandeyes.String("foo"),
			"monitor_id":   thousandeyes.Int(1),
			"monitor_type": thousandeyes.String("Public"),
		},
		map[string]interface{}{
			"monitor_name": thousandeyes.String("bar"),
			"monitor_id":   thousandeyes.Int(2),
			"monitor_type": thousandeyes.String("Private"),
		},
	}
	monitorsTarget := []interface{}{
		map[string]interface{}{
			"monitor_id": thousandeyes.Int(2),
		},
	}
	output, err = FixReadValues(monitorsInput, "bgp_monitors")
	if err != nil {
		t.Errorf("bgp_monitors input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, monitorsTarget) != true {
		t.Errorf("Values not stripped correctly from bgp_monitors input: Received %#v Expected %#v", output, monitorsTarget)
	}

	// groups
	groupsInput := []interface{}{
		map[string]interface{}{
			"group_name": "foo",
			"group_id":   1,
		},
		map[string]interface{}{
			"group_name": "bar",
			"group_id":   2,
		},
	}
	groupsTarget := []interface{}{
		map[string]interface{}{
			"group_id": 1,
		},
		map[string]interface{}{
			"group_id": 2,
		},
	}
	output, err = FixReadValues(groupsInput, "groups")
	if err != nil {
		t.Errorf("groups input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, groupsTarget) != true {
		t.Errorf("Values not stripped correctly from groups input: Received %#v Expected %#v", output, groupsTarget)
	}

	//	shared_with_accounts
	account_group_id = 2
	accountsInput := []interface{}{
		map[string]interface{}{
			"name": thousandeyes.String("foo"),
			"aid":  thousandeyes.Int64(1),
		},
		map[string]interface{}{
			"name": thousandeyes.String("bar"),
			"aid":  thousandeyes.Int64(2),
		},
	}
	accountsTarget := []interface{}{
		map[string]interface{}{
			"aid": thousandeyes.Int64(1),
		},
	}
	output, err = FixReadValues(accountsInput, "shared_with_accounts")
	if err != nil {
		t.Errorf("shared_with_accounts input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, accountsTarget) != true {
		t.Errorf("Values not stripped correctly from shared_with_accounts input: Received %#v Expected %#v", output, accountsTarget)
	}
	//  We should fail if account_group_id isn't set and the list of account groups is > 1
	account_group_id = 0
	output, err = FixReadValues(accountsInput, "shared_with_accounts")
	if err == nil {
		t.Errorf("Error was not returned when shared_with_accounts length was > 1 and account_group_id  was not set")
	}
	// We should not fail if account_group_id isn't set and the list of account groups is < 2
	accountsInput = []interface{}{
		map[string]interface{}{
			"name": "bar",
			"aid":  2,
		},
	}
	output, err = FixReadValues(accountsInput, "shared_with_accounts")
	if err != nil {
		t.Errorf("shared_with_accounts input returned error when shared_with_accounts wasn't set despite list of account groups being < 2: %s", err.Error())
	}
	if output != nil {
		t.Errorf("Values not stripped correctly from shared_with_accounts input: Received %#v Expected %#v", output, accountsTarget)
	}

	// target_sip_credentials
	sipCredsInput := map[string]interface{}{
		"sip_proxy": "foo.com",
	}
	sipCredsTarget := []interface{}{
		map[string]interface{}{
			"sip_proxy": "foo.com",
		},
	}
	output, err = FixReadValues(sipCredsInput, "target_sip_credentials")
	if err != nil {
		t.Errorf("target_sip_credentials input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, sipCredsTarget) != true {
		t.Errorf("Values not stripped correctly from target_sip_credentials input: Received %#v Expected %#v", output, accountsTarget)
	}

	//	tests
	testsInput := []interface{}{
		map[string]interface{}{
			"test_name": "foo",
			"test_id":   "1",
		},
		map[string]interface{}{
			"test_name": "bar",
			"test_id":   "2",
		},
	}
	testsTarget := []interface{}{
		map[string]interface{}{
			"test_id": "1",
		},
		map[string]interface{}{
			"test_id": "2",
		},
	}
	output, err = FixReadValues(testsInput, "tests")
	if err != nil {
		t.Errorf("tests input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, testsTarget) != true {
		t.Errorf("Values not stripped correctly from tests input: Received %#v Expected %#v", output, nil)
	}

	// thirdParty notifications
	thirdPartyNotificationsInput := []interface{}{
		map[string]interface{}{
			"integration_id":   "sl-0000",
			"integration_type": "SLACK",
			"integration_name": "bitconnect",
			"target":           "https://slack.com/waso",
			"channel":          "#terraform",
		},
		map[string]interface{}{
			"integration_id":   "pgd-0000",
			"integration_type": "PAGER_DUTY",
			"integration_name": "PagerDuty notification",
			"auth_method":      "Auth Token",
		},
	}
	thirdPartyNotificationsTarget := []interface{}{
		map[string]interface{}{
			"integration_id":   "sl-0000",
			"integration_type": "SLACK",
		},
		map[string]interface{}{
			"integration_id":   "pgd-0000",
			"integration_type": "PAGER_DUTY",
		},
	}

	output, err = FixReadValues(thirdPartyNotificationsInput, "third_party")
	if err != nil {
		t.Errorf("third party notifications input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, thirdPartyNotificationsTarget) != true {
		t.Errorf("Values not stripped correctly from third party notifications input: Received %#v Expected %#v", output, thirdPartyNotificationsTarget)
	}
}

func TestResourceUpdate(t *testing.T) {

}

func TestResourceSchemaBuild(t *testing.T) {
	type refStruct struct {
		FieldName string `json:"fieldName"`
		Port      int    `json:"port"`
	}

	refSchema := map[string]*schema.Schema{
		"field_name": {
			Type: schema.TypeString,
		},
		"port": {
			Type:     schema.TypeInt,
			Default:  41953,
			Optional: true,
		},
	}

	// test with no schema override
	schm := ResourceSchemaBuild(refStruct{}, refSchema, nil)

	for k, v := range refSchema {
		if _, ok := schm[k]; !ok {
			t.Errorf("Key %s missing from generated schema", k)
		}
		if reflect.DeepEqual(*v, *schm[k]) != true {
			t.Errorf("Schemas not equal: Reference schema is %+v\nNew Schema is %+v", refSchema, schm)
		}
	}

	// test with a schema override
	schemaOverride := map[string]*schema.Schema{
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	expectedSchemaWithOverride := map[string]*schema.Schema{
		"field_name": {
			Type: schema.TypeString,
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	ovrSchm := ResourceSchemaBuild(refStruct{}, refSchema, schemaOverride)

	for k, v := range expectedSchemaWithOverride {
		if _, ok := ovrSchm[k]; !ok {
			t.Errorf("Key %s missing from generated schema", k)
		}
		if reflect.DeepEqual(*v, *ovrSchm[k]) != true {
			t.Errorf("Schemas not equal: Expected schema is %+v\nNew Schema is %+v", expectedSchemaWithOverride, ovrSchm)
		}
	}
}

func TestFillValue(t *testing.T) {
	type refStruct struct {
		FieldName string `json:"fieldName"`
	}
	attrs := map[string]string{}
	testSchemas := map[string]*schema.Schema{
		"testInt": {
			Type: schema.TypeInt,
		},
		"testIntAsString": {
			Type: schema.TypeString,
		},
		"testString": {
			Type: schema.TypeString,
		},
		"testSlice": {
			Type: schema.TypeList,
			Elem: schema.TypeInt,
		},
		"testStruct": {
			Type: schema.TypeMap,
		},
		"testStructSlice": {
			Type: schema.TypeList,
			Elem: schema.TypeMap,
		},
	}
	d := getReferenceData(testSchemas, attrs)

	var err error

	// Integer test
	err = d.Set("testInt", 42)
	if err != nil {
		t.Errorf("Error setting resourceData for 'testInt': %+v", err.Error())
	}
	testInt := FillValue(d.Get("testInt"), 1).(int)
	if testInt != 42 {
		t.Errorf("Expected int '42' for testInt, but received '%+v' of type %+v", reflect.ValueOf(testInt), reflect.TypeOf(testInt))
	}
	err = d.Set("testIntAsString", "27")
	if err != nil {
		t.Errorf("Error setting resourceData for 'testInt': %+v", err.Error())
	}
	testIntAsString := FillValue(d.Get("testIntAsString"), 1).(int)
	if testIntAsString != 27 {
		t.Errorf("Expected int '27' for testIntAsString, but received '%+v' of type %+v", reflect.ValueOf(testIntAsString), reflect.TypeOf(testIntAsString))
	}

	// String test - should test default return path
	err = d.Set("testString", "foo")
	if err != nil {
		t.Errorf("Error setting resourceData for 'testString': %+v", err.Error())
	}
	testString := FillValue(d.Get("testString"), "bar").(string)
	if testString != "foo" {
		t.Errorf("Expected string 'foo', but received %+v of type %+v", reflect.ValueOf(testString), reflect.TypeOf(testString))
	}

	// Slice test
	err = d.Set("testSlice", []int{0, 1})
	if err != nil {
		t.Errorf("Error setting resourceData for 'testSlice': %+v", err.Error())
	}
	testSlice := FillValue(d.Get("testSlice"), []int{}).([]int)
	if reflect.DeepEqual(testSlice, []int{0, 1}) != true {
		t.Errorf("Expected []int{0,1} for testSlice, but received '%+v' of type %+v", reflect.ValueOf(testSlice), reflect.TypeOf(testSlice))
	}

	// Struct test
	refMap := map[string]string{
		"field_name": "value foo",
	}
	cmpStruct := refStruct{
		FieldName: "value foo",
	}
	err = d.Set("testStruct", refMap)
	if err != nil {
		t.Errorf("Error setting resourceData for 'testStruct': %+v", err.Error())
	}
	testStruct := FillValue(d.Get("testStruct"), refStruct{}).(refStruct)
	if reflect.DeepEqual(testStruct, cmpStruct) != true {
		t.Errorf("testStruct doesn't match cmpStruct\ntestStruct: %+v\ncmpStruct: %+v", testStruct, cmpStruct)
	}

	// Struct from slice test
	refSlice := []map[string]string{
		refMap,
	}
	err = d.Set("testStructSlice", refSlice)
	if err != nil {
		t.Errorf("Error setting resourceData for 'testStruct': %+v", err.Error())
	}
	testStruct = FillValue(d.Get("testStruct"), refStruct{}).(refStruct)
	if reflect.DeepEqual(testStruct, cmpStruct) != true {
		t.Errorf("testStruct doesn't match cmpStruct\ntestStruct: %+v\ncmpStruct: %+v", testStruct, cmpStruct)
	}

}

func TestUnderscoreToLowerCamelCase(t *testing.T) {
	s1 := UnderscoreToLowerCamelCase("a_test")
	if s1 != "aTest" {
		t.Errorf("String 'a_test' was not transformed to 'aTest', but instead '%s'", s1)
	}
	s2 := UnderscoreToLowerCamelCase("ateSt")
	if s2 != "atest" {
		t.Errorf("String 'aeSt' was not transformed to 'atest', but instead '%s'", s2)
	}
	s3 := UnderscoreToLowerCamelCase("_a_test")
	if s3 != "aTest" {
		t.Errorf("String '_atest' was not transformed to 'atest', but instead '%s'", s3)
	}
}

func TestCamelCaseToUnderscore(t *testing.T) {
	s1 := CamelCaseToUnderscore("CamelCase")
	if s1 != "camel_case" {
		t.Errorf("String 'CamelCase' was not transformed to 'camel_case', but instead '%s'", s1)
	}
	s2 := CamelCaseToUnderscore("camelCase")
	if s2 != "camel_case" {
		t.Errorf("String 'camelCase' was not transformed to 'camel_case', but instead '%s'", s2)
	}
	s3 := CamelCaseToUnderscore("camelCase")
	if s3 != "camel_case" {
		t.Errorf("String '_camelCase' was not transformed to 'camel_case', but instead '%s'", s3)
	}
}

func TestGetJSONKey(t *testing.T) {
	type testStruct struct {
		TestField int `json:"testField"`
	}
	ts := testStruct{}
	v := reflect.ValueOf(&ts).Elem()
	vt := reflect.TypeOf(v.Interface())
	tag := GetJSONKey(vt.Field(0))
	if tag != "testField" {
		t.Errorf("Field name should be 'testField', but received '%s'", tag)
	}
}
