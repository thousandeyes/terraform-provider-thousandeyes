package thousandeyes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

type resourceKeyType string
type setInConfigKeyType string

const tagsKey resourceKeyType = "tags"
const setInConfigKey setInConfigKeyType = "is_set"

var sensitiveFields = []string{"password", "custom_headers", "headers", "bearer_token", "client_id", "client_secret"}

// emptyStringToNilTypes contains type names that should be converted to nil when their value is an empty string.
// This handles cases where the Terraform SDK v2 converts null values to empty strings, but the API expects
// the field to be omitted (via omitempty) rather than sent as an empty string.
// Store the actual type names as they appear in the SDK (e.g., "ApiClientAuthentication")
var emptyStringToNilTypes = []string{
	"ApiClientAuthentication",
}

type ResourceReadFunc func(client *client.APIClient, id string) (interface{}, error)

type RequestWithAid[T any] interface {
	Aid(aid string) T
}

func IsNotFoundError(err error) bool {
	notFoundPatterns := []string{"404", "not found"}
	for _, pattern := range notFoundPatterns {
		if strings.Contains(strings.ToLower(err.Error()), pattern) {
			return true
		}
	}
	return false
}

func expandAgents(v interface{}) []tests.TestAgentRequest {
	agents := make([]tests.TestAgentRequest, 0)
	var agentsIDs []interface{}
	if rawAgents, ok := v.(*schema.Set); ok {
		agentsIDs = rawAgents.List()
	}
	for _, item := range agentsIDs {
		id := item.(string)
		if len(id) == 0 {
			continue
		}
		agents = append(agents, tests.TestAgentRequest{
			AgentId: id,
		})
	}
	return agents
}

// ResourceBuildStruct fills the struct at a given address by querying a
// schema.ResourceData object for the matching field.  It discovers the
// matching value name by getting the JSON key from the struct field,
// and then fills in the value according to the struct field's type.
func ResourceBuildStruct[T any](d *schema.ResourceData, structPtr *T) *T {
	v := reflect.ValueOf(structPtr).Elem()
	t := reflect.TypeOf(v.Interface())
	for i := 0; i < v.NumField(); i++ {
		tag := GetJSONKey(t.Field(i))
		tfName := CamelCaseToUnderscore(tag)
		val, ok := d.GetOkExists(tfName)
		if ok {
			newVal := FillValue(val, v.Field(i).Interface())
			setVal := reflect.ValueOf(newVal)
			v.Field(i).Set(setVal)
		}
	}
	return resourceFixups(d, structPtr)
}

// GetResource is a generic function for reading resources.
func GetResource(ctx context.Context, d *schema.ResourceData, m interface{}, readFunc ResourceReadFunc) error {
	apiClient := m.(*client.APIClient)
	if aid, ok := apiClient.GetConfig().Context.Value(accountGroupIdKey).(string); ok {
		ctx = context.WithValue(ctx, accountGroupIdKey, aid)
	}

	log.Printf("[INFO] Reading Thousandeyes Resource %s", d.Id())
	remote, err := readFunc(apiClient, d.Id())

	// Check if the resource no longer exists
	if err != nil && IsNotFoundError(err) {
		log.Printf("[INFO] Resource was deleted - will recreate it")
		d.SetId("") // Set ID to empty to mark the resource as non-existent
		return nil
	} else if err != nil {
		return err
	}

	// Continue with updating the state
	err = ResourceRead(ctx, d, remote)
	if err != nil {
		return err
	}

	return nil
}

// ResourceRead sets values for a schema.ResourceData object by names derived
// from the fields of the struct at the provided pointer.
func ResourceRead(ctx context.Context, d *schema.ResourceData, structPtr interface{}) error {
	v := reflect.ValueOf(structPtr).Elem()
	t := reflect.TypeOf(v.Interface())

	targetMaps := getTargetFieldsMaps(structPtr)

	for i := 0; i < v.NumField(); i++ {
		tag := GetJSONKey(t.Field(i))
		tfName := CamelCaseToUnderscore(tag)
		_, ok := d.GetOk(tfName)

		if slices.Contains(sensitiveFields, tfName) {
			if ok {
				if err := d.Set(tfName, nil); err != nil {
					return err
				}
			}
			continue
		}

		if v.Field(i).Kind() == reflect.Ptr && v.Field(i).IsNil() {
			continue
		}

		val, err := ReadValue(v.Field(i).Interface())
		if err != nil {
			return err
		}

		ctx = context.WithValue(ctx, setInConfigKey, ok)
		val, err = FixReadValues(ctx, targetMaps, val, &tfName)
		if err != nil {
			return err
		}
		if len(tfName) == 0 {
			continue
		}

		err = d.Set(tfName, val)
		if err != nil {
			return err
		}
	}

	for k, v := range targetMaps {
		if err := d.Set(k, []interface{}{v}); err != nil {
			return err
		}
	}

	return nil
}

