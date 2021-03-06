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
	assert.Equal(t, in, got[0].DataPoints)
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
	assert.Equal(t, in, got[0].DataPoints)
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

// Prevents nondeterministic map iteration behavior from shuffling dates.
func TestSplitByDay_DateOrderingPreserved(t *testing.T) {
	in, err := AggregateIntoHourWindows([]csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(
			time.Date(2021, 01, 01, 12, 00, 00, 0, time.UTC),
			3.0),
		csvparser.NewRowWith15MinuteDuration(
			time.Date(2020, 01, 02, 13, 00, 00, 0, time.UTC),
			4.0)})
	assert.NoError(t, err)

	got, err := SplitByDay(in)
	assert.NoError(t, err)

	assert.Len(t, got, 2)
	assert.Equal(t, in[0].StartTime().Year(), got[0].Day.Year())
	assert.Equal(t, in[1].StartTime().Year(), got[1].Day.Year())
}

func TestUsageDay_EndTime_ReturnsEndFromLastDataPoint(t *testing.T) {
	in, err := AggregateIntoHourWindows([]csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(
			time.Date(2020, 01, 01, 12, 00, 00, 0, time.UTC),
			3.0),
		csvparser.NewRowWith15MinuteDuration(
			time.Date(2020, 01, 01, 15, 00, 00, 0, time.UTC),
			4.0)})
	assert.NoError(t, err)

	got, err := SplitByDay(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, time.Date(2020, 01, 01, 15, 15, 0, 0, time.UTC), got[0].EndTime())
}
