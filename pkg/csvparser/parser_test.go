package csvparser

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOneDayHasFifteenMinutePoints(t *testing.T) {
	got, err := Parse(readOrDie("one_day_constant_power.csv"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, got, 24*4)
}

func TestOneDayAt1KwIs24Kwh(t *testing.T) {
	got, err := Parse(readOrDie("one_day_constant_power.csv"))
	if err != nil {
		t.Fatal(err)
	}

	sum := 0.0
	for _, p := range got {
		sum += p.UsageKwh
	}

	assert.Equal(t, 24.0, sum)
}

func TestTwoDaysIgnoresIntermediateHeaders(t *testing.T) {
	got, err := Parse(readOrDie("two_days_constant_power.csv"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, got, 24*4*2)

	sum := 0.0
	for _, p := range got {
		sum += p.UsageKwh
	}
	assert.Equal(t, 48.0, sum)
}

func TestSingleLineParsesIntoTimesAndUsage(t *testing.T) {
	file := addHeaderTo([]string{
		`"2020-01-01 00:00:00 to 2020-01-01 00:15:00","1234",""`,
	})
	got, err := Parse(file)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, got, 1)
	assert.Equal(t, time.Date(2020, 01, 01, 00, 00, 0, 0, time.UTC), got[0].StartTime)
	assert.Equal(t, time.Date(2020, 01, 01, 00, 15, 0, 0, time.UTC), got[0].EndTime)
	assert.Equal(t, 1234.0, got[0].UsageKwh)
}

func TestDurationWorksOnSpringDaylightSavings(t *testing.T) {
	file := addHeaderTo([]string{
		`"2021-03-14 01:45:00 to 2021-03-14 02:00:00","1",""`,
	})
	got, err := Parse(file)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, 15*time.Minute, got[0].Duration())
}

func TestDurationWorksOnFallDaylightSavings(t *testing.T) {
	file := addHeaderTo([]string{
		`"2020-11-01 01:45:00 to 2020-11-01 02:00:00","1.0",""`,
	})
	got, err := Parse(file)
	assert.NoError(t, err)

	assert.Len(t, got, 1)
	assert.Equal(t, 15*time.Minute, got[0].Duration())
}

func readOrDie(file string) string {
	fileBytes, err := ioutil.ReadFile("testdata/" + file)
	if err != nil {
		panic(err)
	}
	return string(fileBytes)
}

func addHeaderTo(lines []string) string {
	header := readOrDie("header_only.csv")
	return strings.Join(append([]string{header}, lines...), "\n")
}