// getTargetFieldsMaps gets a map of target fields for a specific resource when multiple fields need to be set in a single target map.
func getTargetFieldsMaps(structPtr interface{}) map[string]map[string]interface{} {
	switch structPtr.(type) {
	// Example:
	// case (tests.Example):
	// 	res := make(map[string]map[string]interface{})
	// 	res["TARGET_FIELD"] = map[string]interface{}{
	// 		"SOURCE_FIELD_1":     nil,
	// 		"SOURCE_FIELD_2":     nil,
	// 		"SOURCE_FIELD_3":     nil,
	// 		...
	// 	}
	// 	return res
	case (*tests.SipServerTestResponse):
		res := make(map[string]map[string]interface{})
		res["target_sip_credentials"] = map[string]interface{}{
			"auth_user":     nil,
			"port":          nil,
			"protocol":      nil,
			"sip_registrar": nil,
			"user":          nil,
		}
		return res
	}

	return nil
}

// FixReadValues adjusts certain values returned from ThousandEyes to make them
// processable by this Terraform plugin.  This includes removing extraneous
// information that ThousandEyes returns when querying certain resources (ie,
// when querying a group it may return a list of associated tests with details)
// and transforms certain values to match the expected schema.
// We need to account for this data on so that it does not get saved to state and
// cause conflict with configuration.
func FixReadValues(ctx context.Context, targetMaps map[string]map[string]interface{}, m interface{}, name *string) (interface{}, error) {
	aid, _ := ctx.Value(accountGroupIdKey).(string)

	// Set fields into map to match schema
	for targetField := range targetMaps {
		if _, ok := targetMaps[targetField][*name]; ok {
			targetMaps[targetField][*name] = m
			*name = ""
			return nil, nil
		}
	}

	switch *name {
	// Remove all fields from user definitions except for user ID.
	case "users":
		for i, v := range m.([]interface{}) {
			user := v.(map[string]interface{})
			m.([]interface{})[i] = user["uid"]
		}

	// Remove all fields from permission definitions except for permission ID.
	case "permissions":
		for i, v := range m.([]interface{}) {
			perm := v.(map[string]interface{})
			m.([]interface{})[i] = perm["permission_id"]
		}

	// Get login account group ID
	case "login_account_group":
		*name = "login_account_group_id"
		if lag, ok := m.(map[string]interface{}); ok {
			m = lag["aid"]
		}

	// Remove all fields from all account group roles definitions except for ID.
	case "all_account_group_roles":
		*name = "all_account_group_role_ids"
		for i, v := range m.([]interface{}) {
			role := v.(map[string]interface{})
			m.([]interface{})[i] = role["role_id"]
		}

	// Parse account group roles to schema.
	case "account_group_roles":
		for i, v := range m.([]interface{}) {
			agr := v.(map[string]interface{})
			roles := agr["roles"].([]interface{})
			ag := agr["account_group"].(map[string]interface{})

			roleIds := make([]string, 0, len(roles))
			for _, r := range roles {
				role := r.(map[string]interface{})
				roleIds = append(roleIds, *role["role_id"].(*string))
			}

			m.([]interface{})[i] = map[string]interface{}{
				"account_group_id": ag["aid"],
				"role_ids":         roleIds,
			}
		}

	// Remove all fields from agent definitions except for agent ID.
	case "agents":
		for i, v := range m.([]interface{}) {
			agent := v.(map[string]interface{})
			m.([]interface{})[i] = agent["agent_id"]
		}

	// Ignore emulated device ID if it wasn't set
	case "emulated_device_id":
		if isSet, _ := ctx.Value(setInConfigKey).(bool); !isSet {
			*name = ""
			return nil, nil
		}

	// Ignore MTU measurements if it wasn't set, prevents API side effects (network_measurements) from being saved to state
	case "mtu_measurements":
		if isSet, _ := ctx.Value(setInConfigKey).(bool); !isSet {
			*name = ""
			return nil, nil
		}

	// Ignore num path traces if it wasn't set, prevents API side effects (network_measurements) from being saved to state
	case "num_path_traces":
		if isSet, _ := ctx.Value(setInConfigKey).(bool); !isSet {
			*name = ""
			return nil, nil
		}

	// Return only host when host:port pattern obtained
	case "server":
		m = strings.Split(m.(string), ":")[0]

	// Remove all alert rule fields except for rule ID. Ignore default rules.
	// Remove all alert rule fields except for rule ID.
	case "alert_rules":
		for i, v := range m.([]interface{}) {
			rule := v.(map[string]interface{})
			m.([]interface{})[i] = rule["rule_id"]
		}

	// Remove all public BGP monitors. (ThousandEyes does not allow
	// specifying individual public BGP monitors, and all available
	// public BGP monitors are returned if public BGP monitors are enabled.)
	case "monitors":
		monitors := m.([]interface{})
		// Edit the monitors slice in place, to return the same type.
		i := 0
		for i < len(monitors) {
			monitor := monitors[i].(map[string]interface{})
			if *monitor["monitor_type"].(*tests.MonitorType) == tests.MONITORTYPE_PUBLIC {
				// Remove this item from the slice
				monitors = append(monitors[:i], monitors[i+1:]...)
			} else {
				monitors[i] = monitor["monitor_id"]
				i = i + 1
			}
		}
		m = monitors

	// Remove all dns_server fields except for the server name.
	case "dns_servers":
		for i, v := range m.([]interface{}) {
			servers := v.(map[string]interface{})
			m.([]interface{})[i] = servers["server_name"]
		}

	// Remove all label fields except for the label ID.
	// This is required because some customers are still using the old provider v1
	// In v1 it was possible to assign labels to tests by ID, but in v2+ it is not.
	// Eventually this can be removed
	case "labels":
		for i, v := range m.([]interface{}) {
			label := v.(map[string]interface{})
			m.([]interface{})[i] = label["label_id"]
		}

	// custom_headers is currently unsupported due to complications with Terraform
	// and the object schema.  It will presently be removed from state, and when
	// a solution is found it will be transformed here according to the specification
	// of that solution.
	case "custom_headers":
		m = nil

	// download_limit may appear as a string instead of an integer.
	case "download_limit":
		var err error
		if reflect.TypeOf(m) == reflect.TypeOf("") {
			if m.(string) == "" {
				m = 0
			} else {
				m, err = strconv.Atoi(m.(string))
				if err != nil {
					return nil, err
				}
			}
		}

	// Remove the owning account from the list of shared accounts.
	case "shared_with_accounts":
		accounts := m.([]interface{})
		if len(aid) == 0 {
			if len(accounts) > 1 {
				return nil, errors.New("Resources are shared between account groups, but account_group_id is not set.")
			}
			// A single listed account should be the owning account group.
			if len(accounts) == 1 {
				return nil, nil
			}
		}
		i := 0
		for i < len(accounts) {
			account := accounts[i].(map[string]interface{})
			//  Compare to account group ID stored in global variable.
			shared_aid := account["aid"].(*string)
			if *shared_aid == aid {
				// Remove this item from the slice
				accounts = append(accounts[:i], accounts[i+1:]...)
			} else {
				accounts[i] = shared_aid
				i = i + 1
			}
		}
		m = accounts

	case "notifications":
		var e interface{}
		var err error

		notifications := m.(map[string]interface{})

		// this is a special case to handle internal email structure inside the notifications block
		e, err = FixReadValues(ctx, nil, notifications["email"].(map[string]interface{}), getPointer("email"))
		if err != nil {
			return nil, err
		}

		// third party notifications
		var tp interface{}
		if _, ok := notifications["third_party"]; ok {
			tp, err = FixReadValues(ctx, nil, notifications["third_party"].([]interface{}), getPointer("third_party"))
			if err != nil {
				return nil, err
			}
		} else {
			tp = nil
		}

		// webhook notifications
		var w interface{}
		if _, ok := notifications["webhook"]; ok {
			w, err = FixReadValues(ctx, nil, notifications["webhook"].([]interface{}), getPointer("webhook"))
			if err != nil {
				return nil, err
			}
		} else {
			w = nil
		}

		// custom webhook notifications
		var cw interface{}
		if _, ok := notifications["custom_webhook"]; ok {
			cw, err = FixReadValues(ctx, nil, notifications["custom_webhook"].([]interface{}), getPointer("custom_webhook"))
			if err != nil {
				return nil, err
			}
		} else {
			cw = nil
		}

		// update the notifications block if the email block is present and contains recipients, or
		// the third party notifications are present, or webhook notifications are present.
		// Otherwise set the whole notifications block to nil
		if (e == nil || len(e.([]interface{})) == 0) &&
			(tp == nil || len(tp.([]interface{})) == 0) &&
			(w == nil || len(w.([]interface{})) == 0) &&
			(cw == nil || len(cw.([]interface{})) == 0) {
			// *name = ""
			m = nil
		} else {
			// Add the third party map and or webhook map to the notifications map if they are present
			// if they're not configured, then the API doesn't return them at all
			if tp != nil {
				notifications["third_party"] = tp
			}

			if w != nil {
				notifications["webhook"] = w
			}

			if cw != nil {
				notifications["custom_webhook"] = cw
			}

			notifications["email"] = e
			m = []interface{}{
				notifications,
			}
		}

	case "email":
		if _, ok := m.(*string); !ok {
			if len(m.(map[string]interface{})["recipients"].([]interface{})) == 0 {
				m = nil
			} else {
				m = []interface{}{
					m.(map[string]interface{}),
				}
			}
		}

	// remove all fields except the integration ID and type to
	// mimick the behavior of the example in our docs for a
	// regular create API request for alert rules, where only
	// these two fields are passed
	case "third_party":
		for i, v := range m.([]interface{}) {
			tpn := v.(map[string]interface{})
			m.([]interface{})[i] = map[string]interface{}{
				"integration_id":   tpn["integration_id"],
				"integration_type": tpn["integration_type"],
			}
		}

	case "webhook":
		for i, v := range m.([]interface{}) {
			webhookNotifications := v.(map[string]interface{})
			m.([]interface{})[i] = map[string]interface{}{
				"integration_id":   webhookNotifications["integration_id"],
				"integration_type": webhookNotifications["integration_type"],
			}
		}

	case "custom_webhook":
		for i, v := range m.([]interface{}) {
			webhookNotifications := v.(map[string]interface{})
			m.([]interface{})[i] = map[string]interface{}{
				"integration_id":   webhookNotifications["integration_id"],
				"integration_type": webhookNotifications["integration_type"],
			}
		}

	case "tests":
		*name = "test_ids"
		for i, v := range m.([]interface{}) {
			test := v.(map[string]interface{})
			m.([]interface{})[i] = test["test_id"]
		}

	case "_links":
		*name = "link"
		if self, ok := m.(map[string]interface{})["self"].(map[string]interface{}); ok {
			m = self["href"]
		}

	case "created_date", "modified_date", "date_registered", "last_login":
		{
			m = m.(*time.Time).Format(time.RFC3339)
		}

	// Ignore nullable fields (already set);  skip assignments for Tags (used in Tags Assignments)
	case "icon", "description", "legacy_id", "assignments":
		isTags := ctx.Value(tagsKey)
		if isTags != nil {
			*name = ""
			return nil, nil
		}

	case "aid":
		isTags := ctx.Value(tagsKey)
		if isTags != nil {
			tmp, _ := m.(*int32)
			if tmp != nil {
				m = fmt.Sprintf("%v", *tmp)
			} else {
				*name = ""
				return nil, nil
			}
		}

	case "o_auth":
		*name = "oauth"

	}

	return m, nil
}

