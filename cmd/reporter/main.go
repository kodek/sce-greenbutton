package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
	"time"

	"github.com/kodek/sce-greenbutton/pkg/costcalculator"

	"github.com/kodek/sce-greenbutton/pkg/analyzer"
	"github.com/kodek/sce-greenbutton/pkg/csvparser"
)

var inputFilePath = flag.String("input_file_path", "", "Path to input CSV file from GreenButton.")

func main() {
	flag.Parse()
	if *inputFilePath == "" {
		panic("Must specify --input_file_path")
	}
	file, err := ioutil.ReadFile(*inputFilePath)
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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	_, _ = fmt.Fprintln(w, "Month\tDays\tUsage\tAveDailyUsage\tDOMESTIC\tTOU-D-A\t")
	for _, month := range months {
		domesticCost := costcalculator.CalculateDomesticCost(month)
		totalDomesticCost += domesticCost

		touDACost := costcalculator.CalculateTouDACostForMonth(month)
		totalTouDACost += touDACost

		_, _ = fmt.Fprintf(w, "%d-%d\t%d\t%.2f\t%.2f\t$%.2f\t$%.2f\t\n", month.Month.Year(), month.Month.Month(), len(month.UsageDays), month.UsageKwh, month.AverageDailyUsageKwh(), domesticCost, touDACost)
	}
	_ = w.Flush()
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
