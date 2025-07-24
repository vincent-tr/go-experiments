package indicators

import (
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

type cache struct {
	indicators map[string][]float64
}

func NewCache() context.IndicatorCache {
	return &cache{
		indicators: make(map[string][]float64),
	}
}

func (c *cache) Tick() {
	c.indicators = make(map[string][]float64)
}

func (c *cache) access(key string, computer func() []float64) []float64 {
	if data, found := c.indicators[key]; found {
		return data
	}
	data := computer()
	c.indicators[key] = data
	return data
}

type Indicator interface {
	formatter.Formatter
	Values(ctx context.TraderContext) []float64
}

type indicator struct {
	compute func(ctx context.TraderContext) []float64
	format  func() *formatter.FormatterNode
}

func newIndicator(compute func(ctx context.TraderContext) []float64, format func() *formatter.FormatterNode) Indicator {
	return &indicator{
		compute: compute,
		format:  format,
	}
}

func (i *indicator) Values(ctx context.TraderContext) []float64 {
	c := ctx.IndicatorCache().(*cache)
	key := i.format().String()

	return c.access(key, func() []float64 {
		return i.compute(ctx)
	})
}

func (i *indicator) Format() *formatter.FormatterNode {
	return i.format()
}
