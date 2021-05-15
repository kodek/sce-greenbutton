package costcalculator

import (
	"time"

	"github.com/kodek/sce-greenbutton/pkg/analyzer"
)

const SUMMER_SUPER_OFF_PER_KWH = 0.12
const SUMMER_OFF_PER_KWH = 0.28
const SUMMER_ON_PER_KWH = 0.48

const WINTER_SUPER_OFF_PER_KWH = 0.13
const WINTER_OFF_PER_KWH = 0.27
const WINTER_ON_PER_KWH = 0.36

//const DAILY_BASIC_CHARGE = 0.03
//const MINIMUM_DAILY_CHARGE = 0.34
const BASELINE_CREDIT_PER_KWH = 0.08

//func CalculateTouDANonBypassableChargesForMonth(m analyzer.UsageMonth) float64 {
//	cumulativeCost := 0.0
//	for _, d := range m.UsageDays {
//		cumulativeCost += 0.03
//		if d.UsageKwh < MINIMUM_DAILY_CHARGE {
//			dailyCost := touDACostForDay(d)
//			if dailyCost > 0 {
//				cumulativeCost += MINIMUM_DAILY_CHARGE - dailyCost
//			} else {
//				cumulativeCost += MINIMUM_DAILY_CHARGE
//			}
//		}
//	}
//	return cumulativeCost
//}

func CalculateTouDACostForMonth(m analyzer.UsageMonth) float64 {
	cumulativeCost := 0.0
	cumulativeKwh := 0.0
	for _, d := range m.UsageDays {
		cumulativeCost += touDACostForDay(d)
		cumulativeKwh += d.UsageKwh
	}

	baselineAllocation := baselineAllocationForMonth(m.Month)
	if cumulativeKwh > baselineAllocation {
		cumulativeCost -= baselineAllocation * BASELINE_CREDIT_PER_KWH
	} else {
		cumulativeCost -= cumulativeKwh * BASELINE_CREDIT_PER_KWH
	}
	return cumulativeCost
}

func touDACostForDay(d analyzer.UsageDay) float64 {
	cumulativeCost := 0.0
	for _, h := range d.DataPoints {
		cumulativeCost += h.UsageKwh() * calculateTouDARateForHour(h)
	}
	return cumulativeCost
}

func calculateTouDARateForHour(hour analyzer.UsageHour) float64 {
	t := hour.StartTime
	if isOnPeak(t) {
		if isSummerMonth(t.Month()) {
			return SUMMER_ON_PER_KWH
		} else {
			return WINTER_ON_PER_KWH
		}
	} else if isOffPeak(t) {
		if isSummerMonth(t.Month()) {
			return SUMMER_OFF_PER_KWH
		} else {
			return WINTER_OFF_PER_KWH
		}
	} else if isSuperOffPeak(t) {
		if isSummerMonth(t.Month()) {
			return SUMMER_SUPER_OFF_PER_KWH
		} else {
			return WINTER_SUPER_OFF_PER_KWH
		}
	} else {
		panic("Unexpected")
	}
}

func isOnPeak(t time.Time) bool {
	if isWeekend(t) {
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
