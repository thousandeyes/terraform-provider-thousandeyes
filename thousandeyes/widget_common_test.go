package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestSetCommonBuilderFields_PreservesFalseBooleans(t *testing.T) {
	widget := dashboards.NewApiGeoMapWidget("Map")

	setCommonBuilderFields(widget, map[string]interface{}{
		"is_embedded": false,
		"should_exclude_alert_suppression_windows": false,
	})

	isEmbedded, ok := widget.GetIsEmbeddedOk()
	assert.True(t, ok)
	if assert.NotNil(t, isEmbedded) {
		assert.False(t, *isEmbedded)
	}

	shouldExclude, ok := widget.GetShouldExcludeAlertSuppressionWindowsOk()
	assert.True(t, ok)
	if assert.NotNil(t, shouldExclude) {
		assert.False(t, *shouldExclude)
	}
}

func TestSetCommonBuilderFields_PreservesTrueBooleans(t *testing.T) {
	widget := dashboards.NewApiGeoMapWidget("Map")

	setCommonBuilderFields(widget, map[string]interface{}{
		"is_embedded": true,
		"should_exclude_alert_suppression_windows": true,
	})

	isEmbedded, ok := widget.GetIsEmbeddedOk()
	assert.True(t, ok)
	if assert.NotNil(t, isEmbedded) {
		assert.True(t, *isEmbedded)
	}

	shouldExclude, ok := widget.GetShouldExcludeAlertSuppressionWindowsOk()
	assert.True(t, ok)
	if assert.NotNil(t, shouldExclude) {
		assert.True(t, *shouldExclude)
	}
}