// ReadValue returns a value with key names for which Terraform will be able to
// identify in the Schema.  This is required because calling the Set function on
// a struct results in the JSON tag name (instead of the Terraform config key)
// being used for schema lookups.
func ReadValue(structPtr interface{}) (interface{}, error) {
	var err error
	v := reflect.Indirect(reflect.ValueOf(structPtr))
	t := reflect.TypeOf(v.Interface())
	eltype := v.Type()
	switch t.Kind() {
	case reflect.Struct:
		// For structs, return a map with key names set to be translations of
		// the JSON key names.
		if (eltype == reflect.TypeOf(time.Time{})) {
			return structPtr, nil
		}
		newMap := make(map[string]interface{})
		for i := 0; i < v.NumField(); i++ {
			tag := GetJSONKey(t.Field(i))
			tfName := CamelCaseToUnderscore(tag)

			if slices.Contains(sensitiveFields, tfName) {
				newMap[tfName] = nil
				continue
			}
			if v.Field(i).Kind() == reflect.Ptr && v.Field(i).IsNil() {
				continue
			}
			// check for unexported fields
			if v.Field(i).CanInterface() {
				newMap[tfName], err = ReadValue(v.Field(i).Interface())
			}
		}
		if err != nil {
			return nil, err
		}
		return newMap, nil
	case reflect.Slice:
		// If it's a list, create an empty version of
		// that collection type, and then recurse for each child value (passing the
		// extended key name).
		var newSlice []interface{}
		for i := 0; i < v.Len(); i++ {
			newVal, err := ReadValue(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			newSlice = append(newSlice, newVal)
		}
		return newSlice, nil

	default:
		return structPtr, nil
	}
}

// resourceFixups sanitizes values to ensure that the ThousandEyes API
// behavior does not surprise Terraform's state.
func resourceFixups[T any](d *schema.ResourceData, structPtr *T) *T {
	v := reflect.ValueOf(structPtr).Elem()
	t := reflect.TypeOf(v.Interface())

	// When changing networkMeasurements, the ThousandEyes API
	// modifies other flags as well.
	if d.HasChange("network_measurements") {
		// TE API automatically sets bandwidthMeasurements to
		// true. This is not ideal when using cloud agents, as it's
		// not supported. Better to let users explicitly set it to
		// true.
		_, bandwidthMeasurementsSet := d.GetOk("bandwidth_measurements")
		_, hasBandwidthMeasurementsField := t.FieldByName("BandwidthMeasurements")
		if hasBandwidthMeasurementsField && !bandwidthMeasurementsSet {
			setVal := reflect.ValueOf(getPointer(false))
			v.FieldByName("BandwidthMeasurements").Set(setVal)
			d.Set("bandwidth_measurements", false)
		}

		// TE API automatically sets bgpMeasurements to
		// true. This is not ideal when using cloud agents, as it's
		// not supported. Better to let users explicitly set it to
		// true.
		_, bgpMeasurementsSet := d.GetOk("bgp_measurements")
		_, hasBgpMeasurements := t.FieldByName("BgpMeasurements")
		if hasBgpMeasurements && !bgpMeasurementsSet {
			setVal := reflect.ValueOf(getPointer(false))
			v.FieldByName("BgpMeasurements").Set(setVal)
			d.Set("bgp_measurements", false)
		}
	}

	_, hasAgents := t.FieldByName("Agents")
	if hasAgents && t.Name() != "AccountGroupRequest" {
		scrappedAgents := expandAgents(d.Get("agents"))
		v.FieldByName("Agents").Set(reflect.ValueOf(scrappedAgents))
	}

	_, hasOAuth := t.FieldByName("OAuth")
	if hasOAuth {
		val, ok := d.GetOk("oauth")
		if ok {
			newVal := FillValue(val, v.FieldByName("OAuth").Interface())
			setVal := reflect.ValueOf(newVal)
			v.FieldByName("OAuth").Set(setVal)
		}
	}

	// When protocol is ICMP, ensure port is explicitly set to nil
	// This is required when updating from TCP/UDP to ICMP
	_, hasProtocol := t.FieldByName("Protocol")
	_, hasPort := t.FieldByName("Port")
	if hasProtocol && hasPort {
		protocol, ok := d.GetOk("protocol")
		if ok && protocol == "icmp" {
			// Explicitly set port to nil for ICMP tests
			v.FieldByName("Port").Set(reflect.Zero(v.FieldByName("Port").Type()))
		}
	}

	return structPtr
}

// ResourceUpdate updates values of a struct for the provided pointer if
// matching changes for those values are found in a provided
// schema.ResourceData object.
func ResourceUpdate[T any](d *schema.ResourceData, structPtr *T) *T {
	d.Partial(true)
	v := reflect.ValueOf(structPtr).Elem()
	t := reflect.TypeOf(v.Interface())
	for i := 0; i < v.NumField(); i++ {
		tag := GetJSONKey(t.Field(i))
		tfName := CamelCaseToUnderscore(tag)
		if d.HasChange(tfName) {
			newVal := FillValue(d.Get(tfName), v.Field(i).Interface())
			setVal := reflect.ValueOf(newVal)
			v.Field(i).Set(setVal)
		}
	}
	d.Partial(false)
	return resourceFixups(d, structPtr)
}

// ResourceSchemaBuild creates a map of schemas based on the fields
// of the provided struct.
func ResourceSchemaBuild(referenceStruct interface{}, schemas map[string]*schema.Schema, schemasOverride map[string]*schema.Schema) map[string]*schema.Schema {
	newSchema := map[string]*schema.Schema{}
	v := reflect.ValueOf(referenceStruct)
	t := reflect.TypeOf(referenceStruct)

	for i := 0; i < v.NumField(); i++ {
		tag := GetJSONKey(t.Field(i))
		tfName := CamelCaseToUnderscore(tag)

		// use the override if there is one
		if len(schemasOverride) > 0 {
			if val, ok := schemasOverride[tfName]; ok {
				newSchema[tfName] = val
			} else if val, ok := schemas[tfName]; ok {
				newSchema[tfName] = val
			}
		} else {
			if val, ok := schemas[tfName]; ok {
				newSchema[tfName] = val
			}
		}
	}

	// instead of "_links"
	if _, ok := schemas["link"]; ok {
		newSchema["link"] = schemas["link"]
	}

	return newSchema
}

// FillValue takes a value from the Terraform resource data and translates it
// to the correct type, based on the type of the target parameter.
func FillValue(source interface{}, target interface{}) interface{} {
	// We determine how to interpret the supplied value based on
	// the type of the target argument.
	vt := reflect.ValueOf(target)
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)
	switch vt.Kind() {
	case reflect.Ptr:
		// Check if this is an empty string that should be converted to nil for specific types
		// This is determined by checking if the target type name is in our emptyStringToNilTypes list
		if str, isString := source.(string); isString && str == "" {
			targetType := reflect.TypeOf(target)
			if targetType.Elem().Kind() == reflect.String {
				// For string-based types (including enums), check the type name
				typeName := targetType.Elem().Name()
				if slices.Contains(emptyStringToNilTypes, typeName) {
					// Return nil pointer for this type
					return reflect.Zero(targetType).Interface()
				}
			}
		}
		p := reflect.New(reflect.TypeOf(target).Elem())
		newVal := FillValue(source, p.Elem().Interface())
		p.Elem().Set(reflect.ValueOf(newVal))
		return p.Interface()
	case reflect.Slice:
		// When the target is a slice, we create a new slice of the same type,
		// then recurse with the value of each element.
		tt := reflect.TypeOf(target)
		tte := reflect.TypeOf(target).Elem() // The type of items in the slice
		ntte := reflect.New(tte).Elem()
		newSlice := reflect.New(tt).Elem()

		vs := reflect.ValueOf(source)
		// If source is a *Set, we dereference it and convert it to a
		// List so we can iterate over its elements.
		if vs.Type() == reflect.TypeOf(&schema.Set{}) {
			vs = reflect.ValueOf(source.(*schema.Set).List())
		}

		for i := 0; i < vs.Len(); i++ {
			toAppend := FillValue(vs.Index(i).Interface(), ntte.Interface())
			appendVal := reflect.ValueOf(toAppend)
			newSlice = reflect.Append(newSlice, appendVal)
		}
		return newSlice.Interface()
	case reflect.Struct:
		// When the target is a struct, we assume that the source is a map
		// containing values corresponding to the struct's fields, then
		// recurse on each value looked up to get the value to be set.

		// Due to limitations of Terraform's schema handling, some maps may
		// be delivered inside single-item slices.  This occurs when maps
		// must be declared as lists of terraform resources, whether to
		// define specific key names or to have values of mixed types,
		// neither of which is supported by Terraform's implementation of
		// maps.
		vs := reflect.ValueOf(source)
		structSource := source
		if vs.Kind() == reflect.Slice {
			structSource = source.([]interface{})[0]
		} else if vs.Kind() == reflect.Ptr {
			structSource = source.(*schema.Set).List()
			if len(structSource.([]interface{})) != 0 {
				structSource = structSource.([]interface{})[0]
			} else {
				source = nil
			}
		}
		t := reflect.TypeOf(vt.Interface())
		newStruct := reflect.New(t).Interface()
		setStruct := reflect.ValueOf(newStruct).Elem()
		if source != nil {
			m, ok := structSource.(map[string]interface{})
			if !ok {
				return setStruct.Interface()
			}
			for i := 0; i < vt.NumField(); i++ {
				tag := GetJSONKey(t.Field(i))
				tfName := CamelCaseToUnderscore(tag)
				if mv, ok := m[tfName]; ok {
					newVal := FillValue(mv, vt.Field(i).Interface())
					setStruct.Field(i).Set(reflect.ValueOf(newVal))
				}
			}
		}
		return setStruct.Interface()
	case reflect.Int:
		// Values destined to be ints may come to us as strings.
		if reflect.TypeOf(source).Kind() == reflect.String {
			i, _ := strconv.ParseInt(source.(string), 10, 32)
			return int(i)
		}

		return source

	case reflect.Int64:
		// Values destined to be int64 may come to us as strings.
		if sourceType.ConvertibleTo(vt.Type()) {
			return sourceValue.Convert(vt.Type()).Interface()
		}
		if reflect.TypeOf(source).Kind() == reflect.String {
			i, _ := strconv.ParseInt(source.(string), 10, 64)
			return i
		}

		return int64(source.(int))

	case reflect.Int32:
		// Values destined to be int32 may come to us as strings.
		if sourceType.ConvertibleTo(vt.Type()) {
			return sourceValue.Convert(vt.Type()).Interface()
		}
		if reflect.TypeOf(source).Kind() == reflect.String {
			i, _ := strconv.ParseInt(source.(string), 10, 32)
			return int32(i)
		}

		return int32(source.(int))

	default:
		// If we haven't matched one of the above cases, then there
		// is likely no reason to translate.
		if sourceType.ConvertibleTo(vt.Type()) {
			return sourceValue.Convert(vt.Type()).Interface()
		}
		return source
	}
}

