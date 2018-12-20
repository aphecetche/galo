package galo

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strings"

	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func CreateHistogramCollection() *Collection {
	hc := NewCollection("clusters")
	h := hbook.NewH1D(500, 0, 500)
	h.Annotation()["name"] = "multiplicity"
	hc.Add(h)
	h = hbook.NewH1D(128, 0, 1)
	h.Annotation()["name"] = "residual"
	hc.Add(h)
	return hc
}

func FillHistogramCollection(tc *TaggedClusters, selected []int, hc *Collection, declupos DEClusterPositioner) {
	for _, i := range selected {
		hname := "multiplicity"
		h, err := hc.H1D("multiplicity")
		if err != nil {
			log.Fatalf("could not get histogram %s\n", hname)
		}
		clu := tc.Clusters()[i]
		h.Fill(float64(clu.Pre.NofPads()), 1.0)
		res := DEClusterResidual(tc.declu, i, declupos)
		hname = "residual"
		h, err = hc.H1D(hname)
		if err != nil {
			log.Fatalf("could not get histogram %s\n", hname)
		}
		h.Fill(res, 1.0)
	}
}

type ClipScale struct {
	Min  float64
	Max  float64
	Norm plot.Normalizer
}

func (cs ClipScale) Normalize(min, max, x float64) float64 {
	min = math.Max(cs.Min, min)
	max = math.Min(cs.Max, max)
	switch {
	case x < cs.Min:
		x = cs.Min
	case x > cs.Max:
		x = cs.Max
	}
	return cs.Norm.Normalize(min, max, x)
}

var _ plot.Normalizer = ClipScale{}

type ClipTicker struct {
	Min    float64
	Max    float64
	Ticker plot.Ticker
}

var _ plot.Ticker = ClipTicker{}

func (ct ClipTicker) Ticks(min, max float64) []plot.Tick {
	min = math.Max(min, ct.Min)
	max = math.Min(max, ct.Max)
	return ct.Ticker.Ticks(min, max)
}

func plotMultiplicity(h *hbook.H1D) *hplot.Plot {
	if h.Entries() == 0 {
		return nil
	}
	fmt.Println(h.Annotation()["name"])
	p := hplot.New()
	p.Y.Min = 0.5
	p.Y.Scale = ClipScale{p.Y.Min, math.Inf(+1), plot.LogScale{}}
	p.Y.Tick.Marker = ClipTicker{p.Y.Min, math.Inf(+1), plot.LogTicks{}}
	p.Add(hplot.NewH1D(h))
	return p
}

func plotResidual(h *hbook.H1D) *hplot.Plot {

	p := hplot.New()
	p.X.Label.Text = "Distance (cm)"
	h.Scale(1 / h.Integral())

	hh := hplot.NewH1D(h)
	p.Add(hh)

	res, err := fit.H1D(
		h,
		fit.Func1D{
			F: func(x float64, params []float64) float64 {
				return Moyal(x, params[0], params[1], params[2])
			},
			N: 3,
		},
		nil,
		&optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}

	f := plotter.NewFunction(func(x float64) float64 {
		return Moyal(x, res.X[0], res.X[1], res.X[2])
	})
	if res.X[1] > 0 {
		h.Ann["mu"] = res.X[1]
		h.Ann["sigma"] = res.X[2]
	}
	f.Color = color.RGBA{R: 255, A: 255}
	f.Samples = 1000
	p.Add(f)

	p.X.Max = 1.0
	return p
}

func PlotHistogramCollection(hc *Collection) []*hplot.Plot {
	var plots []*hplot.Plot
	for _, h := range hc.H1Ds() {
		if h == nil || h.Entries() == 0 {
			continue
		}
		// fmt.Printf("%40s entries %4d Xmean %7.2f\n", h.Name(), h.Entries(), h.XMean())
		name := (h.Annotation()["name"]).(string)
		if strings.Contains(name, "residual") {
			plots = append(plots, plotResidual(h))
		}
		if strings.Contains(name, "multiplicity") {
			plots = append(plots, plotMultiplicity(h))
		}
	}
	hc.Print(os.Stdout)
	return plots
}
