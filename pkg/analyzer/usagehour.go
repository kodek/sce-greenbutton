package analyzer

import (
	"fmt"
	"sort"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/csvparser"
)

type UsageHour interface {
	UsageKwh() float64
	StartTime() time.Time
	EndTime() time.Time
}

type GreenButtonHour struct {
	startTime  time.Time
	DataPoints []csvparser.CsvRow
}

func (h *GreenButtonHour) UsageKwh() float64 {
	total := 0.0
	for _, p := range h.DataPoints {
		total += p.UsageKwh
	}
	return total
}

func (h *GreenButtonHour) StartTime() time.Time {
	return h.DataPoints[0].StartTime
}

func (h *GreenButtonHour) EndTime() time.Time {
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

	hours := make([]GreenButtonHour, 0)
	// Map iteration is random in Go, but we will sort the data points in chronological order later.
	for k, v := range valuesByHour {
		hours = append(hours, GreenButtonHour{
			startTime:  k,
			DataPoints: v,
		})
	}
	sort.Slice(hours, func(i, j int) bool {
		return hours[i].StartTime().Before(hours[j].StartTime())
	})

	out := make([]UsageHour, len(hours))
	for i, _ := range hours {
		out[i] = &hours[i]
	}
	return out, nil
}

func truncateToHour(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), in.Hour(), 0, 0, 0, in.Location())
}
