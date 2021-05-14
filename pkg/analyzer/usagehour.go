package analyzer

import (
	"fmt"
	"sort"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/csvparser"
)

type UsageHour struct {
	StartTime time.Time
	UsageKwh  float64
}

type SortedUsageHours []UsageHour

func AggregateIntoHourWindows(parsedFile csvparser.CsvFile) (SortedUsageHours, error) {
	usageByHour := make(map[time.Time]float64)

	for _, v := range parsedFile {
		hr := truncateToHour(v.StartTime)
		hrEnd := truncateToHour(v.EndTime.Add(-1 * time.Second))
		if hr != hrEnd {
			return nil, fmt.Errorf("start and end of data point should be within the same hour for %+v", v)
		}
		usageByHour[hr] += v.UsageKwh
	}

	out := make([]UsageHour, 0)
	for k, v := range usageByHour {
		out = append(out, UsageHour{
			StartTime: k,
			UsageKwh:  v,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartTime.Before(out[j].StartTime)
	})
	return out, nil
}

func truncateToHour(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), in.Hour(), 0, 0, 0, in.Location())
}
