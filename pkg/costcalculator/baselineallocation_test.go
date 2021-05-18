package costcalculator

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMedicalBaselineEnabled_BaselineIsHigher(t *testing.T) {
	err := flag.CommandLine.Parse([]string{"--use_medical_baseline=false"})
	assert.NoError(t, err)
	noMedicalBaseline := GetDailyAllocation(now)

	err = flag.CommandLine.Parse([]string{"--use_medical_baseline=true"})
	assert.NoError(t, err)
	withMedicalBaseline := GetDailyAllocation(now)

	assert.Greater(t, withMedicalBaseline, noMedicalBaseline)
}
