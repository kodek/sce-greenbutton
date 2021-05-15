package analyzer

import (
	"fmt"
	"io"
)

func CalculateAverageUsageByHour(in []UsageMonth, w io.Writer) {
	usageByHour := make(map[int]float64)
	countByHour := make(map[int]int)
	for _, m := range in {
		for _, d := range m.UsageDays {
			for _, h := range d.DataPoints {
				usageByHour[h.StartTime.Hour()] += h.UsageKwh()
				countByHour[h.StartTime.Hour()]++
			}
		}
	}

	_, _ = fmt.Fprintln(w, "Hour of day\tData Points\tAverage Usage\t")
	for i := 0; i < 24; i++ {
		ave := usageByHour[i] / float64(countByHour[i])
		_, _ = fmt.Fprintf(w, "%d\t%d\t%.2f\t\n", i, countByHour[i], ave)
	}
}
