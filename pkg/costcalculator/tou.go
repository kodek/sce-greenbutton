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
	SummerOnPeak

	WinterSuperOffPeak
	WinterOffPeak
	WinterOnPeak
)

func (cost *CostPeriod) Cost() float64 {
	switch *cost {
	case SummerOffPeak:
		return 0.34
	case SummerOnPeak:
		return 0.61
	case SummerSuperOffPeak:
		return 0.16
	case WinterOffPeak:
		return 0.30
	case WinterOnPeak:
		return 0.40
	case WinterSuperOffPeak:
		return 0.16
	}
	panic("unexpected")
}

func (cost *CostPeriod) Name() string {
	switch *cost {
	case SummerOffPeak:
		return "Off-Peak (Summer)"
	case SummerOnPeak:
		return "On-Peak (Summer)"
	case SummerSuperOffPeak:
		return "Super Off-Peak (Summer)"
	case WinterOffPeak:
		return "Off-Peak (Winter)"
	case WinterOnPeak:
		return "On-Peak (Winter)"
	case WinterSuperOffPeak:
		return "Super Off-Peak (Winter)"
	}
	panic("unexpected")
}

const BaselineCreditPerKwh = -0.07848
const NonBypassableChargePerKwh = 0.01362

func CalculateTouDACostForDays(days []analyzer.UsageDay) TouBillSummary {

	bucket := TouBillSummary{
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

type TouBillSummary struct {
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
		total += usage * period.Cost()
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
	return b.energyImported * NonBypassableChargePerKwh
}

func (b *TouBillSummary) MaxBaselineAllowance() float64 {
	total := 0.0
	for _, p := range b.days {
		total += GetDailyAllocation(p.Day)
	}
	return total
}

func (b *TouBillSummary) BaselineCredit() float64 {
	actualUsage := b.NetEnergyUsage()
	maxBaseline := b.MaxBaselineAllowance()

	absActualUsage := math.Abs(actualUsage)
	absAllowance := math.Min(absActualUsage, maxBaseline)

	return math.Copysign(absAllowance, actualUsage) * BaselineCreditPerKwh
}

func (b *TouBillSummary) AverageDailyUsage() float64 {
	return float64(b.NetEnergyUsage()) / float64(len(b.days))
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
		period := calculateTouDARateForHour(h.StartTime)
		b.usageKwhByPeriod[period] += h.UsageKwh()
		b.hoursByPeriod[period] += 1
		if h.UsageKwh() > 0 {
			b.energyImported += h.UsageKwh()
		} else {
			b.energyExported += h.UsageKwh()
		}
	}
}

func calculateTouDARateForHour(t time.Time) CostPeriod {
	if isOnPeak(t) {
		if isSummerMonth(t.Month()) {
			return SummerOnPeak
		} else {
			return WinterOnPeak
		}
	} else if isOffPeak(t) {
		if isSummerMonth(t.Month()) {
			return SummerOffPeak
		} else {
			return WinterOffPeak
		}
	} else if isSuperOffPeak(t) {
		if isSummerMonth(t.Month()) {
			return SummerSuperOffPeak
		} else {
			return WinterSuperOffPeak
		}
	} else {
		panic("Unexpected")
	}
}

func isOnPeak(t time.Time) bool {
	if isWeekend(t) || analyzer.IsHoliday(t) {
		return false
	}
	h := t.Hour()
	return h >= 14 && h < 20
}

func isWeekend(t time.Time) bool {
	weekDay := t.Weekday()
	return weekDay == time.Saturday || weekDay == time.Sunday
}

func isOffPeak(t time.Time) bool {
	h := t.Hour()
	return h >= 8 && h < 22 && !isOnPeak(t)
}

func isSuperOffPeak(t time.Time) bool {
	return !isOnPeak(t) && !isOffPeak(t)
}
