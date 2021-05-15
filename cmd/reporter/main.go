package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)

	hours, err := analyzer.AggregateIntoHourWindows(csv)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d data points (%d days).\n", len(csv), len(hours)/24)
	fmt.Printf("First date: %+v\n", hours[0].StartTime)
	fmt.Printf("Last date: %+v\n", hours[len(hours)-1].StartTime)

	totalUsage := 0.0
	for _, hr := range csv {
		totalUsage += hr.UsageKwh
	}
	fmt.Printf("Total usage: %.2f kWh.\n\n", totalUsage)

	days, err := analyzer.SplitByDay(hours)
	if err != nil {
		panic(err)
	}

	domesticCost := costcalculator.CalculateDomesticCost(days)
	touBill := costcalculator.CalculateTouDACostForDays(days)
	touDANemCost := touBill.NetMeteredCost()
	nbc := touBill.NonBypassableCharges()
	baselineDiscount := touBill.BaselineCredit()
	touDATotal := touDANemCost + nbc + baselineDiscount

	_, _ = fmt.Fprintf(w, "Start\t%s\t\n", days[0].Day.Format("2006-01-02"))
	_, _ = fmt.Fprintf(w, "End (excl.)\t%s\t\n", days[len(days)-1].EndTime().Format("2006-01-02"))
	_, _ = fmt.Fprintf(w, "Time\t%d\tDays\t\n", len(days))
	_, _ = fmt.Fprintf(w, "Energy exported\t%.2f\tKWh\t\n", -1*touBill.EnergyExported())
	_, _ = fmt.Fprintf(w, "Energy imported\t%.2f\tKWh\t\n", touBill.EnergyImported())
	_, _ = fmt.Fprintf(w, "Net usage\t%.2f\tKWh\t\n", touBill.NetEnergyUsage())
	_, _ = fmt.Fprintf(w, "Average daily usage\t%.2f\tKWh\t\n", touBill.AverageDailyUsage())
	_, _ = fmt.Fprintf(w, "-------\t-------\t\n")
	_, _ = fmt.Fprintf(w, "Domestic est.\t%.2f\t$\t\n", domesticCost)
	_, _ = fmt.Fprintf(w, "-------\t-------\t\n")

	for period, usage := range touBill.UsageByPeriod() {
		_, _ = fmt.Fprintf(w, "%s\t%.2f\tKWh\t\n", period.Name(), usage)
	}
	_, _ = fmt.Fprintf(w, "-------\t-------\t\n")
	_, _ = fmt.Fprintf(w, "Non-bypassable charges\t%.2f\t$\t\n", nbc)
	_, _ = fmt.Fprintf(w, "Max baseline allocation\t%.2f\tKWh\t\n", touBill.MaxBaselineAllowance())
	_, _ = fmt.Fprintf(w, "Baseline discount\t%.2f\t$\t\n", baselineDiscount)
	_, _ = fmt.Fprintf(w, "TOU-D-A (NEM).\t%.2f\t$\t\n", touDANemCost)
	_, _ = fmt.Fprintf(w, "TOU-D-A (Total).\t%.2f\t$\t\n", touDATotal)
	_, _ = fmt.Fprintln(w)
	_ = w.Flush()
}
