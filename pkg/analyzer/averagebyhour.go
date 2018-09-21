package analyzer

import "fmt"

func CalculateAverageUsageByHour(in []UsageMonth) {
	usageByHour := make(map[int]float64)
	countByHour := make(map[int]int)
	for _, m := range in {
		for _, d := range m.UsageDays {
			for _, h := range d.DataPoints {
				usageByHour[h.StartTime.Hour()] += h.UsageKwh
				countByHour[h.StartTime.Hour()]++
			}
		}
	}

	for i := 0; i < 24; i++ {
		ave := usageByHour[i] / float64(countByHour[i])
		fmt.Printf("%.2f\n", ave)
	}
}
