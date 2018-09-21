package costcalculator

import (
	"time"

	"github.com/kodek/sce-greenbutton/pkg/analyzer"
)

const TIER1_PER_KWH = 0.18
const TIER2_PER_KWH = 0.25
const TIER3_PER_KWH = 0.35

func CalculateDomesticCost(month analyzer.UsageMonth) float64 {
	baselineAllocation := baselineAllocationForMonth(month.Month)

	usage := month.UsageKwh
	cumulativeCost := 0.0

	if usage <= baselineAllocation {
		return TIER1_PER_KWH * usage
	} else {
		cumulativeCost += TIER1_PER_KWH * baselineAllocation
		usage -= baselineAllocation
	}
	if usage <= 3*baselineAllocation {
		cumulativeCost += TIER2_PER_KWH * usage
		return cumulativeCost
	} else {
		cumulativeCost += TIER2_PER_KWH * 3 * baselineAllocation
		usage -= 3 * baselineAllocation
	}
	if usage < 0 {
		panic("Usage should be positive if we're in tier 3")
	}
	cumulativeCost += TIER3_PER_KWH * usage
	return cumulativeCost
}

func baselineAllocationForMonth(t time.Time) float64 {
	daysInMonth := daysInMonth(t)
	return float64(daysInMonth) * getDailyAllocation(t)
}

func daysInMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
