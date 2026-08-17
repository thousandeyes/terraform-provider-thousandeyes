package thousandeyes

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

const oauth2TokenVariableFormat = "OAuth2_Token_Step_%d"

// normalizeAPIResponseRequests removes OAuth2 implementation details that the
// API adds to a request but the user did not configure. Keeping these values in
// state causes Terraform to plan their removal after every refresh.
func normalizeAPIResponseRequests(d rawConfigReader, response *tests.ApiTestResponse) {
	if response == nil {
		return
	}

	configuredRequests, configAvailable := rawConfiguredAPIRequests(d)
	executionStep := 1
	for i := range response.Requests {
		request := &response.Requests[i]
		if request.AuthType != nil && *request.AuthType == tests.APIREQUESTAUTHTYPE_OAUTH2 {
			tokenName := fmt.Sprintf(oauth2TokenVariableFormat, executionStep)
			headerValue := "{{" + tokenName + "}}"

			request.Headers = filterAPIRequestHeaders(
				request.Headers,
				configuredRequests,
				configAvailable,
				i,
				headerValue,
			)
			request.Variables = filterAPIRequestVariables(
				request.Variables,
				configuredRequests,
				configAvailable,
				i,
				tokenName,
			)

			// OAuth2 token acquisition is an execution step of its own.
			executionStep++
		}

		// Every configured request consumes one execution step.
		executionStep++
	}
}

func filterAPIRequestHeaders(headers []tests.ApiRequestHeader, configuredRequests []cty.Value, configAvailable bool, requestIndex int, headerValue string) []tests.ApiRequestHeader {
	want := map[string]string{
		"key":   "Authorization",
		"value": headerValue,
	}
	dropSynthetic := shouldDropGeneratedAPIRequestValue(configuredRequests, configAvailable, requestIndex, "headers", want)

	filtered := make([]tests.ApiRequestHeader, 0, len(headers))
	for _, header := range headers {
		if dropSynthetic && header.Key != nil && header.Value != nil &&
			*header.Key == "Authorization" && *header.Value == headerValue {
			continue
		}
		filtered = append(filtered, header)
	}
	return filtered
}

func filterAPIRequestVariables(variables []tests.ApiRequestVariable, configuredRequests []cty.Value, configAvailable bool, requestIndex int, tokenName string) []tests.ApiRequestVariable {
	want := map[string]string{
		"name":  tokenName,
		"value": "",
	}
	dropSynthetic := shouldDropGeneratedAPIRequestValue(configuredRequests, configAvailable, requestIndex, "variables", want)

	filtered := make([]tests.ApiRequestVariable, 0, len(variables))
	for _, variable := range variables {
		if dropSynthetic && variable.Name != nil && variable.Value != nil &&
			*variable.Name == tokenName && *variable.Value == "" {
			continue
		}
		filtered = append(filtered, variable)
	}
	return filtered
}

func rawConfiguredAPIRequests(d rawConfigReader) ([]cty.Value, bool) {
	if d == nil {
		return nil, false
	}

	raw, diags := d.GetRawConfigAt(cty.Path{cty.GetAttrStep{Name: "requests"}})
	if diags.HasError() || !raw.IsKnown() || raw.IsNull() {
		return nil, false
	}

	configured := make([]cty.Value, 0, raw.LengthInt())
	it := raw.ElementIterator()
	for it.Next() {
		_, request := it.Element()
		configured = append(configured, request)
	}
	return configured, true
}

func shouldDropGeneratedAPIRequestValue(configuredRequests []cty.Value, configAvailable bool, requestIndex int, field string, want map[string]string) bool {
	if !configAvailable || requestIndex >= len(configuredRequests) {
		return true
	}

	configured, known := configuredAPIRequestCollectionContains(configuredRequests[requestIndex], field, want)
	if configured {
		return false
	}
	// If configuration contains unknown values, preserve the remote value
	// because it may resolve to the generated-looking value the user intended.
	return known
}

func configuredAPIRequestCollectionContains(request cty.Value, field string, want map[string]string) (bool, bool) {
	if !request.IsKnown() {
		return false, false
	}
	if request.IsNull() || !request.Type().IsObjectType() || !request.Type().HasAttribute(field) {
		return false, true
	}

	collection := request.GetAttr(field)
	if !collection.IsKnown() {
		return false, false
	}
	if collection.IsNull() {
		return false, true
	}

	allKnown := true
	it := collection.ElementIterator()
	for it.Next() {
		_, item := it.Element()
		matches, known := configuredAPIRequestObjectMatches(item, want)
		if matches {
			return true, true
		}
		allKnown = allKnown && known
	}
	return false, allKnown
}

func configuredAPIRequestObjectMatches(item cty.Value, want map[string]string) (bool, bool) {
	if !item.IsKnown() {
		return false, false
	}
	if item.IsNull() || !item.Type().IsObjectType() {
		return false, true
	}

	hasUnknown := false
	for field, expected := range want {
		if !item.Type().HasAttribute(field) {
			return false, true
		}
		value := item.GetAttr(field)
		if !value.IsKnown() {
			hasUnknown = true
			continue
		}
		if value.IsNull() || value.AsString() != expected {
			return false, true
		}
	}
	if hasUnknown {
		return false, false
	}
	return true, true
}
