package csvparser

import (
	"io/ioutil"
	"testing"

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

func readOrDie(file string) string {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(fileBytes)
}
