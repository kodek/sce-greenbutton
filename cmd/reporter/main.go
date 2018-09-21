package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kodek/greenbutton/pkg/costcalculator"

	"github.com/kodek/greenbutton/pkg/analyzer"
	"github.com/kodek/greenbutton/pkg/csvparser"
)

func main() {
	filePath := "/Users/hcosi/Downloads/Marter - SCE_Usage_3-027-2852-20_08-31-17_to_08-31-18.csv"

	//filePath := "/Users/hcosi/Downloads/sce marter partial bill (2018 01 to 08).csv"
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	csv, err := csvparser.Parse(string(file))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Read %d hours (%d days).\n", len(csv), len(csv)/24)

	addCarSimulation(csv)

	annualUsage := 0.0
	for _, hr := range csv {
		annualUsage += hr.UsageKwh
	}
	fmt.Printf("Total annual usage: %.2f kWh.\n\n", annualUsage)

	// Split into months
	hours, err := analyzer.ParseIntoHours(csv)
	if err != nil {
		panic(err)
	}

	days, err := analyzer.SplitByDay(hours)
	if err != nil {
		panic(err)
	}

	months, err := analyzer.SplitByMonth(days)
	if err != nil {
		panic(err)
	}

	totalDomesticCost := 0.0
	totalTouDACost := 0.0
	fmt.Println("Month, Days, Usage, AveDailyUsage, DOMESTIC, TOU-D-A")
	for _, month := range months {
		domesticCost := costcalculator.CalculateDomesticCost(month)
		totalDomesticCost += domesticCost

		touDACost := costcalculator.CalculateTouDACostForMonth(month)
		totalTouDACost += touDACost

		fmt.Printf("%d-%d, %d, %.2f,%.2f, $%.2f, $%.2f\n", month.Month.Year(), month.Month.Month(), len(month.UsageDays), month.UsageKwh, month.AverageDailyUsageKwh(), domesticCost, touDACost)
	}
	fmt.Printf("Total DOMESTIC: $%.2f.\n", totalDomesticCost)
	fmt.Printf("Total TOU-D-A: $%.2f.\n", totalTouDACost)

	analyzer.CalculateAverageUsageByHour(months)
}
func addCarSimulation(f csvparser.CsvFile) {
	for i, hr := range f {
		if hr.StartOfHour.Before(time.Date(2018, 07, 10, 0, 0, 0, 0, time.UTC)) {
			h := hr.StartOfHour.Hour()
			if !(h >= 1 && h < 22) {
				f[i].UsageKwh += 8
			}
		}
	}
}
