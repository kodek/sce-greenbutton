package analyzer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSplitByMonth_ConvertsDayToMonth(t *testing.T) {
	in := []UsageDay{
		{
			Day:      time.Date(2020, 01, 01, 00, 00, 00, 0, time.UTC),
			UsageKwh: 1,
		},
	}

	got, err := SplitByMonth(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, time.Date(2020, 01, 01, 00, 00, 00, 0, time.UTC), got[0].Month)
	assert.Equal(t, 1.0, got[0].UsageKwh)
	assert.Equal(t, in, got[0].UsageDays)
}

func TestSplitByMonth_MultipleDays_AggregatesIntoSingleElement(t *testing.T) {
	in := []UsageDay{
		{
			Day:      time.Date(2020, 01, 01, 0, 00, 00, 0, time.UTC),
			UsageKwh: 1,
		},
		{
			Day:      time.Date(2020, 01, 02, 0, 00, 00, 0, time.UTC),
			UsageKwh: 2,
		},
	}

	got, err := SplitByMonth(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, time.Date(2020, 01, 01, 00, 00, 00, 0, time.UTC), got[0].Month)
	assert.Equal(t, 3.0, got[0].UsageKwh)
	assert.Equal(t, in, got[0].UsageDays)
}

func TestSplitByMonth_DifferentMonths_AggregatesSeparately(t *testing.T) {
	in := []UsageDay{
		{
			Day:      time.Date(2020, 01, 01, 00, 00, 00, 0, time.UTC),
			UsageKwh: 3,
		},
		{
			Day:      time.Date(2020, 02, 01, 00, 00, 00, 0, time.UTC),
			UsageKwh: 4,
		},
	}

	got, err := SplitByMonth(in)
	assert.NoError(t, err)

	assert.Len(t, got, 2)
	assert.Equal(t, 3.0, got[0].UsageKwh)
	assert.Equal(t, 4.0, got[1].UsageKwh)
}
