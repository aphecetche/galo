package run2

import (
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/aphecetche/galo/hist"
	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var (
	ClusterPosFuncs = []ClusterPosFunc{
		// {cogWithSquaredWeight, "w2"},
		{cogNoWeight, "cog"},
		// {cogWithRegularWeight, "cogW"},
	}
	ClusterSelFuncs = []ClusterSelFunc{
		{allClusters, "allClusters"},
		{simpleClusters, "simpleClusters"},
		{splitClusters, "splitClusters"},
		{dupClusters, "dupClusters"},
	}
)

func createHistogramCollection() *hist.Collection {

	hc := hist.NewCollection("plot")

	for _, cluselfunc := range ClusterSelFuncs {
		for _, cluposfunc := range ClusterPosFuncs {
			hname := "/" + cluselfunc.Name + "/" + cluposfunc.Name
			createResidualHisto(hc, hname)
		}
	}

	return hc
}

func fillHistogramCollection(ec *EventClusters, hc *hist.Collection, cc *hist.CounterCollection) {

	//here hc should be hl <=> hist.Library <=> map of path -> hist.Collection

	for i := 0; i < ec.E.ClustersLength(); i++ {
		for _, cluselfunc := range ClusterSelFuncs {
			if cluselfunc.F(ec, i) == false {
				continue
			}
			(*cc).Incr(cluselfunc.Name)
			for _, cluposfunc := range ClusterPosFuncs {
				res := getClusterResidual(ec, i, cluposfunc)
				// here should instead get a "collection" with all the
				// histos/numbers/drawings I'd want to fill
				// for this cluster :
				//
				// - H1 residuals
				// - H2 positions
				// - SVG <=> canvas ?
				//
				hname := "/" + cluselfunc.Name + "/" + cluposfunc.Name
				h, err := hc.H1D(hname)
				if err != nil {
					log.Fatalf("could not get histogram %s\n", hname)
				}
				h.Fill(res, 1.0)
			}
		}
	}
}

func PlotClusters(r io.ReaderAt, maxEvents int, outputFileName string) {

	hc := createHistogramCollection()

	cc := hist.NewCounterCollection()

	ForEachEvent(r, func(ec *EventClusters) {
		cc.Incr("events")
		fillHistogramCollection(ec, hc, cc)
	}, maxEvents)

	plotHistogramCollection(hc, outputFileName)
	fmt.Println(cc)
}

func getClusterResidual(ec *EventClusters, i int, cluposfunc ClusterPosFunc) float64 {
	var clu Cluster
	ec.E.Clusters(&clu, i)
	x, y := cluposfunc.F(&clu)
	pos := clu.Pos(nil)
	dx := x - float64(pos.X())
	dy := y - float64(pos.Y())
	return math.Sqrt(dx*dx + dy*dy)
}

func plotHisto(h *hbook.H1D, outputFileName string) {

	p := hplot.New()
	p.X.Label.Text = "Distance (cm)"
	h.Scale(1 / h.Integral())

	hh := hplot.NewH1D(h)
	p.Add(hh)

	gaus := func(x, cst, mu, sigma float64) float64 {
		v := (x - mu) / sigma
		return cst * math.Exp(-0.5*v*v)
	}

	res, err := fit.H1D(
		h,
		fit.Func1D{
			F: func(x float64, params []float64) float64 {
				return gaus(x, params[0], params[1], params[2])
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

	fmt.Println("mu=", res.X[1])
	f := plotter.NewFunction(func(x float64) float64 {
		return gaus(x, res.X[0], res.X[1], res.X[2])
	})
	f.Color = color.RGBA{R: 255, A: 255}
	f.Samples = 1000
	p.Add(f)

	p.X.Max = 1.0
	savePlot(p, outputFileName, h.Name())
}

func savePlot(p *hplot.Plot, outputFileName string, name string) {
	fname := strings.TrimSuffix(outputFileName, filepath.Ext(outputFileName))

	fname += strings.Replace(name, "/", "_", -1) + ".pdf"
	err := p.Save(20*vg.Centimeter, -1, fname)
	if err != nil {
		log.Fatalf("Cannot save histogram:%s", err)
	}
}

func createResidualHisto(hc *hist.Collection, name string) {
	h := hbook.NewH1D(128, 0, 1)
	h.Annotation()["name"] = name
	hc.Add(h)
}

func plotHistogramCollection(hc *hist.Collection, outputFileName string) {
	hc.Print(os.Stdout)
	for _, h := range hc.H1Ds() {
		if h == nil || h.Entries() == 0 {
			continue
		}
		fmt.Printf("%40s entries %4d Xmean %7.2f\n", h.Name(), h.Entries(), h.XMean())
		plotHisto(h, outputFileName)
	}
}
