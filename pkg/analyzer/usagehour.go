package analyzer

import (
	"github.com/kodek/greenbutton/pkg/csvparser"
	"sort"
	"time"
)

type UsageHour struct {
	StartTime time.Time
	UsageKwh  float64
}

type SortedUsageHours []UsageHour

func ParseIntoHours(parsedFile csvparser.CsvFile) (SortedUsageHours, error) {
	out := make([]UsageHour, len(parsedFile))
	for i, v := range parsedFile {
		out[i] = UsageHour{
			StartTime: v.StartOfHour,
			UsageKwh:  v.UsageKwh,
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartTime.Before(out[j].StartTime)
	})
	return out, nil
}
