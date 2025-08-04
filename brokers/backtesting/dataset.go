package backtesting

import (
	"fmt"
	"go-experiments/common"
	"path"
	"runtime"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/source"
)

const dataPath = "brokers/backtesting/data"

const MaxGap = time.Minute // Maximum allowed gap between ticks

// https://www.histdata.com/download-free-forex-historical-data/?/ascii/tick-data-quotes/EURUSD

type Dataset struct {
	ticks     []tick
	symbol    string
	beginDate time.Time
	endDate   time.Time
}

func (d *Dataset) Symbol() string {
	return d.symbol
}

func (d *Dataset) BeginDate() time.Time {
	return d.beginDate
}

func (d *Dataset) EndDate() time.Time {
	return d.endDate
}

func (d *Dataset) TickCount() int {
	return len(d.ticks)
}

func LoadDataset(begin, end common.Month, symbol string) (*Dataset, error) {
	beginTime := time.Now()

	files := make([]*file, 0) // Preallocate for 12 months
	beginDate := begin.FirstDay()
	endDate := end.LastDay()

	for d := beginDate; d.Before(endDate); d = d.AddDate(0, 1, 0) {
		f, err := openFile(d.Year(), int(d.Month()), symbol)
		if err != nil {
			return nil, err
		}

		defer f.Close()
		files = append(files, f)
	}

	tickCount := 0
	for _, f := range files {
		tickCount += f.TickCount()
	}

	ticks := make([]tick, tickCount)
	offset := 0
	for _, f := range files {
		if err := f.ReadTicks(ticks, offset); err != nil {
			return nil, fmt.Errorf("failed to read ticks from file: %v", err)
		}
		offset += f.TickCount()
	}

	// Mark gaps in the data
	markGaps(ticks)

	endTime := time.Now()
	duration := endTime.Sub(beginTime)
	log.Debug("⏱️  Read %d ticks from %d file(s) in %s.", tickCount, len(files), duration)
	log.Info("📈 Loaded dataset from %s to %s", begin.String(), end.String())

	return &Dataset{ticks: ticks}, nil
}

func markGaps(ticks []tick) {
	for i := 1; i < len(ticks)-1; i++ {
		previous := &ticks[i-1]
		current := &ticks[i]

		if current.Timestamp.Sub(previous.Timestamp) > MaxGap {
			current.IsGap = true
			previous.IsGap = true
		}
	}
}

// Tick represents one row of tick data
type tick struct {
	Timestamp time.Time
	Bid       float64
	Ask       float64
	IsGap     bool // Indicates if there is a gap in the data before or after this tick
}

// Use intermediate struct with int64 timestamp
type parquetTick struct {
	Timestamp int64   `parquet:"name=timestamp, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
	Bid       float64 `parquet:"name=bid, type=DOUBLE"`
	Ask       float64 `parquet:"name=ask, type=DOUBLE"`
}

func (t *tick) Price() float64 {
	// For simplicity, we return the average of bid and ask as the price.
	// In a real implementation, you might want to use bid or ask based on your strategy.
	return (t.Bid + t.Ask) / 2
}

type file struct {
	pFile  source.ParquetFile
	reader *reader.ParquetReader
}

func openFile(year int, month int, symbol string) (*file, error) {
	parquetFile := path.Join(dataPath, fmt.Sprintf("HISTDATA_COM_%s_T%04d%02d.parquet", symbol, year, month))

	// Open Parquet file
	pFile, err := local.NewLocalFileReader(parquetFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open Parquet file '%s': %v", parquetFile, err)
	}

	reader, err := reader.NewParquetReader(pFile, new(parquetTick), int64(runtime.NumCPU()))
	if err != nil {
		pFile.Close()
		return nil, fmt.Errorf("failed to create Parquet reader for '%s': %v", parquetFile, err)
	}

	return &file{pFile, reader}, nil
}

func (f *file) Close() error {
	f.reader.ReadStop()
	return f.pFile.Close()
}

func (f *file) TickCount() int {
	return int(f.reader.GetNumRows())
}

func (f *file) ReadTicks(array []tick, offset int) error {

	rows := make([]parquetTick, f.TickCount())
	if err := f.reader.Read(&rows); err != nil {
		return fmt.Errorf("failed to read Parquet rows: %v", err)
	}
	for i, r := range rows {
		array[offset+i] = tick{
			Timestamp: time.UnixMilli(r.Timestamp),
			Bid:       r.Bid,
			Ask:       r.Ask,
			IsGap:     false, // Default to false, will be updated later
		}
	}

	return nil
}
