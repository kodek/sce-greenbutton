package analyzer

import (
	"testing"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/csvparser"
	"github.com/stretchr/testify/assert"
)

var now = time.Date(2020, 1, 1, 1, 0, 0, 0, time.Local)

func TestConvertsDataPointToUsageHour(t *testing.T) {
	parsed := csvparser.CsvFile{
		csvparser.CsvRow{
			StartTime:      now,
			EndTime:        now.Add(15 * time.Minute),
			UsageKwh:       123,
			ReadingQuality: "",
		},
	}

	got, err := AggregateIntoHourWindows(parsed)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, now, got[0].StartTime())
	assert.Equal(t, 123.0, got[0].UsageKwh())
}

func TestTwoDataPointsInOneHourAggregated(t *testing.T) {
	parsed := csvparser.CsvFile{
		csvparser.NewRowWith15MinuteDuration(now, 1),
		csvparser.NewRowWith15MinuteDuration(now.Add(15*time.Minute), 2),
	}

	got, err := AggregateIntoHourWindows(parsed)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(got))
	assert.Equal(t, now, got[0].StartTime())
	assert.Equal(t, 3.0, got[0].UsageKwh())
}
func TestWindowOverlapsHourFails(t *testing.T) {
	parsed := csvparser.CsvFile{
		csvparser.CsvRow{
			StartTime:      time.Date(2020, 01, 01, 01, 50, 0, 0, time.UTC),
			EndTime:        time.Date(2020, 01, 01, 02, 10, 0, 0, time.UTC),
			UsageKwh:       1,
			ReadingQuality: "",
		},
	}

	_, err := AggregateIntoHourWindows(parsed)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "within the same hour")
	}
}

func TestUsageHour_EndTime_ReturnsEndTimeFromLastPoint(t *testing.T) {
	parsed := csvparser.CsvFile{
		csvparser.NewRowWith15MinuteDuration(now, 1),
		csvparser.NewRowWith15MinuteDuration(now.Add(15*time.Minute), 2),
	}

	got, err := AggregateIntoHourWindows(parsed)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(got))
	assert.Equal(t, now.Add(30*time.Minute), got[0].EndTime())
}
