package analyzer

import (
	"fmt"
	"sort"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/csvparser"
)

type UsageHour struct {
	StartTime  time.Time
	DataPoints []csvparser.CsvRow
}

func (h *UsageHour) UsageKwh() float64 {
	total := 0.0
	for _, p := range h.DataPoints {
		total += p.UsageKwh
	}
	return total
}

func (h *UsageHour) EndTime() time.Time {
	return h.DataPoints[len(h.DataPoints)-1].EndTime
}

func AggregateIntoHourWindows(parsedFile csvparser.CsvFile) ([]UsageHour, error) {
	valuesByHour := make(map[time.Time][]csvparser.CsvRow)

	for _, v := range parsedFile {
		hr := truncateToHour(v.StartTime)
		hrEnd := truncateToHour(v.EndTime.Add(-1 * time.Second))
		if hr != hrEnd {
			return nil, fmt.Errorf("start and end of data point should be within the same hour for %+v", v)
		}
		valuesByHour[hr] = append(valuesByHour[hr], v)
	}

	out := make([]UsageHour, 0)
	for k, v := range valuesByHour {
		out = append(out, UsageHour{
			StartTime:  k,
			DataPoints: v,
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
