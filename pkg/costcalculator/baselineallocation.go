package costcalculator

import "time"

const SIMI_SUMMER_DAILY_ALLOCATION = 13.8
const SIMI_WINTER_DAILY_ALLOCATION = 10.6

// Returns the daily allocation based on the following data:
func getDailyAllocation(t time.Time) float64 {
	if isSummerMonth(t.Month()) {
		return SIMI_SUMMER_DAILY_ALLOCATION
	} else {
		return SIMI_WINTER_DAILY_ALLOCATION
	}
}

// Returns true if it's a summer month.
// Winter allocation: October through May
// Summer allocation: June through September
func isSummerMonth(m time.Month) bool {
	return m >= 6 && m <= 9
}