// UnderscoreToLowerCamelCase translates from words separated by
// underscores to camel case with initial lowercase.
// ie, a_string would become aString
func UnderscoreToLowerCamelCase(s string) string {
	// We have a map of exceptions to the usual conversion logic.
	exceptions := map[string]string{
		"ip_addresses": "IPAddresses",
	}
	if val, ok := exceptions[s]; ok {
		return val
	}
	s = strings.ToLower(s)
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	s = strings.Replace(s, " ", "", -1)
	firstChar := string(s[0])
	s = strings.Replace(s, firstChar, strings.ToLower(firstChar), 1)
	return s
}

// CamelCaseToUnderscore translates from camel case (with any leading case)
// to underscore separated words.
// ie, either aString and AString would become a_string
func CamelCaseToUnderscore(s string) string {
	input := []rune(s)
	output := []rune{}
	for i, r := range input {
		if unicode.IsUpper(r) {
			if i != 0 && i < len(input)-1 {
				if unicode.IsLower(input[i+1]) {
					output = append(output, []rune("_")[0])
				}
			}
			output = append(output, unicode.ToLower(r))
		} else {
			output = append(output, r)
		}
	}
	return string(output)
}

// GetJSONKey returns the JSON object key for the struct which is represented
// by the passed reflect.StructField instance.
func GetJSONKey(v reflect.StructField) string {
	s := v.Tag.Get("json")
	return strings.Split(s, ",")[0]
}

func SetAidFromContext[T RequestWithAid[T]](ctx context.Context, req T) T {
	aid, ok := ctx.Value(accountGroupIdKey).(string)
	if ok && len(aid) > 0 {
		return req.Aid(aid)
	}
	return req
}

func getPointer[T any](v T) *T {
	return &v
}

func checkDomainRecordTypeExists(domain string) bool {
	matched, _ := regexp.MatchString(`^.+ ([A-Z]+)$`, domain)
	return matched
}
