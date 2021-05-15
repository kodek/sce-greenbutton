package costcalculator

import (
	"time"
)

type TouDAPlan struct{}

func NewTouDAPlan() TouPlan {
	return &TouDAPlan{}
}

func (p *TouDAPlan) Name() string { return "TOU-D-A" }

func (p *TouDAPlan) Cost(period CostPeriod) float64 {
	switch period {
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

func (p *TouDAPlan) IsOnPeak(t time.Time) bool {
	if !isWeekday(t) {
		return false
	}
	h := t.Hour()
	return h >= 14 && h < 20
}

func (p *TouDAPlan) IsMidPeak(_ time.Time) bool {
	return false
}
func (p *TouDAPlan) IsOffPeak(t time.Time) bool {
	h := t.Hour()
	return h >= 8 && h < 22 && !p.IsOnPeak(t)
}

func (p *TouDAPlan) IsSuperOffPeak(t time.Time) bool {
	return !p.IsOnPeak(t) && !p.IsOffPeak(t)
}

func (p *TouDAPlan) DailyBasicCharge() float64 {
	return 0.031
}

func (p *TouDAPlan) MinimumDailyCharge() float64 {
	return 0.35
}

func (p *TouDAPlan) HasBaselineAllocation() bool {
	return true
}
