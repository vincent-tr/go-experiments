package indicators

import (
	"encoding/json"
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"

	"github.com/markcheno/go-talib"
)

func ATR(period int) Indicator {
	return newIndicator(
		func(ctx context.TraderContext) []float64 {
			history := ctx.HistoricalData()
			return talib.Atr(history.GetHighPrices(), history.GetLowPrices(), history.GetClosePrices(), period)
		},
		func() *formatter.FormatterNode {
			return formatter.Format("ATR", formatter.Format(fmt.Sprintf("Period: %d", period)))
		},
		func() (string, any) {
			return "atr", period
		},
	)
}

func init() {
	jsonParsers.RegisterParser("atr", func(arg json.RawMessage) (Indicator, error) {
		var period int
		if err := json.Unmarshal(arg, &period); err != nil {
			return nil, fmt.Errorf("failed to parse ATR period: %w", err)
		}

		return ATR(period), nil
	})
}
