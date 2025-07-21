package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type parquetTick struct {
	Timestamp int64   `parquet:"name=timestamp, type=INT64, convertedtype=TIMESTAMP_MILLIS"` // Store as milliseconds since epoch
	Bid       float64 `parquet:"name=bid, type=DOUBLE"`
	Ask       float64 `parquet:"name=ask, type=DOUBLE"`
}

const dataPath = "brokers/backtesting/data"

// https://www.histdata.com/download-free-forex-historical-data/?/ascii/tick-data-quotes/EURUSD

func loadCsvZip(zipFile string) ([]parquetTick, error) {

	// Unzip CSV
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open ZIP archive '%s': %v", zipFile, err)
	}
	defer r.Close()

	var csvFile io.ReadCloser

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".csv") {
			csvFile, err = f.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open CSV file '%s' in ZIP archive '%s': %v", f.Name, zipFile, err)
			}

			break
		}
	}

	if csvFile == nil {
		return nil, fmt.Errorf("no CSV file found in ZIP archive '%s'", zipFile)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = ','

	// https://www.histdata.com/f-a-q/
	// The timezone of all data is: Eastern Standard Time (EST) time-zone WITHOUT Day Light Savings adjustments.
	est := time.FixedZone("EST", -5*60*60) // -5 hours in seconds

	ticks := make([]parquetTick, 0)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %v", err)
		}
		if len(row) < 3 {
			return nil, fmt.Errorf("expected at least 3 columns in CSV row, got %d: %v", len(row), row)
		}

		dtStr := row[0]
		bid, _ := strconv.ParseFloat(row[1], 64)
		ask, _ := strconv.ParseFloat(row[2], 64)

		// Add ms separator because go cannot parse without it
		splitIndex := len(dtStr) - 3
		dtStr = dtStr[:splitIndex] + "." + dtStr[splitIndex:]

		t, err := time.ParseInLocation("20060102 150405.000", dtStr, est)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date '%s': %v", dtStr, err)
		}

		tick := parquetTick{
			Timestamp: t.UnixMilli(),
			Bid:       bid,
			Ask:       ask,
		}
		ticks = append(ticks, tick)
	}

	return ticks, nil
}

func writeParquet(filename string, ticks []parquetTick) error {
	// Create file
	fw, err := local.NewLocalFileWriter(filename)
	if err != nil {
		return err
	}
	defer fw.Close()

	// Create Parquet writer
	pw, err := writer.NewParquetWriter(fw, new(parquetTick), 4)
	if err != nil {
		return err
	}
	defer pw.WriteStop()

	pw.RowGroupSize = 128 * 1024 * 1024 // 128MB
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	// Write all ticks
	for _, tick := range ticks {
		if err := pw.Write(tick); err != nil {
			return err
		}
	}

	fmt.Printf("📊 Wrote %d ticks to %s\n", len(ticks), filename)
	return nil
}

func convertMissingParquetFiles() error {
	files, err := filepath.Glob(filepath.Join(dataPath, "HISTDATA_COM_ASCII_*.zip"))
	if err != nil {
		return fmt.Errorf("failed to list zip files: %v", err)
	}

	for _, zipFile := range files {
		base := filepath.Base(zipFile)
		base = strings.Replace(base, "HISTDATA_COM_ASCII_", "HISTDATA_COM_", 1)
		parquetName := strings.TrimSuffix(base, ".zip") + ".parquet"
		parquetPath := filepath.Join(dataPath, parquetName)

		if _, err := os.Stat(parquetPath); err == nil {
			fmt.Printf("✅ Parquet exists: %s (skipping)\n", parquetName)
			continue
		}

		fmt.Printf("📦 Converting: %s → %s\n", base, parquetName)

		ticks, err := loadCsvZip(zipFile)
		if err != nil {
			return fmt.Errorf("failed to load CSV: %v", err)
		}

		if err := writeParquet(parquetPath, ticks); err != nil {
			return fmt.Errorf("failed to write parquet: %v", err)
		}
	}

	return nil
}

func main() {
	// Convert all csv files where the target does not exist
	err := convertMissingParquetFiles()
	if err != nil {
		fmt.Printf("❌ Conversion failed: %v\n", err)
	}

}
