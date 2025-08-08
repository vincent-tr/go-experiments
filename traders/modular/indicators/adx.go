package indicators

import (
	"encoding/json"
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"

	"github.com/markcheno/go-talib"
)

func ADX(period int) Indicator {
	return newIndicator(
		func(ctx context.TraderContext) []float64 {
			history := ctx.HistoricalData()
			return talib.Adx(history.GetHighPrices(), history.GetLowPrices(), history.GetClosePrices(), period)
		},
		func() *formatter.FormatterNode {
			return formatter.Format("ADX", formatter.Format(fmt.Sprintf("Period: %d", period)))
		},
	)
}

func init() {
	jsonParsers.RegisterParser("adx", func(arg json.RawMessage) (Indicator, error) {
		var period int
		if err := json.Unmarshal(arg, &period); err != nil {
			return nil, fmt.Errorf("failed to parse ADX period: %w", err)
		}

		return ADX(period), nil
	})
}
