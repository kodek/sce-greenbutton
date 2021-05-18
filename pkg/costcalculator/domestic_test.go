package costcalculator

import (
	"testing"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/analyzer"
	"github.com/stretchr/testify/assert"
)

func TestCalculateDomesticForDays(t *testing.T) {
	winterDay := time.Date(2020, 01, 01, 00, 00, 00, 00, time.UTC)
	dailyAllocation := GetDailyAllocation(winterDay)
	assert.Equal(t, 28.8, dailyAllocation)

	in := []analyzer.UsageDay{
		{
			Day:      winterDay.Add(0 * 24 * time.Hour),
			UsageKwh: 100,
		},
		{
			Day:      winterDay.Add(1 * 24 * time.Hour),
			UsageKwh: 200,
		},
		{
			Day:      winterDay.Add(2 * 24 * time.Hour),
			UsageKwh: 1000,
		},
	}

	actual := CalculateDomesticForDays(in)

	assert.Equal(t, 3, actual.Days)
	assert.Equal(t, 1300.0, actual.UsageKwh)
	assert.Equal(t, 1300.0, actual.UsageKwh)
	assert.InEpsilon(t, 86.4, actual.Tier1UsageKwh, 0.01)
	assert.InEpsilon(t, 259.6, actual.Tier2UsageKwh, 0.01)
	assert.InEpsilon(t, 954.4, actual.Tier3UsageKwh, 0.01)
	assert.InEpsilon(t, 0.093, actual.DailyCharges, 0.01)
	assert.Equal(t, 0.0, actual.MinCharges)
}

func TestCalculateDomesticForDays_MinCharge(t *testing.T) {
	winterDay := time.Date(2020, 01, 01, 00, 00, 00, 00, time.UTC)

	in := []analyzer.UsageDay{
		{
			Day:      winterDay.Add(0 * 24 * time.Hour),
			UsageKwh: 1,
		},
	}

	actual := CalculateDomesticForDays(in)

	assert.Greater(t, actual.MinCharges, 0.0)
}
