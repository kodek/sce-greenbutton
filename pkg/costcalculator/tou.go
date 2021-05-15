package costcalculator

import (
	"math"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/analyzer"
)

type CostPeriod float64

const (
	SummerSuperOffPeak CostPeriod = iota
	SummerOffPeak
	SummerMidPeak
	SummerOnPeak

	WinterSuperOffPeak
	WinterOffPeak
	WinterMidPeak
	WinterOnPeak
)

func (cost *CostPeriod) Name() string {
	switch *cost {
	case SummerOffPeak:
		return "Off-Peak (Summer)"
	case SummerMidPeak:
		return "Mid-Peak (Summer)"
	case SummerOnPeak:
		return "On-Peak (Summer)"
	case SummerSuperOffPeak:
		return "Super Off-Peak (Summer)"
	case WinterOffPeak:
		return "Off-Peak (Winter)"
	case WinterMidPeak:
		return "Mid-Peak (Winter)"
	case WinterOnPeak:
		return "On-Peak (Winter)"
	case WinterSuperOffPeak:
		return "Super Off-Peak (Winter)"
	}
	panic("unexpected")
}

const BaselineCreditPerKwh = -0.07848
const Nem2NonBypassableChargePerKwh = 0.01362
const StateTaxPerKwh = 0.00030

func CalculateTouDACostForDays(days []analyzer.UsageDay) TouBillSummary {
	return CalculateWithTouPlan(days, NewTouDAPlan())
}

func CalculateWithTouPlan(days []analyzer.UsageDay, plan TouPlan) TouBillSummary {

	bucket := TouBillSummary{
		touPlan:          plan,
		days:             days,
		usageKwhByPeriod: make(map[CostPeriod]float64),
		hoursByPeriod:    make(map[CostPeriod]int),
		energyImported:   0,

		weekdays: 0,
		weekends: 0,
		holidays: 0,
	}

	for _, d := range days {
		bucket.add(d)
	}

	return bucket
}

type TouPlan interface {
	Name() string
	Cost(period CostPeriod) float64
	IsOnPeak(t time.Time) bool
	IsMidPeak(t time.Time) bool
	IsOffPeak(t time.Time) bool
	IsSuperOffPeak(t time.Time) bool
	DailyBasicCharge() float64
	MinimumDailyCharge() float64
	HasBaselineAllocation() bool
}

type TouBillSummary struct {
	touPlan          TouPlan
	days             []analyzer.UsageDay
	usageKwhByPeriod map[CostPeriod]float64
	hoursByPeriod    map[CostPeriod]int
	energyExported   float64
	energyImported   float64

	weekdays int
	weekends int
	holidays int
}

func (b *TouBillSummary) NetMeteredCost() float64 {
	total := 0.0
	for period, usage := range b.usageKwhByPeriod {
		total += usage * b.touPlan.Cost(period)
	}
	return total
}

func (b *TouBillSummary) NetEnergyUsage() float64 {
	total := 0.0
	for _, usage := range b.usageKwhByPeriod {
		total += usage
	}
	return total
}

func (b *TouBillSummary) Taxes() float64 {
	usage := b.NetEnergyUsage()
	if usage > 0 {
		return usage * StateTaxPerKwh
	} else {
		return 0
	}

}

func (b *TouBillSummary) UsageByPeriod() map[CostPeriod]float64 {
	return b.usageKwhByPeriod
}

func (b *TouBillSummary) EnergyExported() float64 {
	return b.energyExported
}

func (b *TouBillSummary) EnergyImported() float64 {
	return b.energyImported
}

func (b *TouBillSummary) NonBypassableCharges() float64 {
	return b.energyImported * Nem2NonBypassableChargePerKwh
}

func (b *TouBillSummary) MaxBaselineAllowance() float64 {
	if !b.touPlan.HasBaselineAllocation() {
		return 0
	}
	total := 0.0
	for _, p := range b.days {
		total += GetDailyAllocation(p.Day)
	}
	return total
}

func (b *TouBillSummary) TotalBasicCharge() float64 {
	return float64(len(b.days)) * b.touPlan.DailyBasicCharge()
}

func (b *TouBillSummary) BaselineCredit() float64 {
	actualUsage := b.NetEnergyUsage()
	maxBaseline := b.MaxBaselineAllowance()

	absActualUsage := math.Abs(actualUsage)
	absAllowance := math.Min(absActualUsage, maxBaseline)

	return math.Copysign(absAllowance, actualUsage) * BaselineCreditPerKwh
}

func (b *TouBillSummary) AverageDailyUsage() float64 {
	return b.NetEnergyUsage() / float64(len(b.days))
}

func (b *TouBillSummary) add(d analyzer.UsageDay) {
	if isWeekend(d.Day) {
		b.weekends++
	}
	if analyzer.IsHoliday(d.Day) {
		b.holidays++
	}
	if !isWeekend(d.Day) && !analyzer.IsHoliday(d.Day) {
		b.weekdays++
	}
	for _, h := range d.DataPoints {
		period := calculateTouRateForHour(h.StartTime, b.touPlan)
		b.usageKwhByPeriod[period] += h.UsageKwh()
		b.hoursByPeriod[period] += 1
		if h.UsageKwh() > 0 {
			b.energyImported += h.UsageKwh()
		} else {
			b.energyExported += h.UsageKwh()
		}
	}
}

func calculateTouRateForHour(t time.Time, plan TouPlan) CostPeriod {
	if plan.IsOnPeak(t) {
		if isSummerMonth(t.Month()) {
			return SummerOnPeak
		} else {
			return WinterOnPeak
		}
	} else if plan.IsMidPeak(t) {
		if isSummerMonth(t.Month()) {
			return SummerMidPeak
		} else {
			return WinterMidPeak
		}
	} else if plan.IsOffPeak(t) {
		if isSummerMonth(t.Month()) {
			return SummerOffPeak
		} else {
			return WinterOffPeak
		}
	} else if plan.IsSuperOffPeak(t) {
		if isSummerMonth(t.Month()) {
			return SummerSuperOffPeak
		} else {
			return WinterSuperOffPeak
		}
	} else {
		panic("Unexpected")
	}
}

func isWeekend(t time.Time) bool {
	weekDay := t.Weekday()
	return weekDay == time.Saturday || weekDay == time.Sunday
}

func isWeekday(t time.Time) bool {
	return !isWeekend(t) && !analyzer.IsHoliday(t)
}
