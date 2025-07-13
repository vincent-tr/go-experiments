package brokers

import "time"

type Timeframe time.Duration

const (
	Timeframe1Minute   Timeframe = Timeframe(1 * time.Minute)
	Timeframe5Minutes  Timeframe = Timeframe(5 * time.Minute)
	Timeframe15Minutes Timeframe = Timeframe(15 * time.Minute)
)

type Candle struct {
	Open  float64
	Close float64
	High  float64
	Low   float64
}

type PositionDirection int

const (
	// PositionDirectionLong means the position is a long position, i.e. buying low and selling high.
	PositionDirectionLong PositionDirection = iota

	// PositionDirectionShort means the position is a short position, i.e. selling high and buying low.
	PositionDirectionShort
)

func (d PositionDirection) String() string {
	switch d {
	case PositionDirectionLong:
		return "long"
	case PositionDirectionShort:
		return "short"
	default:
		return "unknown"
	}
}

// Order represents an order to enter a position in the market.
type Order struct {
	// Direction of the position (long or short)
	Direction PositionDirection

	// Number of lot to buy or sell
	// This is the number of lots, not the total amount of money invested.
	//
	// For example, if the lot size is 100 and Quantity is 10, and the price is 50,
	// the total amount of money invested is 10 * 100 * 50 = 50000.
	Quantity int

	// Price at which to stop loss the position
	StopLoss float64

	// Price at which to take profit on the position
	TakeProfit float64

	// Reason for the order
	Reason string
}

// Position represents a trading position in the market.
type Position interface {
	// Direction of the position (long or short)
	Direction() PositionDirection

	// Quantity of the position
	// This is the number of lots, not the total amount of money invested.
	Quantity() int

	// Price at which the position was opened
	OpenPrice() float64

	// Time at which the position was opened
	OpenTime() time.Time

	// Price at which the position was closed
	ClosePrice() float64

	// Time at which the position was closed
	CloseTime() time.Time

	// Whether the position is closed or not
	Closed() bool
}

// Broker is an interface that defines the methods required to interact with a trading broker.
// A broker is responsible for providing market data, executing orders, and managing the trading account.
type Broker interface {
	// Get the size of a single lot for the trading instrument.
	GetLotSize() int

	// Get the leverage for the trading account.
	GetLeverage() float64

	// Get the current capital of the trading account.
	GetCapital() float64

	// Register a callback to receive market data for a specific timeframe.
	RegisterMarketDataCallback(timeframe Timeframe, callback func(candle Candle))

	// Get the current time.
	// It is important to use this rather than time.Now() because when running in a backtest, the time may be simulated and not the real time.
	GetCurrentTime() time.Time

	// Place an order to enter a position in the market.
	PlaceOrder(order *Order) (Position, error)
}

// BacktestingBroker extends the Broker interface to include methods specific to backtesting scenarios.
type BacktestingBroker interface {
	Broker

	// Run the backtesting simulation.
	Run() error
}
