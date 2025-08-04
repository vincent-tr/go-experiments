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

func (h *History) IsUsable() bool {
	if len(h.candles) < h.maxSize {
		return false
	}

	for _, candle := range h.candles {
		if !candle.Usable {
			return false
		}
	}

	return true
}

func (h *History) AddCandle(candle brokers.Candle) {
	if len(h.candles) >= h.maxSize {
		h.candles = h.candles[1:] // Remove the oldest candle
	}

	h.candles = append(h.candles, candle)
}

func (h *History) GetPrice() float64 {
	return h.candles[len(h.candles)-1].Close
}

func (h *History) GetClosePrices() []float64 {
	size := len(h.candles)
	prices := make([]float64, size)

	for i := 0; i < size; i++ {
		prices[i] = h.candles[i].Close
	}

	return prices
}

func (h *History) GetHighPrices() []float64 {
	size := len(h.candles)
	prices := make([]float64, size)

	for i := 0; i < size; i++ {
		prices[i] = h.candles[i].High
	}

	return prices
}

func (h *History) GetLowPrices() []float64 {
	size := len(h.candles)
	prices := make([]float64, size)

	for i := 0; i < size; i++ {
		prices[i] = h.candles[i].Low
	}

	return prices
}

func (h *History) GetLowest(timeperiod int) float64 {
	startIndex := len(h.candles) - timeperiod

	lowest := h.candles[startIndex].Low
	for i := startIndex + 1; i < len(h.candles); i++ {
		if h.candles[i].Low < lowest {
			lowest = h.candles[i].Low
		}
	}

	return lowest
}

func (h *History) GetHighest(timeperiod int) float64 {
	startIndex := len(h.candles) - timeperiod

	highest := h.candles[startIndex].High
	for i := startIndex + 1; i < len(h.candles); i++ {
		if h.candles[i].High > highest {
			highest = h.candles[i].High
		}
	}

	return highest
}
