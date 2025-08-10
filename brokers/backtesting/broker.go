package backtesting

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/common"
	"math"
	"slices"
	"time"
)

var log = common.NewLogger("backtesting")

type Config struct {
	LotSize        int     // Size of the lot to trade
	Leverage       float64 // Leverage to use for trading
	InitialCapital float64 // Initial capital for the backtesting account
}

type Metrics struct {
	// TotalTrades is the number of trades taken during the test period.
	// Important to ensure statistical significance.
	TotalTrades int

	// WinRate is the percentage of trades that were profitable.
	// Calculated as (WinningTrades / TotalTrades) * 100.
	WinRate float64 // in percent

	// NetPnL is the total profit or loss at the end of the test period.
	// Usually expressed in base currency (e.g., USD).
	NetPnL float64

	// ProfitFactor is the ratio of gross profit to gross loss.
	// A value > 1 means the system was profitable overall.
	ProfitFactor float64

	// MaxDrawdownPct is the largest percentage drop from a peak in equity.
	// Shows the worst-case capital exposure during the test.
	MaxDrawdownPct float64 // in percent

	// ExpectedValueR is the average return per trade in R-multiples.
	// Helps understand the return relative to risk.
	ExpectedValueR float64

	// AvgTradeDuration is the average duration a trade was open.
	// Useful for evaluating intraday vs swing behavior.
	AvgTradeDuration time.Duration

	// LongTrades is the number of trades taken in the long (buy) direction.
	LongTrades int

	// ShortTrades is the number of trades taken in the short (sell) direction.
	ShortTrades int
}

type broker struct {
	config           *Config
	ticks            []tick
	currentIndex     int
	capital          float64
	openPositions    map[*position]struct{}
	callbacks        map[brokers.Timeframe][]func(candle brokers.Candle)
	positionsHistory []*position
}

// Run implements brokers.BacktestingBroker.
func (b *broker) Run() error {
	log.Debug("🚀 Starting backtest with %d ticks and initial capital %.2f", len(b.ticks), b.capital)

	for {
		b.processTick()

		if b.currentIndex == len(b.ticks)-1 {
			break
		}

		b.currentIndex++
	}

	b.closeAllOpenPositions()

	log.Debug("✅ Backtest completed.")
	// b.printSummary()

	return nil
}

// GetLotSize implements brokers.Broker.
func (b *broker) GetLotSize() int {
	return b.config.LotSize
}

// GetCapital implements brokers.Broker.
func (b *broker) GetCapital() float64 {
	return b.capital
}

// GetLeverage implements brokers.Broker.
func (b *broker) GetLeverage() float64 {
	return b.config.Leverage
}

// GetCurrentTime implements brokers.Broker.
func (b *broker) GetCurrentTime() time.Time {
	return b.currentTick().Timestamp
}

// GetMarketDataChannel implements brokers.Broker.
func (b *broker) RegisterMarketDataCallback(timeframe brokers.Timeframe, callback func(candle brokers.Candle)) {
	b.callbacks[timeframe] = append(b.callbacks[timeframe], callback)
}

// PlaceOrder implements brokers.Broker.
func (b *broker) PlaceOrder(order *brokers.Order) (brokers.Position, error) {
	pos := newPosition(b.currentTick(), b.GetCapital(), order)
	margin := pos.getMargin(b.GetLeverage())

	if margin > b.capital {
		return nil, fmt.Errorf("insufficient capital: cannot place order for %d lots at price %.4f (margin: %.2f, capital:  %.2f)", pos.Quantity(), pos.OpenPrice(), margin, b.capital)
	}

	b.capital -= margin
	b.openPositions[pos] = struct{}{}
	b.positionsHistory = append(b.positionsHistory, pos)

	log.Debug("📈 Placed order: Direction=%s, Quantity=%d, OpenPrice=%.5f, StopLoss=%.5f, TakeProfit=%.5f, Reason=%s",
		pos.Direction(), pos.Quantity(), pos.openPrice, order.StopLoss, order.TakeProfit,
		order.Reason)

	return pos, nil
}

