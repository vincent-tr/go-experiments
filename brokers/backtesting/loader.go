package backtesting

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"path"
	"strconv"
	"strings"
	"time"
)

const dataPath = "brokers/backtesting"

// https://www.histdata.com/download-free-forex-historical-data/?/ascii/tick-data-quotes/EURUSD

// Tick represents one row of tick data
type Tick struct {
	Timestamp time.Time
	Bid       float64
	Ask       float64
}

func loadFile(arrayPtr *[]Tick, year int, month int, symbol string) error {

	zipFile := path.Join(dataPath, fmt.Sprintf("HISTDATA_COM_ASCII_%s_T%04d%02d.zip", symbol, year, month))

	// Unzip CSV
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var csvFile io.ReadCloser

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".csv") {
			csvFile, err = f.Open()
			if err != nil {
				log.Fatal(err)
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

		t, err := time.Parse("20060102 150405.000", dtStr)
		if err != nil {
			return fmt.Errorf("failed to parse date '%s': %v", dtStr, err)
		}

		tick := Tick{
			Timestamp: t,
			Bid:       bid,
			Ask:       ask,
		}
		*arrayPtr = append(*arrayPtr, tick)
	}

	return nil
}

func loadData(beginDate time.Time, endDate time.Time, symbol string) ([]Tick, error) {
	var ticks []Tick

	// Loop through each month in the date range
	for d := beginDate; d.Before(endDate); d = d.AddDate(0, 1, 0) {
		year := d.Year()
		month := int(d.Month())

		err := loadFile(&ticks, year, month, symbol)
		if err != nil {
			return nil, fmt.Errorf("failed to load file for %s %d-%02d: %v", symbol, year, month, err)
		}
	}

	return ticks, nil
}
