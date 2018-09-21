package csvparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type LineParsingError struct {
	CanBeIgnored bool
	Cause        error
	LineNumber   int
	LineText     string
}

func (e *LineParsingError) Error() string {
	return fmt.Sprintf("Error when reading line %d:\n\n%s\n\n\t%s", e.LineNumber, e.LineText, e.Cause.Error())
}

func parseHourConsumption(line string, lineNumber int) (*CsvRow, error) {
	if line == HEADER {
		return nil, &LineParsingError{
			CanBeIgnored: true,
			Cause:        errors.New("Line is a CSV header. This line should be ignored"),
			LineText:     line,
			LineNumber:   lineNumber,
		}
	}

	csvSplit := strings.Split(line, ",")
	if len(csvSplit) != 3 {
		err := errors.New("Expected 3 components in line (time period, usage, reading quality).")
		return nil, &LineParsingError{
			// TODO: This could be a malformed ACTUAL line, so we should be smart about which csv lines to ignore.
			CanBeIgnored: true,
			Cause:        err,
			LineText:     line,
			LineNumber:   lineNumber,
		}
	}

	timePeriod := csvSplit[0]
	usage := csvSplit[1]
	readingQuality := csvSplit[2]

	tStart, tEnd, err := parseTimePeriod(timePeriod)
	if err != nil {
		return nil, &LineParsingError{
			Cause:      err,
			LineText:   line,
			LineNumber: lineNumber,
		}
	}
	err = assertTimePeriodIsOneHour(tStart, tEnd)
	if err != nil {
		return nil, &LineParsingError{
			Cause:      err,
			LineText:   line,
			LineNumber: lineNumber,
		}
	}
	usageNum, err := parseUsage(usage)
	if err != nil {
		return nil, &LineParsingError{
			Cause:      err,
			LineText:   line,
			LineNumber: lineNumber,
		}
	}

	return &CsvRow{
		StartOfHour:    tStart,
		UsageKwh:       usageNum,
		ReadingQuality: readingQuality}, nil
}

const QUOTE = "\""
const HEADER = "Energy consumption time period,Usage(Real energy in kilowatt-hours),Reading quality\r"

// Parses a string in the format of (quotes included):
// "2017-09-01 23:00:00 to 2017-09-02 00:00:00",
func parseTimePeriod(timePeriod string) (time.Time, time.Time, error) {
	noQuotes := removeQuotes(timePeriod)
	split := strings.Split(noQuotes, "\xa0to ")
	if len(split) != 2 {
		return time.Time{}, time.Time{}, errors.New(fmt.Sprintf("Expected time period with 2 dates, but got %d", len(split)))
	}
	beforeStr := split[0]
	afterStr := split[1]

	before, err := parseTime(beforeStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	after, err := parseTime(afterStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return before, after, nil
}

// parseTime parses a date and time string and returns a local time.Time object
// t: a string in the format of "2006-01-02 15:04:05
func parseTime(t string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", t, time.UTC)
}

func assertTimePeriodIsOneHour(before time.Time, after time.Time) error {
	diff := after.Sub(before)
	if diff != 1*time.Hour {
		return errors.New(fmt.Sprintf("Expected time period of 1 hour between %s and %s, but got %s.", before, after, diff))
	}
	return nil
}

func parseUsage(usageStr string) (float64, error) {
	noQuotes := removeQuotes(usageStr)
	return strconv.ParseFloat(noQuotes, 64)
}

func removeQuotes(in string) string {
	return strings.TrimPrefix(strings.TrimSuffix(in, QUOTE), QUOTE)
}
