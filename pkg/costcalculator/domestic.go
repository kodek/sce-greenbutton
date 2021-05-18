package costcalculator

import (
	"math"

	"github.com/kodek/sce-greenbutton/pkg/analyzer"
)

const Tier1KWhCost = 0.06181 + 0.09545
const Tier2KWhCost = 0.10511 + 0.09545
const Tier3KWhCost = 0.15525 + 0.09545

const DomesticMinDailyCharge = 0.35
const DomesticDailyCharge = 0.031

type DomesticBreakdown struct {
	Days                  int
	BaselineAllocationKwh float64

	NemCost      float64
	MinCharges   float64
	DailyCharges float64

	UsageKwh      float64
	Tier1UsageKwh float64
	Tier2UsageKwh float64
	Tier3UsageKwh float64
}

func CalculateDomesticForDays(days []analyzer.UsageDay) DomesticBreakdown {
	out := DomesticBreakdown{
		Days:                  len(days),
		BaselineAllocationKwh: baselineAllocationForDays(days),
	}

	for _, d := range days {
		out.UsageKwh += d.UsageKwh
	}
	rebalanceTiers(&out)

	out.NemCost = out.Tier1UsageKwh*Tier1KWhCost + out.Tier2UsageKwh*Tier2KWhCost + out.Tier3UsageKwh*Tier3KWhCost
	out.DailyCharges = DomesticDailyCharge * float64(out.Days)

	minCharge := DomesticMinDailyCharge * float64(out.Days)
	if minCharge > out.NemCost {
		out.MinCharges = minCharge - math.Max(0.0, out.NemCost)
	}

	return out
}

func rebalanceTiers(d *DomesticBreakdown) {
	baseline := d.BaselineAllocationKwh
	remainingAbsUsage := d.UsageKwh

	if remainingAbsUsage >= 4*baseline {
		t3Usage := remainingAbsUsage - 4*baseline
		remainingAbsUsage -= t3Usage
		d.Tier3UsageKwh = math.Copysign(t3Usage, d.UsageKwh)
	}
	if remainingAbsUsage >= 1*baseline {
		t2Usage := remainingAbsUsage - baseline
		remainingAbsUsage -= t2Usage
		d.Tier2UsageKwh = math.Copysign(t2Usage, d.UsageKwh)
	}
	d.Tier1UsageKwh = math.Copysign(remainingAbsUsage, d.UsageKwh)
}

func baselineAllocationForDays(days []analyzer.UsageDay) float64 {
	totalAllowance := 0.0
	for _, d := range days {
		totalAllowance += GetDailyAllocation(d.Day)
	}
	return totalAllowance
}
