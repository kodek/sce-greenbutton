package analyzer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSplitByDay_ConvertsToDay(t *testing.T) {
	in := []UsageHour{
		{
			StartTime: time.Date(2020, 01, 01, 12, 00, 00, 0, time.UTC),
			UsageKwh:  1,
		},
	}

	got, err := SplitByDay(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, time.Date(2020, 01, 01, 00, 00, 00, 0, time.UTC), got[0].Day)
	assert.Equal(t, 1.0, got[0].UsageKwh)
	assert.Equal(t, in, got[0].DataPoints)
}

func TestSplitByDay_MultipleHours_SumsUsage(t *testing.T) {
	in := []UsageHour{
		{
			StartTime: time.Date(2020, 01, 01, 12, 00, 00, 0, time.UTC),
			UsageKwh:  1,
		},
		{
			StartTime: time.Date(2020, 01, 01, 13, 00, 00, 0, time.UTC),
			UsageKwh:  2,
		},
	}

	got, err := SplitByDay(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, time.Date(2020, 01, 01, 00, 00, 00, 0, time.UTC), got[0].Day)
	assert.Equal(t, 3.0, got[0].UsageKwh)
	assert.Equal(t, in, got[0].DataPoints)
}

func TestSplitByDay_DifferentDays_AggregatesSeparately(t *testing.T) {
	in := []UsageHour{
		{
			StartTime: time.Date(2020, 01, 01, 12, 00, 00, 0, time.UTC),
			UsageKwh:  3,
		},
		{
			StartTime: time.Date(2020, 01, 02, 13, 00, 00, 0, time.UTC),
			UsageKwh:  4,
		},
	}

	got, err := SplitByDay(in)
	assert.NoError(t, err)

	assert.Len(t, got, 2)
	assert.Equal(t, 3.0, got[0].UsageKwh)
	assert.Equal(t, 4.0, got[1].UsageKwh)
}
