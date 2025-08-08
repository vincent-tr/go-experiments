package indicators

import (
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/marshal"
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
	ToJsonSpec() (string, any)
}

type indicator struct {
	compute    func(ctx context.TraderContext) []float64
	format     func() *formatter.FormatterNode
	toJsonSpec func() (string, any)
}

func newIndicator(
	compute func(ctx context.TraderContext) []float64,
	format func() *formatter.FormatterNode,
	toJsonSpec func() (string, any),
) Indicator {
	return &indicator{
		compute:    compute,
		format:     format,
		toJsonSpec: toJsonSpec,
	}
}

func (i *indicator) Values(ctx context.TraderContext) []float64 {
	c := ctx.IndicatorCache().(*cache)
	key := i.format().Compact()

	return c.access(key, func() []float64 {
		return i.compute(ctx)
	})
}

func (i *indicator) Format() *formatter.FormatterNode {
	return i.format()
}

func (i *indicator) ToJsonSpec() (string, any) {
	return i.toJsonSpec()
}

var jsonParsers = marshal.NewRegistry[Indicator]()

func FromJSON(jsonData []byte) (Indicator, error) {
	return jsonParsers.FromJSON(jsonData)
}
