package costcalculator

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/analyzer"
	"github.com/kodek/sce-greenbutton/pkg/csvparser"
	"github.com/stretchr/testify/assert"
)

var now = time.Date(2020, 01, 01, 00, 00, 00, 00, time.UTC)

func TestTouBillSummary_NetEnergyUsage_MatchesImportAndExport(t *testing.T) {
	days := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1.0),
		csvparser.NewRowWith15MinuteDuration(now.Add(24*time.Hour), -1.0),
	})

	bill := CalculateTouDACostForDays(days)

	assert.Equal(t, 0.0, bill.NetEnergyUsage())
	assert.Equal(t, bill.EnergyImported()+bill.EnergyExported(), bill.NetEnergyUsage())
}

func TestTouBillSummary_NetMeteredCost_CreditsCancelOut(t *testing.T) {
	days := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1.0),
		csvparser.NewRowWith15MinuteDuration(now.Add(24*time.Hour), -1.0),
	})

	got := CalculateTouDACostForDays(days)

	assert.Equal(t, 0.0, got.NetMeteredCostNoBaseline())
}

func TestTouBillSummary_AverageDailyUsage_AveragesAcrossDays(t *testing.T) {
	days := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1.0),
		csvparser.NewRowWith15MinuteDuration(now.Add(24*time.Hour), 2.0),
	})

	got := CalculateTouDACostForDays(days)

	assert.Equal(t, 1.5, got.AverageDailyUsage())
}

func TestTouBillSummary_EnergyExported_SameHour_CancelsOut(t *testing.T) {
	days := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1.0),
		csvparser.NewRowWith15MinuteDuration(now.Add(15*time.Minute), -1.0),
	})

	got := CalculateTouDACostForDays(days)

	assert.Equal(t, 0.0, got.EnergyExported())
}

func TestTouBillSummary_EnergyExported_DifferentHours_ReturnsNegativeSum(t *testing.T) {
	days := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1.0),
		csvparser.NewRowWith15MinuteDuration(now.Add(1*time.Hour), -1.0),
	})

	got := CalculateTouDACostForDays(days)

	assert.Equal(t, -1.0, got.EnergyExported())
}

func TestTouBillSummary_EnergyImported_SameHour_CancelsOut(t *testing.T) {
	days := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1.0),
		csvparser.NewRowWith15MinuteDuration(now.Add(15*time.Minute), -1.0),
	})

	got := CalculateTouDACostForDays(days)

	assert.Equal(t, 0.0, got.EnergyImported())
}

func TestTouBillSummary_EnergyImported_DifferentHours_ReturnsPositiveSum(t *testing.T) {
	days := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1.0),
		csvparser.NewRowWith15MinuteDuration(now.Add(1*time.Hour), -1.0),
	})

	got := CalculateTouDACostForDays(days)

	assert.Equal(t, 1.0, got.EnergyImported())
}

func TestTouBillSummary_UsageByPeriod_SumMatchesUsage(t *testing.T) {
	summerWeekday := time.Date(2020, 8, 3, 0, 0, 0, 0, time.UTC)
	days := toDaysOrDie(t, oneDataPointPerHourWithConstantUsage(summerWeekday, 1))
	assert.Len(t, days, 1)
	assert.Len(t, days[0].DataPoints, 24)

	bill := CalculateTouDACostForDays(days)
	got := bill.UsageByPeriod()

	// A summer weekday has 3 cost periods.
	assert.Len(t, got, 3)

	totalUsage := 0.0
	for _, usage := range got {
		totalUsage += usage
	}
	assert.Equal(t, bill.NetEnergyUsage(), totalUsage)
}

func TestTouBillSummary_BaselineCredit_ReducesCostWhenNetConsumption(t *testing.T) {
	days := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 2.0),
		csvparser.NewRowWith15MinuteDuration(now.Add(1*time.Hour), -1.0),
	})

	bill := CalculateTouDACostForDays(days)
	got := bill.BaselineCredit()

	assert.Positive(t, bill.NetEnergyUsage())
	assert.Negative(t, got)
}