var _ brokers.Broker = (*broker)(nil)
var _ brokers.BacktestingBroker = (*broker)(nil)

// NewBroker creates a new instance of the broker.
func NewBroker(config *Config, dataset *Dataset) (brokers.BacktestingBroker, error) {
	b := &broker{
		config:           config,
		ticks:            dataset.ticks,
		currentIndex:     0,
		capital:          config.InitialCapital,
		openPositions:    make(map[*position]struct{}),
		callbacks:        make(map[brokers.Timeframe][]func(candle brokers.Candle)),
		positionsHistory: make([]*position, 0),
	}

	return b, nil
}

func ComputeMetrics(b brokers.BacktestingBroker) (map[common.Month]*Metrics, error) {
	bb, ok := b.(*broker)
	if !ok {
		return nil, fmt.Errorf("invalid broker type: expected *broker, got %T", b)
	}

	metrics := bb.computeMetrics()
	return metrics, nil
}

func (b *broker) currentTick() *tick {
	return &b.ticks[b.currentIndex]
}

func (b *broker) printGap() {
	currentTick := b.currentTick()
	if !currentTick.IsGap || b.currentIndex == 0 {
		return
	}
	previousTick := b.ticks[b.currentIndex-1]
	if !previousTick.IsGap {
		return
	}

	log.Warning("⏳ Gap detected at %s: Previous=%s, Difference=%s",
		currentTick.Timestamp.Format("2006-01-02 15:04:05"),
		previousTick.Timestamp.Format("2006-01-02 15:04:05"),
		currentTick.Timestamp.Sub(previousTick.Timestamp).String())
}

func (b *broker) processTick() {
	currentTick := b.currentTick()
	// b.printGap()
	// log.Debug("📈 Processing tick at %s: Bid=%.5f, Ask=%.5f", currentTick.Timestamp.Format("2006-01-02 15:04:05"), currentTick.Bid, currentTick.Ask)

	if currentTick.IsGap {
		b.cancelAllOpenPositions()
	}

	for pos := range b.openPositions {
		switch pos.isTriggered(currentTick) {
		case CloseTriggerNone:
			// Position is still open, do nothing
			continue
		case CloseTriggerStopLoss, CloseTriggerTakeProfit:
			// Position should be closed
			b.closePosition(pos)

			closeReason := "unknown"
			if pos.isTriggered(currentTick) == CloseTriggerStopLoss {
				closeReason = "stop loss"
			} else if pos.isTriggered(currentTick) == CloseTriggerTakeProfit {
				closeReason = "take profit"
			}

			log.Debug("📉 Position closed (%s) at %s: Direction=%s, Quantity=%d, OpenPrice=%.5f, ClosePrice=%.5f",
				closeReason,
				currentTick.Timestamp.Format("2006-01-02 15:04:05"),
				pos.direction, pos.quantity, pos.openPrice, pos.closePrice)
		}
	}

	// Check if we have a full candle for any registered timeframes
	for timeframe, callbacks := range b.callbacks {
		candle := b.tryCandle(timeframe)

		if candle != nil {
			// log.Debug("📊 New candle for timeframe %s: Open=%.5f, Close=%.5f, High=%.5f, Low=%.5f",
			// 	timeframe, candle.Open, candle.Close, candle.High, candle.Low)

			// Call all registered callbacks for this timeframe
			for _, callback := range callbacks {
				callback(*candle)
			}
		}
	}

}

