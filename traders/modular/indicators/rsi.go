package indicators

import (
	"encoding/json"
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

func init() {
	jsonParsers.RegisterParser("rsi", func(arg json.RawMessage) (Indicator, error) {
		var period int
		if err := json.Unmarshal(arg, &period); err != nil {
			return nil, fmt.Errorf("failed to parse RSI period: %w", err)
		}

		return RSI(period), nil
	})
}