func TestTouBillSummary_BaselineCredit_IncreasesCostWhenNetExport(t *testing.T) {
	days := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1.0),
		csvparser.NewRowWith15MinuteDuration(now.Add(1*time.Hour), -2.0),
	})

	bill := CalculateTouDACostForDays(days)
	got := bill.BaselineCredit()

	assert.Negative(t, bill.NetEnergyUsage())
	assert.Positive(t, got)
}

func TestTouBillSummary_BaselineCredit_DiscountHasCeilingBasedOnUsage(t *testing.T) {
	days1 := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1000.0),
	})
	largeBill := CalculateTouDACostForDays(days1)

	days2 := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 2000.0),
	})
	largerBill := CalculateTouDACostForDays(days2)

	assert.Greater(t, largerBill.NetMeteredCostNoBaseline(), largeBill.NetMeteredCostNoBaseline())
	assert.Equal(t, largeBill.BaselineCredit(), largerBill.BaselineCredit())
}

func TestTouBillSummary_NegativeTrueUp_MuchCheaperDueToSurplusRate(t *testing.T) {
	surplusDays := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, -1000.0),
	})
	surplusBill := CalculateTouDACostForDays(surplusDays)

	deficitDays := toDaysOrDie(t, []csvparser.CsvRow{
		csvparser.NewRowWith15MinuteDuration(now, 1000.0),
	})
	deficitBill := CalculateTouDACostForDays(deficitDays)

	assert.Negative(t, surplusBill.TrueUp())
	assert.Positive(t, deficitBill.TrueUp())
	assert.Equal(t, math.Abs(deficitBill.NetMeteredCostNoBaseline()), math.Abs(surplusBill.NetMeteredCostNoBaseline()))
	assert.Equal(t, math.Abs(deficitBill.BaselineCredit()), math.Abs(surplusBill.BaselineCredit()))
	assert.Greater(t, math.Abs(deficitBill.TrueUp()), math.Abs(surplusBill.TrueUp()))
}

func oneDataPointPerHourWithConstantUsage(day time.Time, usage float64) []csvparser.CsvRow {
	rows := make([]csvparser.CsvRow, 0)
	for i := 0; i < 24; i++ {
		offset := time.Duration(i) * time.Hour
		rows = append(rows, csvparser.NewRowWith15MinuteDuration(day.Add(offset), usage))
	}
	return rows
}

func toDaysOrDie(t *testing.T, rows []csvparser.CsvRow) []analyzer.UsageDay {
	hours, err := analyzer.AggregateIntoHourWindows(rows)
	assert.NoError(t, err)

	days, err := analyzer.SplitByDay(hours)
	assert.NoError(t, err)
	return days
}

