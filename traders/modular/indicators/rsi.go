package indicators

import (
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"

	"github.com/markcheno/go-talib"
)

func RSI(period int) Indicator {
	return newIndicator(
		func(ctx context.TraderContext) []float64 {
			closePrices := ctx.HistoricalData().GetClosePrices()
			return talib.Rsi(closePrices, period)
		},
		func() *formatter.FormatterNode {
			return formatter.Format("RSI",
				formatter.Format(fmt.Sprintf("Period: %d", period)),
			)
		},
	)
}
