package costcalculator

import (
	"time"
)

type TouDPrime struct{}

func NewTouDPrime() TouPlan {
	return &TouDPrime{}
}

func (p *TouDPrime) Name() string { return "TOU-D-PRIME" }

func (p *TouDPrime) Cost(period CostPeriod) float64 {
	switch period {
	case SummerOffPeak:
		return 0.17
	case SummerMidPeak:
		return 0.33
	case SummerOnPeak:
		return 0.44
	case SummerSuperOffPeak:
		panic("should never happen")
	case WinterOffPeak:
		return 0.16
	case WinterMidPeak:
		return 0.41
	case WinterOnPeak:
		panic("should never happen")
	case WinterSuperOffPeak:
		return 0.16
	}
	panic("unexpected")
}

// IsOnPeak is only true on summers, weekdays, between 4-9 pm.
func (p *TouDPrime) IsOnPeak(t time.Time) bool {
	if !isSummerMonth(t.Month()) {
		return false
	}
	if !isWeekday(t) {
		return false
	}
	return is4to9(t)
}

// IsMidPeak is true between 4-9pm all days except on summer weeekdays
func (p *TouDPrime) IsMidPeak(t time.Time) bool {
	if !is4to9(t) {
		return false
	}
	return !(isWeekday(t) && isSummerMonth(t.Month()))
}

func is4to9(t time.Time) bool {
	h := t.Hour()
	return h >= 16 && h < 21
}

// IsOffPeak is true outside 4-9pm except for when it's super-off peak
func (p *TouDPrime) IsOffPeak(t time.Time) bool {
	if is4to9(t) {
		return false
	}
	return !p.IsSuperOffPeak(t)
}

// IsSuperOffPeak is true on winter between 8am-4pm.
func (p *TouDPrime) IsSuperOffPeak(t time.Time) bool {
	if isSummerMonth(t.Month()) {
		return false
	}
	h := t.Hour()
	return h >= 8 && h < 16
}

func (p *TouDPrime) DailyBasicCharge() float64 {
	return 0.40
}

func (p *TouDPrime) MinimumDailyCharge() float64 {
	return 0
}

func (p *TouDPrime) HasBaselineAllocation() bool {
	return false
}
