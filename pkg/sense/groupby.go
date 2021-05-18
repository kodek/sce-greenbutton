package sense

import (
	"log"
	"sort"
	"time"
)

// GroupByDeviceId returns a map keyed by the device ID. The value contains all the data points for that device
// in existing order.
func GroupByDeviceId(in []CsvRow) (map[string][]CsvRow, error) {
	out := make(map[string][]CsvRow)

	for _, point := range in {
		key := point.DeviceId
		out[key] = append(out[key], point)
	}
	return out, nil
}

type Snapshot struct {
	DateTime       time.Time
	ProductionKwh  *float64
	ConsumptionKwh *float64
}

// NetUsageKwh returns the energy consumed from the grid.
func (s *Snapshot) NetUsageKwh() float64 {
	return *s.ConsumptionKwh + *s.ProductionKwh
}

// GroupByTime returns a slice of Snapshot objects summarizing a single point in time.
// The objects are returned in chronological order.
func GroupByTime(in []CsvRow) ([]Snapshot, error) {
	rowsByTime := make(map[time.Time][]CsvRow)

	for _, point := range in {
		key := point.DateTime.Value
		rowsByTime[key] = append(rowsByTime[key], point)
	}

	out := make([]Snapshot, 0)
	for t, rows := range rowsByTime {
		mainsEnergy, err := findEnergy("mains", rows)
		if err != nil {
			return nil, err
		}
		solarEnergy, err := findEnergy("solar", rows)
		if err != nil {
			return nil, err
		}

		out = append(out, Snapshot{
			DateTime:       t,
			ProductionKwh:  solarEnergy,
			ConsumptionKwh: mainsEnergy,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].DateTime.Before(out[j].DateTime)
	})
	return out, nil
}

func findEnergy(deviceId string, rows []CsvRow) (*float64, error) {
	var found *float64 = nil
	for _, row := range rows {
		if row.DeviceId == deviceId {
			if found != nil {
				log.Printf("found device ID '%s' twice within timestamp %+v\n", deviceId, row.DateTime)
			}
			energy := row.EnergyKwh
			found = &energy
		}
	}
	return found, nil
}
