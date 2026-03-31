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
	assert.Error(t, err)
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

func TestMapWidgetUnknownInstanceType_returnsNil(t *testing.T) {
	data, err := mapWidgetWithInstance(dashboards.ApiWidget{}, (*unsupportedWidgetTestStub)(nil))
	require.NoError(t, err)
	assert.Nil(t, data)
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

func TestBuildWidgetRequiresType(t *testing.T) {
	_, err := BuildWidget(map[string]interface{}{"title": "x"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "type")

	_, err = BuildWidget(map[string]interface{}{"type": ""})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "type")
}

func TestBuildWidgetUnsupportedType(t *testing.T) {
	_, err := BuildWidget(map[string]interface{}{"type": "Color Grid"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported widget type")
}

func TestBuildWidgetsRejectsNonObject(t *testing.T) {
	_, err := BuildWidgets([]interface{}{"not-a-map"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "widgets.0")
}

func TestBuildWidgetsWrapsBuildWidgetError(t *testing.T) {
	_, err := BuildWidgets([]interface{}{
		map[string]interface{}{"type": "Time Series: Line", "title": "ok"},
		map[string]interface{}{"type": "Unknown Widget"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "widgets.1")
}

func TestWidgetTypeFromInstance_unmanagedReturnsEmpty(t *testing.T) {
	wType, err := widgetTypeFromInstance((*unsupportedWidgetTestStub)(nil))
	assert.NoError(t, err)
	assert.Equal(t, "", wType)
}

func TestIsUnmanagedWidget(t *testing.T) {
	t.Run("managed widget returns false", func(t *testing.T) {
		w := dashboards.ApiTimeseriesWidgetAsApiWidget(
			dashboards.NewApiTimeseriesWidget("Time Series: Line"),
		)
		assert.False(t, isUnmanagedWidget(w))
	})

	t.Run("nil instance returns false", func(t *testing.T) {
		assert.False(t, isUnmanagedWidget(dashboards.ApiWidget{}))
	})
}

func TestMapAllWidgets_skipsNilResults(t *testing.T) {
	callCount := 0
	mapper := func(w dashboards.ApiWidget) (map[string]interface{}, error) {
		callCount++
		if callCount == 2 {
			return nil, nil
		}
		return map[string]interface{}{"i": callCount}, nil
	}

	widgets := []dashboards.ApiWidget{{}, {}, {}}
	result, err := mapAllWidgets(widgets, mapper)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func newManagedWidget(title string) dashboards.ApiWidget {
	w := dashboards.NewApiTimeseriesWidget("Time Series: Line")
	w.SetTitle(title)
	return dashboards.ApiTimeseriesWidgetAsApiWidget(w)
}

func newUnmanagedWidget() dashboards.ApiWidget {
	return dashboards.ApiWidget{
		ApiNumbersCardWidget: &dashboards.ApiNumbersCardWidget{},
	}
}

func TestMergeUnmanagedWidgets_preservesOrder(t *testing.T) {
	mA := newManagedWidget("A")
	mB := newManagedWidget("B")
	uX := newUnmanagedWidget()
	uY := newUnmanagedWidget()

	current := []dashboards.ApiWidget{mA, uX, mB, uY}
	config := []dashboards.ApiWidget{newManagedWidget("A'"), newManagedWidget("B'")}

	merged := mergeUnmanagedWidgets(config, current)
	require.Len(t, merged, 4)
	assert.Equal(t, "A'", merged[0].ApiTimeseriesWidget.GetTitle())
	assert.NotNil(t, merged[1].ApiNumbersCardWidget)
	assert.Equal(t, "B'", merged[2].ApiTimeseriesWidget.GetTitle())
	assert.NotNil(t, merged[3].ApiNumbersCardWidget)
}

func TestMergeUnmanagedWidgets_deleteManaged(t *testing.T) {
	mA := newManagedWidget("A")
	mB := newManagedWidget("B")
	uX := newUnmanagedWidget()

	current := []dashboards.ApiWidget{mA, uX, mB}
	config := []dashboards.ApiWidget{newManagedWidget("A'")}

	merged := mergeUnmanagedWidgets(config, current)
	require.Len(t, merged, 2)
	assert.Equal(t, "A'", merged[0].ApiTimeseriesWidget.GetTitle())
	assert.NotNil(t, merged[1].ApiNumbersCardWidget)
}

func TestMergeUnmanagedWidgets_addManaged(t *testing.T) {
	mA := newManagedWidget("A")
	uX := newUnmanagedWidget()

	current := []dashboards.ApiWidget{mA, uX}
	config := []dashboards.ApiWidget{newManagedWidget("A'"), newManagedWidget("NEW")}

	merged := mergeUnmanagedWidgets(config, current)
	require.Len(t, merged, 3)
	assert.Equal(t, "A'", merged[0].ApiTimeseriesWidget.GetTitle())
	assert.NotNil(t, merged[1].ApiNumbersCardWidget)
	assert.Equal(t, "NEW", merged[2].ApiTimeseriesWidget.GetTitle())
}

func TestMergeUnmanagedWidgets_emptyConfig(t *testing.T) {
	mA := newManagedWidget("A")
	uX := newUnmanagedWidget()
	uY := newUnmanagedWidget()

	current := []dashboards.ApiWidget{uX, mA, uY}
	merged := mergeUnmanagedWidgets(nil, current)
	require.Len(t, merged, 2)
	assert.NotNil(t, merged[0].ApiNumbersCardWidget)
	assert.NotNil(t, merged[1].ApiNumbersCardWidget)
}

func TestMergeUnmanagedWidgets_noUnmanaged(t *testing.T) {
	mA := newManagedWidget("A")
	mB := newManagedWidget("B")

	current := []dashboards.ApiWidget{mA, mB}
	config := []dashboards.ApiWidget{newManagedWidget("A'"), newManagedWidget("B'")}

	merged := mergeUnmanagedWidgets(config, current)
	require.Len(t, merged, 2)
	assert.Equal(t, "A'", merged[0].ApiTimeseriesWidget.GetTitle())
	assert.Equal(t, "B'", merged[1].ApiTimeseriesWidget.GetTitle())
}
