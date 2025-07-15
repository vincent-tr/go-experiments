package gpt

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/common"
	"go-experiments/traders/tools"
	"math"
	"time"

	"github.com/markcheno/go-talib"
)

var log = common.NewLogger("traders/gpt")

const pipSize = 0.0001

type Config struct {
	HistorySize int // Size of the history buffer for technical indicators

	// Position management

	EmaFastPeriod int     // Fast EMA period
	EmaSlowPeriod int     // Slow EMA period
	RsiPeriod     int     // RSI period
	RsiMin        float64 // RSI neutral zone (30-70)
	RsiMax        float64 // RSI neutral zone (30-70)
	AdxEnabled    bool    // Whether to use ADX for trend strength confirmation
	AdxPeriod     int     // ADX period
	AdxThreshold  float64 // ADX threshold to confirm trend strength (e.g., 20-25)

	// Stop-loss and take-profit
	StopLossPipBuffer    int // Number of pips to buffer for stop-loss
	StopLossLookupPeriod int // Number of minutes to look back for stop-loss calculation

	StopLossAtrEnabled    bool    // Whether to use ATR for stop-loss calculation
	StopLossAtrPeriod     int     // ATR period for stop-loss calculation
	StopLossAtrMultiplier float64 // Multiplier for ATR-based stop-loss

	TakeProfitRatio float64 // Ratio of risk to reward for take-profit

	// Capital risk management

	CapitalRiskPercent float64 // Percentage of capital to risk per trade
}

func Setup(broker brokers.Broker, config *Config) {

	trader := newTrader(broker, config)

	broker.RegisterMarketDataCallback(brokers.Timeframe1Minute, func(candle brokers.Candle) {
		trader.tick(candle)
	})
}

type trader struct {
	broker       brokers.Broker
	config       *Config
	history      *tools.History
	openPosition brokers.Position
}

func newTrader(broker brokers.Broker, config *Config) *trader {

	return &trader{
		broker:  broker,
		config:  config,
		history: tools.NewHistory(config.HistorySize),
	}
}

func (t *trader) tick(candle brokers.Candle) {
	t.history.AddCandle(candle)

	if !t.history.IsComplete() {
		log.Debug("Not enough data to make a decision")
		return
	}

	// Check if we have an open position
	if t.openPosition != nil {
		if t.openPosition.Closed() {
			t.openPosition = nil
		}
	}

	// Only take one position at a time
	if t.openPosition != nil {
		return
	}

	// Check if we should trade regarding the calandar (holidays, weekday, session hours, etc.)
	if !t.shouldTrade() {
		return
	}

	res, direction := t.shouldTakePosition()
	if !res {
		return
	}

	entryPrice := candle.Close
	stopLoss := t.computeStopLoss(direction)
	takeProfit := t.computeTakeProfit(direction, entryPrice, stopLoss)
	positionSize := t.computePositionSize(stopLoss)

	if positionSize == 0 {
		// Not enough capital to take a position
		return
	}

	order := &brokers.Order{
		Direction:  direction,
		Quantity:   positionSize,
		StopLoss:   stopLoss,
		TakeProfit: takeProfit,
		Reason:     fmt.Sprintf("GPT strategy: %s at %.5f", direction, entryPrice),
	}

	position, err := t.broker.PlaceOrder(order)
	if err != nil {
		log.Error("Failed to place order: %v", err)
	}

	t.openPosition = position
}

func (t *trader) shouldTakePosition() (bool, brokers.PositionDirection) {
	var defaultValue brokers.PositionDirection

	closePrices := t.history.GetClosePrices()
	last := len(closePrices) - 1

	// RSI must be between RsiMin and RsiMax (neutral zone)
	rsi := talib.Rsi(closePrices, t.config.RsiPeriod)
	currRSI := rsi[last]

	if currRSI < t.config.RsiMin || currRSI > t.config.RsiMax {
		return false, defaultValue
	}

	// If ADX is enabled, check if the trend is strong enough
	if t.config.AdxEnabled {
		highPrices := t.history.GetHighPrices()
		lowPrices := t.history.GetLowPrices()
		adx := talib.Adx(highPrices, lowPrices, closePrices, t.config.AdxPeriod)
		currAdx := adx[last]

		if currAdx < t.config.AdxThreshold {
			return false, defaultValue
		}
	}

	emaSlow := talib.Ema(closePrices, t.config.EmaSlowPeriod)
	emaFast := talib.Ema(closePrices, t.config.EmaFastPeriod)
	prevFast := emaFast[last-1]
	prevSlow := emaSlow[last-1]
	currFast := emaFast[last]
	currSlow := emaSlow[last]

	// Buy signal: bullish crossover
	if prevFast < prevSlow && currFast > currSlow {
		return true, brokers.PositionDirectionLong
	}

	// Sell signal: bearish crossover
	if prevFast > prevSlow && currFast < currSlow {
		return true, brokers.PositionDirectionShort
	}

	return false, defaultValue
}

