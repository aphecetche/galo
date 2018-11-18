package run2

import (
	"image/color"

	"github.com/aphecetche/galo/f1d"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotter"
)

func landauFunc(mu, sigma float64, r, g, b uint8) *plotter.Function {
	f := plotter.NewFunction(func(x float64) float64 {
		return f1d.Landau(x, mu, sigma) / sigma
	})
	f.Color = color.RGBA{R: r, B: b, G: g, A: 255}
	f.Samples = 1000
	return f
}

func levyFunc(mu, sigma float64, r, g, b uint8) *plotter.Function {
	f := plotter.NewFunction(func(x float64) float64 {
		return f1d.Levy(x, mu, sigma) / sigma
	})
	f.Color = color.RGBA{R: r, B: b, G: g, A: 255}
	f.Samples = 1000
	return f
}

func mathiesonFunc(k3 float64, r, g, b uint8) *plotter.Function {
	f := plotter.NewFunction(func(x float64) float64 {
		return f1d.Mathieson(x, k3)
	})
	f.Color = color.RGBA{R: r, B: b, G: g, A: 255}
	f.Samples = 1000
	return f
}
func saveFunction(outputFileName string) {
	p := hplot.New()
	// p.Add(landauFunc(0.25, 0.2, 0, 0, 255))
	// p.Add(landauFunc(0.25, 0.1, 255, 0, 0))
	// p.Add(landauFunc(0.2, 0.05, 0, 255, 0))
	// p.Add(landauFunc(0.05, 0.02, 55, 55, 55))
	// p.Add(levyFunc(1.0, 1.5, 255, 0, 255))
	p.Add(mathiesonFunc(1.0, 255, 0, 255))
	p.X.Max = 10
	p.X.Min = -10
	p.Y.Max = 10
	p.Y.Min = -10
	savePlot(p, outputFileName, "landau")
}