func (b *broker) tryCandle(timeframe brokers.Timeframe) *brokers.Candle {
	currentTick := b.currentTick()

	var nextTick *tick
	if b.currentIndex+1 < len(b.ticks) {
		nextTick = &b.ticks[b.currentIndex+1]
	}

	currentBucket := getTimeframeBucket(currentTick, timeframe)

	// Is the next tick in a different timeframe?
	if nextTick != nil && getTimeframeBucket(nextTick, timeframe) == currentBucket {
		// No complete candle yet, we need to wait for the next tick
		return nil
	}

	// We have the last tick of the current timeframe

	// Get all ticks for the current timeframe bucket
	timeframeTicks := []*tick{currentTick}
	for i := b.currentIndex - 1; i >= 0; i-- {
		if getTimeframeBucket(&b.ticks[i], timeframe) == currentBucket {
			timeframeTicks = append(timeframeTicks, &b.ticks[i])
		} else {
			break
		}
	}

	slices.Reverse(timeframeTicks)

	// Create a candle from the timeframe ticks
	usable := true
	low := timeframeTicks[0].Price()
	high := timeframeTicks[0].Price()
	for _, t := range timeframeTicks {
		price := t.Price()
		if price < low {
			low = price
		}
		if price > high {
			high = price
		}
		if t.IsGap {
			usable = false // If any tick is a gap, the candle is not usable
		}
	}

	return &brokers.Candle{
		Open:   timeframeTicks[0].Price(),
		Close:  timeframeTicks[len(timeframeTicks)-1].Price(),
		High:   high,
		Low:    low,
		Usable: usable,
	}
}

func getTimeframeBucket(tick *tick, timeframe brokers.Timeframe) string {
	// This function should return the start time of the bucket for the given timeframe.
	// For simplicity, we assume that the tick's timestamp is already aligned with the timeframe.
	return tick.Timestamp.Truncate(time.Duration(timeframe)).Format("2006-01-02 15:04:05")
}

func (b *broker) cancelAllOpenPositions() {
	for pos := range b.openPositions {
		pos.cancelPosition()
		delete(b.openPositions, pos)
		b.positionsHistory = slices.DeleteFunc(b.positionsHistory, func(p *position) bool {
			return p == pos
		})

		b.capital += pos.getMargin(b.GetLeverage()) // Return margin to capital
		// Note: We do not add profit/loss here because the position is canceled, not closed.

		log.Debug("📉 Position canceled at %s: Direction=%s, Quantity=%d, OpenPrice=%.5f",
			b.currentTick().Timestamp.Format("2006-01-02 15:04:05"),
			pos.direction, pos.quantity, pos.openPrice)
	}
}

func (b *broker) closeAllOpenPositions() {
	for pos := range b.openPositions {
		b.closePosition(pos)

		log.Debug("📉 Position closed (end of test) at %s: Direction=%s, Quantity=%d, OpenPrice=%.5f, ClosePrice=%.5f",
			b.currentTick().Timestamp.Format("2006-01-02 15:04:05"),
			pos.direction, pos.quantity, pos.openPrice, pos.closePrice)
	}
}

func (b *broker) closePosition(pos *position) {
	pos.closePosition(b.currentTick())
	delete(b.openPositions, pos)

	b.capital += pos.getMargin(b.GetLeverage())
	b.capital += pos.getProfitAndLoss()
}

func (b *broker) printSummary() {
	log.Info("📊 Backtest Summary:")

	log.Debug("Positions history:")
	var currentMonth string
	var monthProfit float64
	monthCapital := b.config.InitialCapital
	for _, pos := range b.positionsHistory {
		monthKey := pos.openTime.Format("2006-01")
		if monthKey != currentMonth {
			if currentMonth != "" {
				log.Info("📅 Month: %s, Initial capital: %.2f, Profit: %s", currentMonth, monthCapital, b.formatProfit(monthProfit))
			}

			monthCapital += monthProfit
			monthProfit = 0 // Reset for new month
			currentMonth = monthKey
			log.Debug("📅 Month: %s", monthKey)
		}

		profit := pos.getProfitAndLoss()
		monthProfit += profit

		log.Debug(" - Capital: %0.2f, Direction: %s, OpenTime: %s, Profit: %s, Duration: %s",
			pos.capital, pos.direction, pos.openTime.Format("2006-01-02 15:04:05"), b.formatProfit(profit), pos.CloseTime().Sub(pos.OpenTime()).String())
	}

	log.Info("📅 Month: %s, Initial capital: %.2f, Profit: %s", currentMonth, monthCapital, b.formatProfit(monthProfit))
	monthProfit = 0 // Reset for new month

	log.Info("Total positions: %d", len(b.positionsHistory))
	log.Info("Final capital: %.2f", b.capital)

	totalProfit := b.capital - b.config.InitialCapital
	ratio := totalProfit / b.config.InitialCapital * 100
	log.Info("Total profit/loss: %s (%s%%)", b.formatProfit(b.capital-b.config.InitialCapital), b.formatProfit(ratio))
}

