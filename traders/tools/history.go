package tools

import (
	"go-experiments/brokers"
)

type History struct {
	candles []brokers.Candle
	maxSize int
}

func NewHistory(maxSize int) *History {
	return &History{
		candles: make([]brokers.Candle, 0),
		maxSize: maxSize,
	}
}

func (h *History) AddCandle(candle brokers.Candle) {
	if len(h.candles) >= h.maxSize {
		h.candles = h.candles[1:] // Remove the oldest candle
	}

	h.candles = append(h.candles, candle)
}

func (h *History) GetClosePrices() []float64 {
	size := len(h.candles)

	if size < h.maxSize {
		return nil // Not enough data
	}

	prices := make([]float64, size)

	for i := 0; i < size; i++ {
		prices[i] = h.candles[i].Close
	}

	return prices
}
