package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestSetCommonBuilderFields_IgnoresIsEmbedded(t *testing.T) {
	widget := dashboards.NewApiGeoMapWidget("Map")

	setCommonBuilderFields(widget, map[string]interface{}{
		"is_embedded": true,
	})

	isEmbedded, ok := widget.GetIsEmbeddedOk()
	assert.False(t, ok)
	assert.Nil(t, isEmbedded)
}

func TestSetCommonBuilderFields_PreservesFalseBooleans(t *testing.T) {
	widget := dashboards.NewApiGeoMapWidget("Map")

	setCommonBuilderFields(widget, map[string]interface{}{
		"should_exclude_alert_suppression_windows": false,
	})

	shouldExclude, ok := widget.GetShouldExcludeAlertSuppressionWindowsOk()
	assert.True(t, ok)
	if assert.NotNil(t, shouldExclude) {
		assert.False(t, *shouldExclude)
	}
}

func TestSetCommonBuilderFields_PreservesTrueBooleans(t *testing.T) {
	widget := dashboards.NewApiGeoMapWidget("Map")

	setCommonBuilderFields(widget, map[string]interface{}{
		"should_exclude_alert_suppression_windows": true,
	})

	shouldExclude, ok := widget.GetShouldExcludeAlertSuppressionWindowsOk()
	assert.True(t, ok)
	if assert.NotNil(t, shouldExclude) {
		assert.True(t, *shouldExclude)
	}
}
