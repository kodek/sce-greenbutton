package costcalculator

import (
	"github.com/kodek/sce-greenbutton/pkg/analyzer"
)

const TIER1_PER_KWH = 0.23
const TIER2_PER_KWH = 0.30
const TIER3_PER_KWH = 0.37

func CalculateDomesticCost(days []analyzer.UsageDay) float64 {
	baselineAllocation := baselineAllocationForDays(days)

	usage := 0.0
	for _, d := range days {
		usage += d.UsageKwh
	}

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

func baselineAllocationForDays(days []analyzer.UsageDay) float64 {
	totalAllowance := 0.0
	for _, d := range days {
		totalAllowance += GetDailyAllocation(d.Day)
	}
	return totalAllowance
}
