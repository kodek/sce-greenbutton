package analyzer

import (
	"sort"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/csvparser"
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
			StartTime: v.StartTime,
			UsageKwh:  v.UsageKwh,
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartTime.Before(out[j].StartTime)
	})
	return out, nil
}