func (b *broker) formatProfit(value float64) string {
	if value < 0 {
		return fmt.Sprintf("\033[31m%.2f\033[0m", value) // Red for losses
	} else if value > 0 {
		return fmt.Sprintf("\033[32m%.2f\033[0m", value) // Green for profits
	} else {
		return fmt.Sprintf("%.2f", value) // No color for zero
	}
}

func (b *broker) computeMetrics() map[common.Month]*Metrics {
	positionsByMonth := make(map[common.Month][]*position)

	// Group positions by month
	for _, pos := range b.positionsHistory {
		month := common.FromDate(pos.openTime)
		positionsByMonth[month] = append(positionsByMonth[month], pos)
	}

	// Compute metrics for each month
	metrics := make(map[common.Month]*Metrics)
	for month, positions := range positionsByMonth {
		monthlyMetrics := b.computeMonthlyMetrics(positions)
		metrics[month] = monthlyMetrics
	}

	return metrics
}

func (b *broker) computeMonthlyMetrics(positions []*position) *Metrics {
	var totalTrades, winningTrades, longTrades, shortTrades int
	var netPnL, grossProfit, grossLoss, totalR, maxR float64
	var totalDuration time.Duration

	equity := 0.0
	peakEquity := 0.0
	maxDrawdown := 0.0

	metrics := &Metrics{}

	for _, pos := range positions {
		if !pos.closed {
			continue
		}

		totalTrades++

		// Trade direction
		switch pos.direction {
		case brokers.PositionDirectionLong:
			longTrades++
		case brokers.PositionDirectionShort:
			shortTrades++
		}

		// PnL
		pnl := pos.getProfitAndLoss()
		netPnL += pnl
		equity += pnl

		// Track peak equity and drawdown
		if equity > peakEquity || totalTrades == 1 {
			peakEquity = equity
		}
		drawdown := peakEquity - equity
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}

		// R-multiple
		risk := math.Abs(pos.openPrice - pos.stopLoss)
		if risk > 0 {
			r := pnl / (risk * float64(pos.quantity))
			totalR += r
			if r > maxR {
				maxR = r
			}
		}

		// Profit stats
		if pnl > 0 {
			winningTrades++
			grossProfit += pnl
		} else {
			grossLoss += -pnl
		}

		// Duration
		duration := pos.closeTime.Sub(pos.openTime)
		totalDuration += duration
	}

	metrics.TotalTrades = totalTrades
	metrics.NetPnL = netPnL
	metrics.LongTrades = longTrades
	metrics.ShortTrades = shortTrades

	if totalTrades > 0 {
		metrics.WinRate = float64(winningTrades) / float64(totalTrades) * 100
		metrics.ExpectedValueR = totalR / float64(totalTrades)
		metrics.AvgTradeDuration = totalDuration / time.Duration(totalTrades)
	}
	if grossLoss > 0 {
		metrics.ProfitFactor = grossProfit / grossLoss
	}
	if peakEquity != 0 {
		metrics.MaxDrawdownPct = (maxDrawdown / peakEquity) * 100
	}

	return metrics
}
