package main

import (
	"flag"
	"fmt"
	"io/ioutil"

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

	hours, err := analyzer.AggregateIntoHourWindows(csv)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d data points (%d days).\n", len(csv), len(hours)/24)
	fmt.Printf("First date: %+v\n", hours[0])
	fmt.Printf("Last date: %+v\n", hours[len(hours)-1])

	totalUsage := 0.0
	for _, hr := range csv {
		totalUsage += hr.UsageKwh
	}
	fmt.Printf("Total usage: %.2f kWh.\n\n", totalUsage)

	days, err := analyzer.SplitByDay(hours)
	if err != nil {
		panic(err)
	}

	_ = costcalculator.CalculateDomesticCost(days)
	_ = costcalculator.CalculateTouDACostForMonth(days)
}
