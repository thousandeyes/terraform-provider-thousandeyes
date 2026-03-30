package thousandeyes

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

// unsupportedWidgetTestStub is a unique type used only in tests to exercise the
// unknown-widget branch without depending on a specific SDK widget (which may
// become supported later).
type unsupportedWidgetTestStub struct{}

func TestMapWidgetNilInstance(t *testing.T) {
	data, err := MapWidget(dashboards.ApiWidget{})
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestMapWidgetKnownType(t *testing.T) {
	w := dashboards.NewApiTimeseriesWidget("Time Series: Line")
	w.SetTitle("T")
	widget := dashboards.ApiTimeseriesWidgetAsApiWidget(w)

	data, err := MapWidget(widget)
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.Equal(t, "Time Series: Line", data["type"])
	assert.Equal(t, "T", data["title"])
}

func TestMapWidgetUnknownInstanceType(t *testing.T) {
	// Nil pointer value, non-nil dynamic type *unsupportedWidgetTestStub — hits default in widgetTypeFromInstance.
	_, err := mapWidgetWithInstance(dashboards.ApiWidget{}, (*unsupportedWidgetTestStub)(nil))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown widget type")
	assert.Contains(t, err.Error(), "unsupportedWidgetTestStub")
}

func TestMapWidgetsEmpty(t *testing.T) {
	out, err := MapWidgets(nil)
	assert.NoError(t, err)
	assert.Nil(t, out)

	out, err = MapWidgets([]dashboards.ApiWidget{})
	assert.NoError(t, err)
	assert.Nil(t, out)
}

func TestMapAllWidgetsPropagatesMapperError(t *testing.T) {
	boom := errors.New("mapper failed")
	_, err := mapAllWidgets([]dashboards.ApiWidget{{}}, func(dashboards.ApiWidget) (map[string]interface{}, error) {
		return nil, boom
	})
	require.Error(t, err)
	assert.Equal(t, boom, err)
}
