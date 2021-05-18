package sense

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGroupByDeviceId_SameIdGroupedTogether(t *testing.T) {
	in := []CsvRow{
		{DeviceId: "abc", EnergyKwh: 10},
		{DeviceId: "abc", EnergyKwh: 5},
	}

	got, err := GroupByDeviceId(in)
	assert.NoError(t, err)

	assert.ElementsMatch(t, in, got["abc"])
}
func TestGroupByDeviceId_DifferentIdSplit(t *testing.T) {
	in := []CsvRow{
		{DeviceId: "abc", EnergyKwh: 10},
		{DeviceId: "def", EnergyKwh: 5},
	}

	got, err := GroupByDeviceId(in)
	assert.NoError(t, err)

	assert.ElementsMatch(t, []CsvRow{in[0]}, got["abc"])
	assert.ElementsMatch(t, []CsvRow{in[1]}, got["def"])
}

var (
	t0 = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 = t0.Add(1 * time.Hour)
)

func TestGroupByTime_SameTime_ReturnsOneElement(t *testing.T) {
	in := []CsvRow{
		{DeviceId: "mains", EnergyKwh: 10, DateTime: DateTime{Value: t0}},
		{DeviceId: "solar", EnergyKwh: -5, DateTime: DateTime{Value: t0}},
	}

	got, err := GroupByTime(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
}

func TestGroupByTime_DifferentTimes_ReturnsSeparateElements(t *testing.T) {
	in := []CsvRow{
		{DeviceId: "mains", EnergyKwh: 10, DateTime: DateTime{Value: t0}},
		{DeviceId: "solar", EnergyKwh: -5, DateTime: DateTime{Value: t0}},
		{DeviceId: "mains", EnergyKwh: 10, DateTime: DateTime{Value: t1}},
		{DeviceId: "solar", EnergyKwh: -5, DateTime: DateTime{Value: t1}},
	}

	got, err := GroupByTime(in)
	assert.NoError(t, err)

	assert.Len(t, got, 2)
}
func TestGroupByTime_AggregatesProductionAndConsumptionFromMainsAndSolar(t *testing.T) {
	in := []CsvRow{
		{DeviceId: "mains", EnergyKwh: 10, DateTime: DateTime{Value: t0}},
		{DeviceId: "solar", EnergyKwh: -5, DateTime: DateTime{Value: t0}},
	}

	got, err := GroupByTime(in)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, -5.0, *got[0].ProductionKwh)
	assert.Equal(t, 10.0, *got[0].ConsumptionKwh)
}
