package mathieson

import (
	"fmt"
	"image/color"
	"math"

	"github.com/aphecetche/galo"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type MakePlotFunc func(fname string)

func (f MakePlotFunc) MakePlots(fname string) {
	f(fname)
}

var Plotter = MakePlotFunc(MakePlots)

func MakePlots(fname string) {
	plot1D(fname)
}

func plotM1(p *hplot.Plot, dir byte, dist Mathieson1D, r, g, b uint8, dashes []vg.Length) {
	f := plotter.NewFunction(func(x float64) float64 {
		return dist.F(x)
	})
	f.Color = color.RGBA{R: r, B: b, G: g, A: 255}
	f.Samples = 1000
	f.Dashes = dashes
	p.Add(f)
	p.Legend.Add(fmt.Sprintf("K3%c=%7.3f Pitch=%5.2f cm", dir, dist.K3(), dist.Pitch()), f)
}

func gausFunc(mu, sigma float64, r, g, b uint8) *plotter.Function {
	f := plotter.NewFunction(func(x float64) float64 {
		//c := math.Sqrt(2.0*math.Pi) * sigma
		c := 1.0
		return galo.Gaus(x, 1.0/c, mu, sigma)
	})
	f.Color = color.RGBA{R: r, B: b, G: g, A: 255}
	f.Samples = 1000
	return f
}

func plot1D(fname string) {
	p := hplot.New()

	plotM1(p, 'x', *NewMathieson1D(St1.Pitch, St1.K3x), 255, 0, 0, nil)
	plotM1(p, 'y', *NewMathieson1D(St1.Pitch, St1.K3y), 255, 0, 0, []vg.Length{vg.Points(2), vg.Points(2)})
	plotM1(p, 'x', *NewMathieson1D(St2345.Pitch, St2345.K3x), 0, 0, 255, nil)
	plotM1(p, 'y', *NewMathieson1D(St2345.Pitch, St2345.K3y), 0, 0, 255, []vg.Length{vg.Points(2), vg.Points(2)})

	d := *NewMathieson1D(St1.Pitch, St1.K3x)

	m7 := NewMathieson1D(0.25, 0.7)
	m1 := NewMathieson1D(0.25, 1.0)

	sigma := d.FWHM() / (2 * math.Sqrt(2*math.Log(2)))

	fmt.Printf("d=%7.2f sigma=%7.2f\n", d.FWHM(), sigma)
	fmt.Printf("m1=%7.2f m7=%7.2f\n", m1.FWHM(), m7.FWHM())

	g := gausFunc(0, sigma, 128, 128, 128)
	p.Add(g)

	p.Legend.Add(fmt.Sprintf("Gaus sigma=%7.2f cm", sigma), g)

	p.X.Min = -3
	p.X.Max = 3
	p.Y.Min = 0
	p.Y.Max = 1

	p.X.Label.Text = "Lambda"
	p.Y.Min = 0.01
	// p.Y.Scale = &plot.LogScale{}

	font, err := vg.MakeFont("Helvetica", 12)
	if err != nil {
		panic(err)
	}

	p.X.Label.Font = font
	p.Legend.Font = font
	// p.X.Label.Text = "λ"

	p.Legend.Top = true
	galo.SavePlot(p, fname, "1d")
}
