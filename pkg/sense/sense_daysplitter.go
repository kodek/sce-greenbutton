// DO NOT SUBMIT: This is the same as pkg/analyzer/daysplitter.go, but the types are different. Refactor
package sense

import "time"

type Day struct {
	Day            time.Time
	DataPoints     []Snapshot
	ProductionKwh  float64
	ConsumptionKwh float64
}

// NetUsageKwh returns the energy consumed from the grid.
func (d *Day) NetUsageKwh() float64 {
	return d.ConsumptionKwh + d.ProductionKwh
}

func SplitByDay(in []Snapshot) ([]Day, error) {
	yearMonthDayMap := make(map[time.Time]*Day)
	// Keep track of the order in which keys were added.
	sortedKeys := make([]time.Time, 0)

	for _, hr := range in {
		dayBucket := truncateToDay(hr.DateTime)
		day, ok := yearMonthDayMap[dayBucket]
		if !ok {
			sortedKeys = append(sortedKeys, dayBucket)
			day = &Day{
				Day:            dayBucket,
				DataPoints:     make([]Snapshot, 0),
				ProductionKwh:  0,
				ConsumptionKwh: 0}
			yearMonthDayMap[dayBucket] = day
		}
		// Append the data to the day
		day.DataPoints = append(day.DataPoints, hr)
		day.ProductionKwh += *hr.ProductionKwh
		day.ConsumptionKwh += *hr.ConsumptionKwh
	}

	values := make([]Day, 0)
	for _, k := range sortedKeys {
		values = append(values, *yearMonthDayMap[k])
	}

	// Validate data points
	for _, v := range values {
		expectedYear := v.Day.Year()
		expectedMonth := v.Day.Month()
		expectedDay := v.Day.Day()
		for _, p := range v.DataPoints {
			if p.DateTime.Day() != expectedDay {
				panic("Wrong day within UsageDay")
			}
			if p.DateTime.Month() != expectedMonth {
				panic("Wrong month within UsageDay")
			}
			if p.DateTime.Year() != expectedYear {
				panic("Wrong year within UsageDay")
			}
		}
	}
	return values, nil
}

func truncateToDay(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, in.Location())
}
