package plot

import (
	"image/color"

	"github.com/aphecetche/galo"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotter"
)

func landauFunc(mu, sigma float64, r, g, b uint8) *plotter.Function {
	f := plotter.NewFunction(func(x float64) float64 {
		return galo.Landau(x, mu, sigma) / sigma
	})
	f.Color = color.RGBA{R: r, B: b, G: g, A: 255}
	f.Samples = 1000
	return f
}

func levyFunc(mu, sigma float64, r, g, b uint8) *plotter.Function {
	f := plotter.NewFunction(func(x float64) float64 {
		return galo.Levy(x, mu, sigma) / sigma
	})
	f.Color = color.RGBA{R: r, B: b, G: g, A: 255}
	f.Samples = 1000
	return f
}

func SaveFunction(outputFileName string) {
	p := hplot.New()
	p.Add(landauFunc(0.25, 0.2, 0, 0, 255))
	p.X.Max = 10
	p.X.Min = -10
	p.Y.Max = 10
	p.Y.Min = -10
	SavePlot(p, outputFileName, "landau")
}
