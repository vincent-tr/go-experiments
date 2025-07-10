package backtesting

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"time"
)

const dataPath = "brokers/backtesting/data"

// https://www.histdata.com/download-free-forex-historical-data/?/ascii/tick-data-quotes/EURUSD

func loadFile(arrayPtr *[]tick, year int, month int, symbol string) error {

	zipFile := path.Join(dataPath, fmt.Sprintf("HISTDATA_COM_ASCII_%s_T%04d%02d.zip", symbol, year, month))

	// Unzip CSV
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("failed to open ZIP archive '%s': %v", zipFile, err)
	}
	defer r.Close()

	var csvFile io.ReadCloser

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".csv") {
			csvFile, err = f.Open()
			if err != nil {
				return fmt.Errorf("failed to open CSV file '%s' in ZIP archive '%s': %v", f.Name, zipFile, err)
			}

			break
		}
	}

	if csvFile == nil {
		return fmt.Errorf("no CSV file found in ZIP archive '%s'", zipFile)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = ','

	// https://www.histdata.com/f-a-q/
	// The timezone of all data is: Eastern Standard Time (EST) time-zone WITHOUT Day Light Savings adjustments.
	est := time.FixedZone("EST", -5*60*60) // -5 hours in seconds

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV row: %v", err)
		}
		if len(row) < 3 {
			return fmt.Errorf("expected at least 3 columns in CSV row, got %d: %v", len(row), row)
		}

		dtStr := row[0]
		bid, _ := strconv.ParseFloat(row[1], 64)
		ask, _ := strconv.ParseFloat(row[2], 64)

		// Add ms separator because go cannot parse without it
		splitIndex := len(dtStr) - 3
		dtStr = dtStr[:splitIndex] + "." + dtStr[splitIndex:]

		t, err := time.ParseInLocation("20060102 150405.000", dtStr, est)
		if err != nil {
			return fmt.Errorf("failed to parse date '%s': %v", dtStr, err)
		}

		tick := tick{
			Timestamp: t,
			Bid:       bid,
			Ask:       ask,
		}
		*arrayPtr = append(*arrayPtr, tick)
	}

	return nil
}

func loadData(beginDate time.Time, endDate time.Time, symbol string) ([]tick, error) {
	var ticks []tick

	// Loop through each month in the date range
	for d := beginDate; d.Before(endDate); d = d.AddDate(0, 1, 0) {
		year := d.Year()
		month := int(d.Month())

		err := loadFile(&ticks, year, month, symbol)
		if err != nil {
			return nil, fmt.Errorf("failed to load file for %s %d-%02d: %v", symbol, year, month, err)
		}
	}

	// Filter ticks within the date range
	var filteredTicks []tick
	for _, t := range ticks {
		if t.Timestamp.After(beginDate) && t.Timestamp.Before(endDate) {
			filteredTicks = append(filteredTicks, t)
		}
	}

	return filteredTicks, nil
}
