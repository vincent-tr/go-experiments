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

// Order represents an order to enter a position in the market.
type Order struct {
	// Direction of the position (long or short)
	Direction PositionDirection

	// Amount of the asset to buy or sell
	Amount float64

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

	// Channel to signal when the position is closed
	ClosedSignal() <-chan struct{}
}

// Broker is an interface that defines the methods required to interact with a trading broker.
// A broker is responsible for providing market data, executing orders, and managing the trading account.
type Broker interface {
	// Get the current capital of the trading account.
	GetCapital() float64

	// Get the current market data channel.
	GetMarketDataChannel(timeframe Timeframe) <-chan Candle

	// Get the current time.
	// It is important to use this rather than time.Now() because when running in a backtest, the time may be simulated and not the real time.
	GetCurrentTime() time.Time

	// Place an order to enter a position in the market.
	PlaceOrder(order *Order) (Position, error)
}
