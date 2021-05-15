package analyzer

import (
	"testing"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/csvparser"
	"github.com/stretchr/testify/assert"
)

func TestSplitByDay_ConvertsToDay(t *testing.T) {
	in, err := AggregateIntoHourWindows([]csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(
			time.Date(2020, 01, 01, 12, 00, 00, 0, time.UTC),
			1.0)})
	assert.NoError(t, err)

	got, err := SplitByDay(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, time.Date(2020, 01, 01, 00, 00, 00, 0, time.UTC), got[0].Day)
	assert.Equal(t, 1.0, got[0].UsageKwh)
	assert.Equal(t, []UsageHour(in), got[0].DataPoints)
}

func TestSplitByDay_MultipleHours_AggregatesIntoSingleElement(t *testing.T) {
	in, err := AggregateIntoHourWindows([]csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(
			time.Date(2020, 01, 01, 12, 00, 00, 0, time.UTC),
			1.0),
		csvparser.NewRowWith15MinuteDuration(
			time.Date(2020, 01, 01, 13, 00, 00, 0, time.UTC),
			2.0)})
	assert.NoError(t, err)

	got, err := SplitByDay(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, time.Date(2020, 01, 01, 00, 00, 00, 0, time.UTC), got[0].Day)
	assert.Equal(t, 3.0, got[0].UsageKwh)
	assert.Equal(t, []UsageHour(in), got[0].DataPoints)
}

func TestSplitByDay_DifferentDays_AggregatesSeparately(t *testing.T) {
	in, err := AggregateIntoHourWindows([]csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(
			time.Date(2020, 01, 01, 12, 00, 00, 0, time.UTC),
			3.0),
		csvparser.NewRowWith15MinuteDuration(
			time.Date(2020, 01, 02, 13, 00, 00, 0, time.UTC),
			4.0)})
	assert.NoError(t, err)

	got, err := SplitByDay(in)
	assert.NoError(t, err)

	assert.Len(t, got, 2)
	assert.Equal(t, 3.0, got[0].UsageKwh)
	assert.Equal(t, 4.0, got[1].UsageKwh)
}
