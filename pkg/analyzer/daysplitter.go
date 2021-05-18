package analyzer

import "time"

type UsageDay struct {
	Day        time.Time
	DataPoints []UsageHour
	UsageKwh   float64
}

func (d *UsageDay) EndTime() time.Time {
	return d.DataPoints[len(d.DataPoints)-1].EndTime()
}

func SplitByDay(in []UsageHour) ([]UsageDay, error) {
	yearMonthDayMap := make(map[time.Time]*UsageDay)
	// Keep track of the order in which keys were added.
	sortedKeys := make([]time.Time, 0)

	for _, hr := range in {
		dayBucket := truncateToDay(hr.StartTime())
		usageDay, ok := yearMonthDayMap[dayBucket]
		if !ok {
			sortedKeys = append(sortedKeys, dayBucket)
			usageDay = &UsageDay{
				Day:        dayBucket,
				DataPoints: make([]UsageHour, 0),
				UsageKwh:   0}
			yearMonthDayMap[dayBucket] = usageDay
		}
		// Append the data to the day
		usageDay.DataPoints = append(usageDay.DataPoints, hr)
		usageDay.UsageKwh += hr.UsageKwh()
	}

	values := make([]UsageDay, 0)
	for _, k := range sortedKeys {
		values = append(values, *yearMonthDayMap[k])
	}

	// Validate data points
	for _, v := range values {
		expectedYear := v.Day.Year()
		expectedMonth := v.Day.Month()
		expectedDay := v.Day.Day()
		for _, p := range v.DataPoints {
			if p.StartTime().Day() != expectedDay {
				panic("Wrong day within UsageDay")
			}
			if p.StartTime().Month() != expectedMonth {
				panic("Wrong month within UsageDay")
			}
			if p.StartTime().Year() != expectedYear {
				panic("Wrong year within UsageDay")
			}
		}
	}
	return values, nil
}

func truncateToDay(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, in.Location())
}