func TestTouDA_SummerWeekday(t *testing.T) {
	expectedHours := []CostPeriod{
		SummerSuperOffPeak, // 00
		SummerSuperOffPeak, // 01
		SummerSuperOffPeak, // 02
		SummerSuperOffPeak, // 03
		SummerSuperOffPeak, // 04
		SummerSuperOffPeak, // 05
		SummerSuperOffPeak, // 06
		SummerSuperOffPeak, // 07
		SummerOffPeak,      // 08
		SummerOffPeak,      // 09
		SummerOffPeak,      // 10
		SummerOffPeak,      // 11
		SummerOffPeak,      // 12
		SummerOffPeak,      // 13
		SummerOnPeak,       // 14
		SummerOnPeak,       // 15
		SummerOnPeak,       // 16
		SummerOnPeak,       // 17
		SummerOnPeak,       // 18
		SummerOnPeak,       // 19
		SummerOffPeak,      // 20
		SummerOffPeak,      // 21
		SummerSuperOffPeak, // 22
		SummerSuperOffPeak, // 23
	}

	// Monday, Aug 4, 2020.
	summerWeekday := time.Date(2020, 8, 3, 0, 0, 0, 0, time.UTC)

	for i, period := range expectedHours {
		t.Run(fmt.Sprintf("Hour %d with expected period %f", i, period), func(t *testing.T) {
			date := summerWeekday.Add(time.Duration(i) * time.Hour)
			assert.Equal(t, period, calculateTouRateForHour(date, NewTouDAPlan()))
		})
	}
}
func TestTouDA_SummerWeekend(t *testing.T) {
	expectedHours := []CostPeriod{
		SummerSuperOffPeak, // 00
		SummerSuperOffPeak, // 01
		SummerSuperOffPeak, // 02
		SummerSuperOffPeak, // 03
		SummerSuperOffPeak, // 04
		SummerSuperOffPeak, // 05
		SummerSuperOffPeak, // 06
		SummerSuperOffPeak, // 07
		SummerOffPeak,      // 08
		SummerOffPeak,      // 09
		SummerOffPeak,      // 10
		SummerOffPeak,      // 11
		SummerOffPeak,      // 12
		SummerOffPeak,      // 13
		SummerOffPeak,      // 14
		SummerOffPeak,      // 15
		SummerOffPeak,      // 16
		SummerOffPeak,      // 17
		SummerOffPeak,      // 18
		SummerOffPeak,      // 19
		SummerOffPeak,      // 20
		SummerOffPeak,      // 21
		SummerSuperOffPeak, // 22
		SummerSuperOffPeak, // 23
	}

	// Saturday, Aug 1, 2020
	summerWeekend := time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC)
	for i, period := range expectedHours {
		t.Run(fmt.Sprintf("Hour %d with expected period %f", i, period), func(t *testing.T) {
			date := summerWeekend.Add(time.Duration(i) * time.Hour)
			assert.Equal(t, period, calculateTouRateForHour(date, NewTouDAPlan()))
		})
	}
}
func TestTouDA_WinterWeekday(t *testing.T) {
	expectedHours := []CostPeriod{
		WinterSuperOffPeak, // 00
		WinterSuperOffPeak, // 01
		WinterSuperOffPeak, // 02
		WinterSuperOffPeak, // 03
		WinterSuperOffPeak, // 04
		WinterSuperOffPeak, // 05
		WinterSuperOffPeak, // 06
		WinterSuperOffPeak, // 07
		WinterOffPeak,      // 08
		WinterOffPeak,      // 09
		WinterOffPeak,      // 10
		WinterOffPeak,      // 11
		WinterOffPeak,      // 12
		WinterOffPeak,      // 13
		WinterOnPeak,       // 14
		WinterOnPeak,       // 15
		WinterOnPeak,       // 16
		WinterOnPeak,       // 17
		WinterOnPeak,       // 18
		WinterOnPeak,       // 19
		WinterOffPeak,      // 20
		WinterOffPeak,      // 21
		WinterSuperOffPeak, // 22
		WinterSuperOffPeak, // 23
	}

	// Friday, Dec 4, 2020
	winterWeekday := time.Date(2020, 12, 4, 0, 0, 0, 0, time.UTC)

	for i, period := range expectedHours {
		t.Run(fmt.Sprintf("Hour %d with expected period %f", i, period), func(t *testing.T) {
			date := winterWeekday.Add(time.Duration(i) * time.Hour)
			assert.Equal(t, period, calculateTouRateForHour(date, NewTouDAPlan()))
		})
	}
}
func TestTouDA_WinterWeekend(t *testing.T) {
	expectedHours := []CostPeriod{
		WinterSuperOffPeak, // 00
		WinterSuperOffPeak, // 01
		WinterSuperOffPeak, // 02
		WinterSuperOffPeak, // 03
		WinterSuperOffPeak, // 04
		WinterSuperOffPeak, // 05
		WinterSuperOffPeak, // 06
		WinterSuperOffPeak, // 07
		WinterOffPeak,      // 08
		WinterOffPeak,      // 09
		WinterOffPeak,      // 10
		WinterOffPeak,      // 11
		WinterOffPeak,      // 12
		WinterOffPeak,      // 13
		WinterOffPeak,      // 14
		WinterOffPeak,      // 15
		WinterOffPeak,      // 16
		WinterOffPeak,      // 17
		WinterOffPeak,      // 18
		WinterOffPeak,      // 19
		WinterOffPeak,      // 20
		WinterOffPeak,      // 21
		WinterSuperOffPeak, // 22
		WinterSuperOffPeak, // 23
	}

	// Saturday, Dec 5, 2020
	winterWeekend := time.Date(2020, 12, 5, 0, 0, 0, 0, time.UTC)
	for i, period := range expectedHours {
		t.Run(fmt.Sprintf("Hour %d with expected period %f", i, period), func(t *testing.T) {
			date := winterWeekend.Add(time.Duration(i) * time.Hour)
			assert.Equal(t, period, calculateTouRateForHour(date, NewTouDAPlan()))
		})
	}
}
