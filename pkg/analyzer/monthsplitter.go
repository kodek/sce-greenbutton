package analyzer

import "time"

type UsageMonth struct {
	Month     time.Time
	UsageDays []UsageDay
	UsageKwh  float64
}

func (m *UsageMonth) AverageDailyUsageKwh() float64 {
	return m.UsageKwh / float64(len(m.UsageDays))
}

func SplitByMonth(in []UsageDay) ([]UsageMonth, error) {
	yearMonthMap := make(map[time.Time]*UsageMonth)
	sortedKeys := make([]time.Time, 0)

	for _, d := range in {
		monthBucket := truncateToMonth(d.Day)
		usageMonth, ok := yearMonthMap[monthBucket]
		if !ok {
			sortedKeys = append(sortedKeys, monthBucket)
			usageMonth = &UsageMonth{
				Month:     monthBucket,
				UsageDays: make([]UsageDay, 0),
				UsageKwh:  0,
			}
			yearMonthMap[monthBucket] = usageMonth
		}
		// Append the data to the month
		usageMonth.UsageDays = append(usageMonth.UsageDays, d)
		usageMonth.UsageKwh += d.UsageKwh
	}

	values := make([]UsageMonth, 0)
	for _, k := range sortedKeys {
		values = append(values, *yearMonthMap[k])
	}

	// Validate data points
	for _, v := range values {
		expectedYear := v.Month.Year()
		expectedMonth := v.Month.Month()

		for _, d := range v.UsageDays {
			if d.Day.Month() != expectedMonth {
				panic("Wrong month within UsageMonth")
			}
			if d.Day.Year() != expectedYear {
				panic("Wrong year within UsageMonth")
			}
		}
	}
	return values, nil
}

func truncateToMonth(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), 1, 0, 0, 0, 0, in.Location())
}
