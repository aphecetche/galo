package run2

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strings"

	"github.com/aphecetche/galo"
	galoplot "github.com/aphecetche/galo/plot"
	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// var (
// 	ClusterPosFuncs = []ClusterPosFunc{
// 		// {cogWithSquaredWeight, "w2"},
// 		{cogNoWeight, "cog"},
// 		// {cogWithRegularWeight, "cogW"},
// 	}
// 	ClusterSelFuncs = []ClusterSelFunc{
// 		{allClusters, "allClusters"},
// 		{simpleClusters, "simpleClusters"},
// 		{splitClusters, "splitClusters"},
// 		{dupClusters, "dupClusters"},
// 		{strangeClusters, "strangeClusters"},
// 		{chargeClusters, "chargeClusters"},
// 	}
// )
//

var (
	cluSelectors   []galo.ClusterSelector
	cluPositioners []galo.ClusterPositioner
)

func init() {
	cluPositioners = append(cluPositioners,
		galo.ClusterPositionerCOG{Wmod: galo.NoWeight})
}

func createHistogramCollection() *galo.Collection {

	hc := galo.NewCollection("plot")

	for _, cluselfunc := range cluSelectors {
		hname := "/" + cluselfunc.Name() + "/multiplicity"
		h := hbook.NewH1D(500, 0, 500)
		h.Annotation()["name"] = hname
		hc.Add(h)
		for _, cluposfunc := range cluPositioners {
			hname := "/" + cluselfunc.Name() + "/residual_" + cluposfunc.Name()
			createResidualHisto(hc, hname)
		}
	}

	return hc
}

func fillHistogramCollection(declu *galo.DEClusters, hc *galo.Collection, cc *galo.CounterCollection) {

	//here hc should be hl <=> galo.Library <=> map of path -> galo.Collection

	for _, clu := range declu.Clusters {

		for _, cluselfunc := range cluSelectors {
			if cluselfunc.Select(clu) == false {
				continue
			}
			(*cc).Incr(cluselfunc.Name())

			hname := "/" + cluselfunc.Name() + "/multiplicity"
			h, err := hc.H1D(hname)
			if err != nil {
				log.Fatalf("could not get histogram %s\n", hname)
			}
			h.Fill(float64(clu.Pre.NofPads()), 1.0)

			for _, cluposfunc := range cluPositioners {
				res := galo.ClusterResidual(clu, cluposfunc)
				// here should instead get a "collection" with all the
				// histos/numbers/drawings I'd want to fill
				// for this cluster :
				//
				// - H1 residuals
				// - H2 positions
				// - SVG <=> canvas ?
				//
				hname := "/" + cluselfunc.Name() + "/residual_" + cluposfunc.Name()
				h, err := hc.H1D(hname)
				if err != nil {
					log.Fatalf("could not get histogram %s\n", hname)
				}
				h.Fill(res, 1.0)
			}
		}
	}
}

func PlotClusters(dec galo.DEClustersDecoder, maxEvents int, outputFileName string) {

	hc := createHistogramCollection()

	cc := galo.NewCounterCollection()

	var declu galo.DEClusters
	for {
		err := dec.Decode(&declu)
		if err != nil {
			break
		}
		cc.Incr("events")
		fillHistogramCollection(&declu, hc, cc)
	}
	plotHistogramCollection(hc, outputFileName)
	galoplot.SaveFunction(outputFileName)
	fmt.Println(cc)
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

func plotMultiplicity(h *hbook.H1D, outputFileName string) {

	if h.Entries() == 0 {
		return
	}

	fmt.Println(h.Annotation()["name"])
	p := hplot.New()
	p.Y.Min = 0.5
	p.Y.Scale = ClipScale{p.Y.Min, math.Inf(+1), plot.LogScale{}}
	p.Y.Tick.Marker = ClipTicker{p.Y.Min, math.Inf(+1), plot.LogTicks{}}
	p.Add(hplot.NewH1D(h))
}

func plotResidual(h *hbook.H1D, outputFileName string) {

	p := hplot.New()
	p.X.Label.Text = "Distance (cm)"
	h.Scale(1 / h.Integral())

	hh := hplot.NewH1D(h)
	p.Add(hh)

	res, err := fit.H1D(
		h,
		fit.Func1D{
			F: func(x float64, params []float64) float64 {
				// return f1d.Gaus(x, params[0], params[1], params[2])
				return galo.Moyal(x, params[0], params[1], params[2])
				// return params[0] * f1d.Landau(x, params[1], params[2])
			},
			N: 3,
			// Ps: []float64{1.0, 0.1},
			// N: 3,
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
		// return gaus(x, res.X[0], res.X[1], res.X[2])
		return galo.Moyal(x, res.X[0], res.X[1], res.X[2])
	})
	if res.X[1] > 0 {
		h.Ann["mu"] = res.X[1]
		h.Ann["sigma"] = res.X[2]
	}
	f.Color = color.RGBA{R: 255, A: 255}
	f.Samples = 1000
	p.Add(f)

	p.X.Max = 1.0
	galoplot.SavePlot(p, outputFileName, h.Name())
}

func createResidualHisto(hc *galo.Collection, name string) {
	h := hbook.NewH1D(128, 0, 1)
	h.Annotation()["name"] = name
	hc.Add(h)
}

func plotHistogramCollection(hc *galo.Collection, outputFileName string) {
	for _, h := range hc.H1Ds() {
		if h == nil || h.Entries() == 0 {
			continue
		}
		// fmt.Printf("%40s entries %4d Xmean %7.2f\n", h.Name(), h.Entries(), h.XMean())
		name := (h.Annotation()["name"]).(string)
		if strings.Contains(name, "residual") {
			plotResidual(h, outputFileName)
		}
		if strings.Contains(name, "multiplicity") {
			plotMultiplicity(h, outputFileName)
		}
	}
	hc.Print(os.Stdout)
}
