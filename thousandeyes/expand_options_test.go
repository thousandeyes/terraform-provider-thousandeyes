package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestKnownExpandOptionsExcludeUnknown(t *testing.T) {
	assert.NotContains(t, knownExpandTestOptions(), tests.EXPANDTESTOPTIONS_UNKNOWN)
	assert.NotContains(t, knownExpandBgpTestOptions(), tests.EXPANDBGPTESTOPTIONS_UNKNOWN)
	assert.NotContains(t, knownExpandAccountGroupOptions(), administrative.EXPANDACCOUNTGROUPOPTIONS_UNKNOWN)
	assert.NotContains(t, knownExpandTagsOptions(), tags.EXPANDTAGSOPTIONS_UNKNOWN)
}
