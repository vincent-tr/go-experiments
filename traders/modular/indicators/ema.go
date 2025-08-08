package indicators

import (
	"encoding/json"
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"

	"github.com/markcheno/go-talib"
)

func EMA(period int) Indicator {
	return newIndicator(
		func(ctx context.TraderContext) []float64 {
			closePrices := ctx.HistoricalData().GetClosePrices()
			return talib.Ema(closePrices, period)
		},
		func() *formatter.FormatterNode {
			return formatter.Format("EMA",
				formatter.Format(fmt.Sprintf("Period: %d", period)),
			)
		},
		func() (string, any) {
			return "ema", period
		},
	)
}

func init() {
	jsonParsers.RegisterParser("ema", func(arg json.RawMessage) (Indicator, error) {
		var period int
		if err := json.Unmarshal(arg, &period); err != nil {
			return nil, fmt.Errorf("failed to parse EMA period: %w", err)
		}

		return EMA(period), nil
	})
}
