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

	assert.Equal(t, len(got), 24*4)
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

	assert.Equal(t, sum, 24.0)
}

func TestSingleLineParsedCorrectly(t *testing.T) {
	file := addHeaderTo([]string{
		`"2020-01-01 00:00:00 to 2020-01-01 00:15:00","1234",""`,
	})
	got, err := Parse(file)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got), 1)
	assert.Equal(t, got[0].StartTime, time.Date(2020, 01, 01, 00, 00, 0, 0, time.Local))
	assert.Equal(t, got[0].EndTime, time.Date(2020, 01, 01, 00, 15, 0, 0, time.Local))
	assert.Equal(t, got[0].UsageKwh, 1234.0)
}

func readOrDie(file string) string {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(fileBytes)
}

func addHeaderTo(lines []string) string {
	header := readOrDie("header_only.csv")
	return strings.Join(append([]string{header}, lines...), "\n")
}
