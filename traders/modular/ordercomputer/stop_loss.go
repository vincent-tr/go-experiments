package ordercomputer

import (
	"encoding/json"
	"fmt"
	"go-experiments/brokers"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/indicators"
	"go-experiments/traders/modular/marshal"
)

func StopLossATR(atr indicators.Indicator, multiplier float64) OrderComputer {
	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			atr := atr.Values(ctx)

			if len(atr) == 0 {
				return fmt.Errorf("not enough data for ATR calculation")
			}

			currAtr := atr[len(atr)-1]
			pipDistance := currAtr * multiplier
			entryPrice := ctx.EntryPrice()

			switch order.Direction {
			case brokers.PositionDirectionLong:
				order.StopLoss = entryPrice - pipDistance
				return nil

			case brokers.PositionDirectionShort:
				order.StopLoss = entryPrice + pipDistance
				return nil

			default:
				panic("invalid position type")
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("StopLossATR",
				atr.Format(),
				formatter.Format(fmt.Sprintf("Multiplier: %.4f", multiplier)),
			)
		},
		func() (string, any) {
			return "stopLossATR", map[string]any{
				"atr":        marshal.ToJSON(atr),
				"multiplier": multiplier,
			}
		},
	)
}

func init() {
	jsonParsers.RegisterParser("stopLossATR", func(arg json.RawMessage) (OrderComputer, error) {
		var params struct {
			ATR        indicators.Indicator `json:"atr"`
			Multiplier float64              `json:"multiplier"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse StopLossATR parameters: %w", err)
		}

		return StopLossATR(params.ATR, params.Multiplier), nil
	})
}

const pipSize = 0.0001

func StopLossPipBuffer(pipBuffer int, lookupPeriod int) OrderComputer {
	pipDistance := float64(pipBuffer) * pipSize

	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {

			switch order.Direction {
			case brokers.PositionDirectionLong:
				// find lowest low in last lookupPeriod minutes
				lowest := ctx.HistoricalData().GetLowest(lookupPeriod)
				order.StopLoss = lowest - pipDistance
				return nil

			case brokers.PositionDirectionShort:
				// find highest high in last lookupPeriod minutes
				highest := ctx.HistoricalData().GetHighest(lookupPeriod)
				order.StopLoss = highest + pipDistance
				return nil

			default:
				return fmt.Errorf("invalid position direction: %s", order.Direction.String())
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("StopLossPipBuffer",
				formatter.Format(fmt.Sprintf("Pip Buffer: %d", pipBuffer)),
				formatter.Format(fmt.Sprintf("Lookup Period: %d", lookupPeriod)),
			)
		},
		func() (string, any) {
			return "stopLossPipBuffer", map[string]any{
				"pipBuffer":    pipBuffer,
				"lookupPeriod": lookupPeriod,
			}
		},
	)
}

func init() {
	jsonParsers.RegisterParser("stopLossPipBuffer", func(arg json.RawMessage) (OrderComputer, error) {
		var params struct {
			PipBuffer    int `json:"pipBuffer"`
			LookupPeriod int `json:"lookupPeriod"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse StopLossPipBuffer parameters: %w", err)
		}

		return StopLossPipBuffer(params.PipBuffer, params.LookupPeriod), nil
	})
}
