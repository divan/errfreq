package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os/exec"
	"runtime"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func (a *Analyzer) PlotHistogram(file string) error {
	// prepare data
	data := make(plotter.Values, 0, len(a.Counts))
	for _, v := range a.Counts {
		data = append(data, float64(v))
	}

	sort.Float64s(data)

	// calculate bins (it's small number so we're good with this approach)
	min, max := math.MaxInt64, 0
	for _, value := range data {
		v := int(value)
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	bins := max - min

	// plot
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "'if err =' count per function"
	h, err := plotter.NewHist(data, bins)
	if err != nil {
		return fmt.Errorf("create histogram: %v", err)
	}
	h.Color = color.RGBA{255, 255, 255, 255}
	h.FillColor = plotutil.Color(2)

	h.Normalize(100)

	p.Add(h)

	stylePlot(p)

	err = p.Save(1920, 1280, file)
	if err != nil {
		return fmt.Errorf("save plot: %v", err)
	}

	fmt.Println("Plot stored to", file)
	return nil
}

// stylePlot applies styling to the plot
func stylePlot(p *plot.Plot) {
	p.Title.Font.SetName("Helvetica")
	p.Title.Font.Size = vg.Points(42)
	p.Title.Padding = vg.Points(100)

	p.X.Label.Text = "number of err checks"
	p.X.Padding = vg.Points(10)
	p.X.Tick.Marker = plot.TickerFunc(HistTick)
	p.X.Label.Font.SetName("Helvetica")
	p.X.Label.Font.Size = vg.Points(32)
	p.X.Tick.Label.Font.SetName("Helvetica")
	p.X.Tick.Label.Font.Size = vg.Points(24)

	p.Y.Label.Text = "Percent (%)"
	p.Y.Tick.Marker = plot.DefaultTicks{}
	p.Y.Label.Font.SetName("Helvetica")
	p.Y.Label.Font.Size = vg.Points(32)
	p.Y.Tick.Label.Font.SetName("Helvetica")
	p.Y.Tick.Label.Font.Size = vg.Points(24)
}

// OpenPlot tries to open the image with a system
// default app.
func OpenPlot(file string) {
	fmt.Println("Opening", file)
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], file)...)
	fmt.Println("If file wasn't opened by your OS, please open it manually in the viewer:", file)
	cmd.Start()
}

func HistTick(min, max float64) []plot.Tick {
	d := int(max - min)
	ticks := make([]plot.Tick, d)
	for i := 0; i < d; i++ {
		value := int(min) + i
		tick := plot.Tick{
			Value: float64(value),
			Label: fmt.Sprintf("%d", value),
		}
		ticks[i] = tick
	}
	return ticks
}
