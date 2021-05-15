package analyzer

import (
	"fmt"
	"time"
)

var stateHolidays = []string{
	// From https://www.sos.ca.gov/state-holidays
	// TODO: SCE may shift some holidays from Sunday to Monday. Confirm if this is the case.
	"2021-01-01",
	"2021-01-18",
	"2021-02-15",
	"2021-03-31",
	"2021-05-31",
	"2021-07-05",
	"2021-09-06",
	"2021-11-11",
	"2021-11-25",
	"2021-11-26",
	"2021-12-25",
}

func IsHoliday(t time.Time) bool {
	for _, hStr := range stateHolidays {
		h, err := time.ParseInLocation("2006-01-02", hStr, time.UTC)
		if err != nil {
			panic(fmt.Sprintf("Unable to parse date %s", hStr))
		}
		if t.Month() == h.Month() && t.Day() == h.Day() {
			return true
		}
	}
	return false
}
