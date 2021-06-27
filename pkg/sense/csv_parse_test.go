package sense

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseCSV_RealSamplePasses(t *testing.T) {
	in := `# Please note that there may be discrepancies between your electric bill and what Sense reports due to data interruptions (such as Wi-Fi signal loss or loss of power) or other issues
DateTime,Device ID,Name,Device Type,Device Make,Device Model,Device Location,Avg Wattage,kWh
2020-01-01 00:00:00,3a9fb50e,Amy’s nightstand,Light,Signify Netherlands B.V.,LCT014,Amy’s bedroom,3.497,0.003
2020-01-01 00:00:00,63ac628e,Dryer,Dryer,,,,66.569,0.067`

	got, err := ParseCSV(in)
	assert.NoError(t, err)

	assert.Len(t, got, 2)
}

func TestParseCSV_ParsesDate(t *testing.T) {
	in := `DateTime,Device ID,Name,Device Type,Device Make,Device Model,Device Location,Avg Wattage,kWh
2020-01-01 00:00:00,63ac628e,Dryer,Dryer,,,,66.569,0.067`

	got, err := ParseCSV(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, time.Date(2020, 01, 01, 00, 00, 00, 00, time.UTC), got[0].DateTime.Value)
}

func TestParseCSV_ParsesPowerAndEnergy(t *testing.T) {
	in := `DateTime,Device ID,Name,Device Type,Device Make,Device Model,Device Location,Avg Wattage,kWh
2020-01-01 00:00:00,63ac628e,Dryer,Dryer,,,,66.569,0.067`

	got, err := ParseCSV(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, 66.569, got[0].AveragePowerWatts)
	assert.Equal(t, 0.067, got[0].EnergyKwh)
}
func TestParseCSV_MissingColumns_Fails(t *testing.T) {
	in := `
DateTime,Device ID,Name,Device Type,Device Make,Device Model,Device Location,Avg Wattage,kWh
2020-01-01 00:00:00,63ac628e`

	_, err := ParseCSV(in)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wrong number of fields")
}

func TestParseCSV_SampleSucceeds(t *testing.T) {
	sample, err := ioutil.ReadFile("testdata/sample.csv")
	assert.NoError(t, err)

	got, err := ParseCSV(string(sample))
	assert.NoError(t, err)

	assert.Len(t, got, 48)
}
