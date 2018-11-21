package mathieson

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/aphecetche/galo/f1d"
	"github.com/aphecetche/galo/util"
	"github.com/aphecetche/pigiron/mapping"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func MakePlots(fname string) {
	plot1D(fname)

	seg := mapping.NewSegmentation(100, true)
	paduid, err := seg.FindPadByPosition(24, 24)
	if err != nil {
		panic("could not get pad")
	}

	mapping.PrintPad(os.Stdout, seg, paduid)
}

func plotM1(p *hplot.Plot, dir byte, dist Dist1D, r, g, b uint8, dashes []vg.Length) {
	f := plotter.NewFunction(func(x float64) float64 {
		return dist.F(x)
	})
	f.Color = color.RGBA{R: r, B: b, G: g, A: 255}
	f.Samples = 1000
	f.Dashes = dashes
	p.Add(f)
	p.Legend.Add(fmt.Sprintf("K3%c=%7.3f Pitch=%5.2f cm", dir, dist.K3(), dist.Pitch()), f)
}

func plotM2(p *hplot.Plot, dist Dist2D, r, g, b uint8) {
	plotM1(p, 'x', dist.X, r, g, b, nil)
	plotM1(p, 'y', dist.Y, r, g, b, []vg.Length{vg.Points(2), vg.Points(2)})
}

func gausFunc(mu, sigma float64, r, g, b uint8) *plotter.Function {
	f := plotter.NewFunction(func(x float64) float64 {
		//c := math.Sqrt(2.0*math.Pi) * sigma
		c := 1.0
		return f1d.Gaus(x, 1.0/c, mu, sigma)
	})
	f.Color = color.RGBA{R: r, B: b, G: g, A: 255}
	f.Samples = 1000
	return f
}

func plot1D(fname string) {
	p := hplot.New()

	plotM2(p, St1, 255, 0, 0)
	plotM2(p, St2345, 0, 0, 255)

	d := St1.X

	m7 := NewDist1D(0.25, 0.7)
	m1 := NewDist1D(0.25, 1.0)

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
	util.SavePlot(p, fname, "1d")
}