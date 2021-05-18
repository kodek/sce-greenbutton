package sense

import (
	"encoding/csv"
	"io"
	time "time"

	"github.com/gocarina/gocsv"
)

// CsvRow represents a single row from a Sense data export file.
// From header:
// DateTime,Device ID,Name,Device Type,Device Make,Device Model,Device Location,Avg Wattage,kWh
type CsvRow struct {
	DateTime          DateTime `csv:"DateTime"`
	DeviceId          string   `csv:"Device ID"`
	Name              string   `csv:"Name"`
	DeviceType        string   `csv:"Device Type"`
	DeviceMake        string   `csv:"Device Make"`
	DeviceModel       string   `csv:"Device Model"`
	DeviceLocation    string   `csv:"Device Location"`
	AveragePowerWatts float64  `csv:"Avg Wattage"`
	EnergyKwh         float64  `csv:"kWh"`
}

type DateTime struct {
	Value time.Time
}

func (date *DateTime) UnmarshalCSV(csv string) error {
	um, err := time.ParseInLocation("2006-01-02 15:04:05", csv, time.UTC)
	if err != nil {
		return err
	}
	date.Value = um
	return nil
}

func ParseCSV(fileIn string) ([]CsvRow, error) {
	gocsv.SetCSVReader(func(reader io.Reader) gocsv.CSVReader {
		r := csv.NewReader(reader)
		r.Comment = '#'
		r.FieldsPerRecord = 9
		r.TrimLeadingSpace = true
		return r
	})

	out := make([]CsvRow, 0)
	if err := gocsv.UnmarshalString(fileIn, &out); err != nil {
		return nil, err
	}
	return out, nil
}
