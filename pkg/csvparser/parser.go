package csvparser

import (
	"strings"
	"time"
)

// Parses a CSV file in the following format. This is two days of data:

/*
Energy Usage Information
"For location: 1234 FAKE ST"

Meter Reading Information
"Type of readings: Electricity"

Summary of Electric Power Usage Information*
"Your download will contain interval usage data that is currently available for your selected Service Account. Based on how our systems process and categorize usage data, your download may contain usage data of the following types: actual, estimated, validated or missing. "

Detailed Usage
"Start date: 2017-08-30 23:00:00  for 366 days"

"Data for period starting: 2017-08-31 00:00:00  for 24 hours"
Energy consumption time period,Usage(Real energy in kilowatt-hours),Reading quality
"2017-08-31 00:00:00 to 2017-08-31 01:00:00","0.610",""
"2017-08-31 01:00:00 to 2017-08-31 02:00:00","0.610",""
"2017-08-31 02:00:00 to 2017-08-31 03:00:00","0.590",""
"2017-08-31 03:00:00 to 2017-08-31 04:00:00","0.600",""
"2017-08-31 04:00:00 to 2017-08-31 05:00:00","0.670",""
"2017-08-31 05:00:00 to 2017-08-31 06:00:00","0.700",""
"2017-08-31 06:00:00 to 2017-08-31 07:00:00","2.720",""
"2017-08-31 07:00:00 to 2017-08-31 08:00:00","1.290",""
"2017-08-31 08:00:00 to 2017-08-31 09:00:00","-0.140",""
"2017-08-31 09:00:00 to 2017-08-31 10:00:00","0.440",""
"2017-08-31 10:00:00 to 2017-08-31 11:00:00","-0.620",""
"2017-08-31 11:00:00 to 2017-08-31 12:00:00","-0.430",""
"2017-08-31 12:00:00 to 2017-08-31 13:00:00","0.800",""
"2017-08-31 13:00:00 to 2017-08-31 14:00:00","0.660",""
"2017-08-31 14:00:00 to 2017-08-31 15:00:00","1.340",""
"2017-08-31 15:00:00 to 2017-08-31 16:00:00","1.240",""
"2017-08-31 16:00:00 to 2017-08-31 17:00:00","2.660",""
"2017-08-31 17:00:00 to 2017-08-31 18:00:00","4.210",""
"2017-08-31 18:00:00 to 2017-08-31 19:00:00","4.540",""
"2017-08-31 19:00:00 to 2017-08-31 20:00:00","2.170",""
"2017-08-31 20:00:00 to 2017-08-31 21:00:00","2.950",""
"2017-08-31 21:00:00 to 2017-08-31 22:00:00","2.550",""
"2017-08-31 22:00:00 to 2017-08-31 23:00:00","2.520",""
"2017-08-31 23:00:00 to 2017-09-01 00:00:00","1.390",""

"Data for period starting: 2017-09-01 00:00:00  for 24 hours"
Energy consumption time period,Usage(Real energy in kilowatt-hours),Reading quality
"2017-09-01 00:00:00 to 2017-09-01 01:00:00","1.440",""
"2017-09-01 01:00:00 to 2017-09-01 02:00:00","0.600",""
"2017-09-01 02:00:00 to 2017-09-01 03:00:00","1.380",""
"2017-09-01 03:00:00 to 2017-09-01 04:00:00","1.270",""
"2017-09-01 04:00:00 to 2017-09-01 05:00:00","0.670",""
"2017-09-01 05:00:00 to 2017-09-01 06:00:00","1.220",""
"2017-09-01 06:00:00 to 2017-09-01 07:00:00","1.060",""
"2017-09-01 07:00:00 to 2017-09-01 08:00:00","1.030",""
"2017-09-01 08:00:00 to 2017-09-01 09:00:00","0.480",""
"2017-09-01 09:00:00 to 2017-09-01 10:00:00","0.750",""
"2017-09-01 10:00:00 to 2017-09-01 11:00:00","0.330",""
"2017-09-01 11:00:00 to 2017-09-01 12:00:00","0.510",""
"2017-09-01 12:00:00 to 2017-09-01 13:00:00","1.550",""
"2017-09-01 13:00:00 to 2017-09-01 14:00:00","1.280",""
"2017-09-01 14:00:00 to 2017-09-01 15:00:00","1.520",""
"2017-09-01 15:00:00 to 2017-09-01 16:00:00","3.170",""
"2017-09-01 16:00:00 to 2017-09-01 17:00:00","3.260",""
"2017-09-01 17:00:00 to 2017-09-01 18:00:00","2.980",""
"2017-09-01 18:00:00 to 2017-09-01 19:00:00","2.880",""
"2017-09-01 19:00:00 to 2017-09-01 20:00:00","0.450",""
"2017-09-01 20:00:00 to 2017-09-01 21:00:00","1.200",""
"2017-09-01 21:00:00 to 2017-09-01 22:00:00","0.610",""
"2017-09-01 22:00:00 to 2017-09-01 23:00:00","0.740",""
"2017-09-01 23:00:00 to 2017-09-02 00:00:00","0.510",""
*/

type CsvRow struct {
	StartTime      time.Time
	EndTime        time.Time
	UsageKwh       float64
	ReadingQuality string
}

func (r *CsvRow) Duration() time.Duration {
	return r.EndTime.Sub(r.StartTime)
}

type CsvFile []CsvRow

func Parse(fileIn string) (CsvFile, error) {
	lines := strings.Split(fileIn, "\n")
	out := make(CsvFile, 0)
	for lNum, l := range lines {
		parsed, err := parseHourConsumption(l, lNum)
		if err != nil {
			if err.(*LineParsingError).CanBeIgnored {
				continue
			}
			return nil, err
		}
		out = append(out, *parsed)
	}
	return out, nil
}
