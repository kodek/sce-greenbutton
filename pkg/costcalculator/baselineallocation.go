package costcalculator

import (
	"flag"
	"time"
)

// From https://www.sce.com/residential/rates/Standard-Residential-Rate-Plan
const SIMI_SUMMER_DAILY_ALLOCATION = 16.5
const SIMI_WINTER_DAILY_ALLOCATION = 12.3

const MedicalBaselineAllocation = 16.5

var useMedicalBaseline = flag.Bool("use_medical_baseline", true, "Add medical baseline")

// GetDailyAllocation returns the daily allocation based on the following data:
func GetDailyAllocation(t time.Time) float64 {
	medicalOffset := 0.0
	if *useMedicalBaseline {
		medicalOffset = MedicalBaselineAllocation
	}
	if isSummerMonth(t.Month()) {
		return SIMI_SUMMER_DAILY_ALLOCATION + medicalOffset
	} else {
		return SIMI_WINTER_DAILY_ALLOCATION + medicalOffset
	}
}

// Returns true if it's a summer month.
// Winter allocation: October through May
// Summer allocation: June through September
func isSummerMonth(m time.Month) bool {
	return m >= 6 && m <= 9
}
