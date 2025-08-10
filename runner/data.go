package runner

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"go-experiments/brokers/backtesting"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "data.db"

type run struct {
	// Config
	Key        string // hash of next fields
	Instrument string
	TimeRange  string
	Strategy   string

	// Results
	backtesting.Metrics
}

type Database struct {
	db *sql.DB
}

func OpenDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Create table if it doesn't exist
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS runs (
        -- Config
        key TEXT PRIMARY KEY,                -- Unique hash for config
        instrument TEXT NOT NULL,            -- e.g., EURUSD
        time_range TEXT NOT NULL,             -- e.g., 202501
        strategy TEXT NOT NULL,               -- Serialized strategy description (JSON)

        -- Metrics
        total_trades INTEGER NOT NULL,        -- Total trades
        win_rate REAL NOT NULL,               -- % of winning trades
        net_pnl REAL NOT NULL,                 -- Net profit/loss in base currency
        profit_factor REAL NOT NULL,           -- Gross profit / gross loss
        max_drawdown_pct REAL NOT NULL,        -- % from peak equity
        expected_value_r REAL NOT NULL,        -- Avg R-multiple return
        avg_trade_duration_seconds INTEGER NOT NULL, -- Duration in seconds
        long_trades INTEGER NOT NULL,          -- Count of long trades
        short_trades INTEGER NOT NULL          -- Count of short trades
    );`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

func (db *Database) ComputeKey(instrument, timeRange, strategy string) string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s:%s:%s", instrument, timeRange, strategy)))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// return nil if run does not exist
func (db *Database) FindRun(instrument, timeRange, strategy string) (*run, error) {
	key := db.ComputeKey(instrument, timeRange, strategy)

	query := `
    SELECT
        key,
        instrument,
        time_range,
        strategy,
        total_trades,
        win_rate,
        net_pnl,
        profit_factor,
        max_drawdown_pct,
        expected_value_r,
        avg_trade_duration_seconds,
        long_trades,
        short_trades
    FROM runs
    WHERE key = ?;`

	var r run
	var tradeDurationSeconds int64

	err := db.db.QueryRow(query, key).Scan(
		&r.Key, &r.Instrument, &r.TimeRange, &r.Strategy,
		&r.TotalTrades, &r.WinRate, &r.NetPnL,
		&r.ProfitFactor, &r.MaxDrawdownPct,
		&r.ExpectedValueR, &tradeDurationSeconds,
		&r.LongTrades, &r.ShortTrades,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Run does not exist
	} else if err != nil {
		return nil, err // Other error
	}

	r.AvgTradeDuration = time.Second * time.Duration(tradeDurationSeconds)

	return &r, nil
}

func (db *Database) SaveRun(r *run) error {
	key := db.ComputeKey(r.Instrument, r.TimeRange, r.Strategy)
	tradeDurationSeconds := int64(r.AvgTradeDuration.Seconds())

	// Insert or update the run
	query := `
    INSERT INTO runs (
        key, instrument, time_range, strategy,
        total_trades, win_rate, net_pnl,
        profit_factor, max_drawdown_pct,
        expected_value_r, avg_trade_duration_seconds,
        long_trades, short_trades
    ) VALUES (?, ?, ?, ?,
        ?, ?, ?,
        ?, ?, ?,
        ?, ?, ?
    );`

	_, err := db.db.Exec(query, key, r.Instrument, r.TimeRange, r.Strategy,
		r.TotalTrades, r.WinRate, r.NetPnL,
		r.ProfitFactor, r.MaxDrawdownPct,
		r.ExpectedValueR, tradeDurationSeconds,
		r.LongTrades, r.ShortTrades,
	)

	return err
}