func (t *trader) shouldTrade() bool {
	currentTime := t.broker.GetCurrentTime()

	weekday := currentTime.Weekday()
	if weekday < time.Tuesday || weekday > time.Thursday {
		return false
	}

	if common.IsUSHoliday(currentTime) || common.IsUKHoliday(currentTime) {
		return false
	}

	if !common.LondonSession.IsOpen(currentTime) || !common.NYSession.IsOpen(currentTime) {
		return false
	}

	return true
}

// Computes the stop-loss price based on the last lookupPeriod minutes of candles.
// For long positions, it is set 3 pips below the lowest low in the last lookupPeriod minutes.
// For short positions, it is set 3 pips above the highest high in the last lookupPeriod minutes.
func (t *trader) computeStopLoss(direction brokers.PositionDirection) float64 {

	if t.config.StopLossAtrEnabled {
		// Use ATR for stop-loss calculation
		highPrices := t.history.GetHighPrices()
		lowPrices := t.history.GetLowPrices()
		closePrices := t.history.GetClosePrices()

		atr := talib.Atr(highPrices, lowPrices, closePrices, t.config.StopLossAtrPeriod)
		last := len(atr) - 1
		currAtr := atr[last]

		pipDistance := currAtr * t.config.StopLossAtrMultiplier
		entryPrice := t.history.GetPrice()

		switch direction {
		case brokers.PositionDirectionLong:
			return entryPrice - pipDistance
		case brokers.PositionDirectionShort:
			return entryPrice + pipDistance
		default:
			panic("invalid position type")
		}
	}

	pipDistance := float64(t.config.StopLossPipBuffer) * pipSize
	lookupPeriod := t.config.StopLossLookupPeriod

	switch direction {
	case brokers.PositionDirectionLong:
		// find lowest low in last lookupPeriod minutes
		lowest := t.history.GetLowest(lookupPeriod)
		// stop loss is 3 pips below that low
		return lowest - pipDistance

	case brokers.PositionDirectionShort:
		// find highest high in last lookupPeriod minutes
		highest := t.history.GetHighest(lookupPeriod)
		// stop loss is 3 pips above that high
		return highest + pipDistance

	default:
		panic("invalid position direction: " + direction.String())
	}
}

// The take-profit is set at a takeProfitRatio reward-to-risk ratio relative to your stop-loss distance.
func (t *trader) computeTakeProfit(direction brokers.PositionDirection, entryPrice, stopLoss float64) float64 {
	switch direction {
	case brokers.PositionDirectionLong:
		risk := entryPrice - stopLoss
		if risk <= 0 {
			panic(fmt.Sprintf("invalid stoploss for long position: entryPrice=%.5f, stopLoss=%.5f", entryPrice, stopLoss))
		}
		return entryPrice + t.config.TakeProfitRatio*risk

	case brokers.PositionDirectionShort:
		risk := stopLoss - entryPrice
		if risk <= 0 {
			panic(fmt.Sprintf("invalid stoploss for short position: entryPrice=%.5f, stopLoss=%.5f", entryPrice, stopLoss))
		}
		return entryPrice - t.config.TakeProfitRatio*risk

	default:
		panic("invalid position direction: " + direction.String())
	}
}

func (t *trader) computePositionSize(stopLoss float64) int {
	accountBalance := t.broker.GetCapital()
	accountRisk := accountBalance * (t.config.CapitalRiskPercent / 100)

	entryPrice := t.history.GetPrice()
	priceDiff := math.Abs(entryPrice - stopLoss)
	if priceDiff <= 0 {
		panic(fmt.Sprintf("Invalid stop loss price: entryPrice=%.5f, stopLoss=%.5f", entryPrice, stopLoss))
	}

	lotSize := float64(t.broker.GetLotSize())
	riskPerLot := lotSize * priceDiff
	positionSize := accountRisk / riskPerLot

	// Ensure position size doesn't exceed account balance
	// Total value = positionSize * lotSize * entryPrice
	maxPositionSize := accountBalance * t.broker.GetLeverage() / (lotSize * entryPrice)
	if positionSize > maxPositionSize {
		positionSize = maxPositionSize
	}

	return int(math.Floor(positionSize))
}
