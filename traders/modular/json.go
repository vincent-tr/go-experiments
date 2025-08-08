package modular

import (
	"encoding/json"
	"fmt"
	"go-experiments/traders/modular/conditions"
	"go-experiments/traders/modular/marshal"
	"go-experiments/traders/modular/ordercomputer"
)

type builderJSON struct {
	HistorySize      int             `json:"historySize"`
	Filter           json.RawMessage `json:"filter"`
	LongTrigger      json.RawMessage `json:"longTrigger"`
	ShortTrigger     json.RawMessage `json:"shortTrigger"`
	StopLoss         json.RawMessage `json:"stopLoss"`
	TakeProfit       json.RawMessage `json:"takeProfit"`
	CapitalAllocator json.RawMessage `json:"capitalAllocator"`
}

func FromJSON(jsonData []byte) (Builder, error) {
	var bjson builderJSON
	err := json.Unmarshal(jsonData, &bjson)
	if err != nil {
		return nil, err
	}

	res := &builder{
		historySize: bjson.HistorySize,
	}

	res.filter, err = conditions.FromJSON(bjson.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filter condition: %w", err)
	}

	res.longTrigger, err = conditions.FromJSON(bjson.LongTrigger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse long trigger condition: %w", err)
	}

	res.shortTrigger, err = conditions.FromJSON(bjson.ShortTrigger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse short trigger condition: %w", err)
	}

	res.stopLoss, err = ordercomputer.FromJSON(bjson.StopLoss)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stop loss order computer: %w", err)
	}

	res.takeProfit, err = ordercomputer.FromJSON(bjson.TakeProfit)
	if err != nil {
		return nil, fmt.Errorf("failed to parse take profit order computer: %w", err)
	}

	res.capitalAllocator, err = ordercomputer.FromJSON(bjson.CapitalAllocator)
	if err != nil {
		return nil, fmt.Errorf("failed to parse capital allocator order computer: %w", err)
	}

	return res, nil
}

func ToJSON(b Builder) ([]byte, error) {
	bu, ok := b.(*builder)
	if !ok {
		return nil, fmt.Errorf("invalid builder type")
	}

	bjson := &builderJSON{
		HistorySize:      bu.historySize,
		Filter:           marshal.ToJSON(bu.filter),
		LongTrigger:      marshal.ToJSON(bu.longTrigger),
		ShortTrigger:     marshal.ToJSON(bu.shortTrigger),
		StopLoss:         marshal.ToJSON(bu.stopLoss),
		TakeProfit:       marshal.ToJSON(bu.takeProfit),
		CapitalAllocator: marshal.ToJSON(bu.capitalAllocator),
	}

	return json.Marshal(bjson)
}
