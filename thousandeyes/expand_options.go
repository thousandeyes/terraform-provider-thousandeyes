package thousandeyes

import (
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func knownExpandTestOptions() []tests.ExpandTestOptions {
	return filterUnknownEnumValues(
		tests.AllowedExpandTestOptionsEnumValues,
		tests.EXPANDTESTOPTIONS_UNKNOWN,
	)
}

func knownExpandBgpTestOptions() []tests.ExpandBgpTestOptions {
	return filterUnknownEnumValues(
		tests.AllowedExpandBgpTestOptionsEnumValues,
		tests.EXPANDBGPTESTOPTIONS_UNKNOWN,
	)
}

func knownExpandAccountGroupOptions() []administrative.ExpandAccountGroupOptions {
	return filterUnknownEnumValues(
		administrative.AllowedExpandAccountGroupOptionsEnumValues,
		administrative.EXPANDACCOUNTGROUPOPTIONS_UNKNOWN,
	)
}

func knownExpandTagsOptions() []tags.ExpandTagsOptions {
	return filterUnknownEnumValues(
		tags.AllowedExpandTagsOptionsEnumValues,
		tags.EXPANDTAGSOPTIONS_UNKNOWN,
	)
}

func filterUnknownEnumValues[T comparable](values []T, unknown T) []T {
	filtered := make([]T, 0, len(values))
	for _, value := range values {
		if value == unknown {
			continue
		}
		filtered = append(filtered, value)
	}
	return filtered
}
