package main

import (
	"fmt"
	"go-experiments/brokers/backtesting"
	"go-experiments/common"
	"image/color"
	"log"
	"os"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type bucket struct {
	date  string
	count int
}

func main() {
	dataset, err := backtesting.LoadDataset(
		common.NewMonth(2024, 1),
		common.NewMonth(2024, 1),
		"EURUSD",
	)

	if err != nil {
		panic(err)
	}

	buckets := make([]*bucket, 0)
	bucketsMap := make(map[string]*bucket)

	// Fill buckets per hour and map for the whole time range
	begin := dataset.BeginDate()
	// to fill last day. only last bucket will be empty.
	// also it looks like the data contains ticks when using UTC that are the day after.
	end := dataset.EndDate().AddDate(0, 0, 2)

	for d := begin; d.Before(end); d = d.Add(time.Hour) {
		bucket := &bucket{
			date:  d.Format("2006-01-02 15"),
			count: 0,
		}
		buckets = append(buckets, bucket)
		bucketsMap[bucket.date] = bucket
	}

	// Map ticks to their respective buckets
	for tick := range dataset.Ticks() {
		tickDate := tick.GetTimestamp().Format("2006-01-02 15")
		bucket := bucketsMap[tickDate]
		bucket.count++
	}

	doplot(buckets)
}

// weekendBand implements plot.Plotter
type weekendBand struct {
	xmin, xmax float64
	ymax       float64
}

func (w weekendBand) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	x0 := trX(w.xmin)
	x1 := trX(w.xmax)
	y0 := trY(0)
	y1 := trY(w.ymax)

	c.FillPolygon(color.RGBA{R: 200, G: 200, B: 200, A: 80},
		[]vg.Point{{X: x0, Y: y0}, {X: x1, Y: y0}, {X: x1, Y: y1}, {X: x0, Y: y1}})
}

func (w weekendBand) DataRange() (xmin, xmax, ymin, ymax float64) {
	return w.xmin, w.xmax, 0, w.ymax
}

func doplot(data []*bucket) {
	p := plot.New()
	p.Title.Text = "Tick Density (per Hour)"
	p.Y.Label.Text = "Tick Count"

	// Convert buckets → plot values
	values := make(plotter.Values, len(data))
	labels := make([]string, len(data))
	maxY := 0.0

	for i, b := range data {
		values[i] = float64(b.count)

		t, _ := time.Parse("2006-01-02 15", b.date)
		if t.Hour() == 0 {
			// Show "01-02 (Mon)" once a day
			prefix := ""
			if t.Weekday() == time.Monday {
				prefix = "| "
			}
			labels[i] = prefix + t.Format("01-02 (Mon)")
		}

		// Add shaded band for weekend hours
		if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
			p.Add(weekendBand{
				xmin: float64(i) - 0.5,
				xmax: float64(i) + 0.5,
				ymax: maxY * 1.1, // slightly above max ticks
			})
		}
	}

	// Bars
	bars, err := plotter.NewBarChart(values, vg.Points(10))
	if err != nil {
		log.Fatal(err)
	}
	bars.Color = color.RGBA{R: 60, G: 120, B: 200, A: 255}
	bars.LineStyle.Width = 0
	p.Add(bars)

	p.NominalX(labels...)
	p.Add(plotter.NewGrid())

	// High-resolution canvas
	width := 2000
	height := 800
	canvas := vgimg.New(vg.Points(float64(width)), vg.Points(float64(height)))
	p.Draw(draw.New(canvas))

	// Save as PNG
	f, err := os.Create("ticks.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	png := vgimg.PngCanvas{Canvas: canvas}
	_, err = png.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Chart saved to ticks.png")
}
