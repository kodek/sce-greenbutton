package costcalculator

import (
	"time"
)

type TouD58 struct{}

func NewTouD58() TouPlan {
	return &TouD58{}
}

func (p *TouD58) Name() string { return "TOU-D-5-8PM" }

func (p *TouD58) Cost(period CostPeriod) float64 {
	switch period {
	case SummerOffPeak:
		return 0.27
	case SummerMidPeak:
		return 0.40
	case SummerOnPeak:
		return 0.54
	case SummerSuperOffPeak:
		panic("should never happen")
	case WinterOffPeak:
		return 0.29
	case WinterMidPeak:
		return 0.44
	case WinterOnPeak:
		panic("should never happen")
	case WinterSuperOffPeak:
		return 0.25
	}
	panic("unexpected")
}

// IsOnPeak is true on summer weekdays between 5-8
func (p *TouD58) IsOnPeak(t time.Time) bool {
	if !isSummerMonth(t.Month()) {
		return false
	}
	if !isWeekday(t) {
		return false
	}
	return is5to8(t)
}

func is5to8(t time.Time) bool {
	h := t.Hour()
	return h >= 17 && h < 20
}

// IsMidPeak is true between 5-8pm all days except on summer weeekdays
func (p *TouD58) IsMidPeak(t time.Time) bool {
	if !is5to8(t) {
		return false
	}
	return !(isWeekday(t) && isSummerMonth(t.Month()))
}

// IsOffPeak is true outside 5-8pm except for when it's super-off peak
func (p *TouD58) IsOffPeak(t time.Time) bool {
	if is5to8(t) {
		return false
	}
	return !p.IsSuperOffPeak(t)
}

// IsSuperOffPeak is true on winter between 8am-5pm.
func (p *TouD58) IsSuperOffPeak(t time.Time) bool {
	if isSummerMonth(t.Month()) {
		return false
	}
	h := t.Hour()
	return h >= 8 && h < 17
}

func (p *TouD58) DailyBasicCharge() float64 {
	return 0.03
}

func (p *TouD58) MinimumDailyCharge() float64 {
	return 0.35
}

func (p *TouD58) HasBaselineAllocation() bool {
	return true
}
