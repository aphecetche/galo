// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aphecetche/galo/dataformats/run2"
	"github.com/spf13/cobra"
	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// plotCmd represents the plot command
var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot clusters",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		plotRun2Clusters(f)
	},
}

// ClusterPosFunc computes the (x,y) position of a cluster
type ClusterPosFunc = func(*run2.Cluster) (float64, float64)

func plotRun2Clusters(r io.Reader) {

	// hc := createHistogramCollection()

	run2.ForEachEvent(r, func(e *run2.EventClusters) {

		names := []string{"ww", "wow"}
		clusterPosFuncs := []ClusterPosFunc{cogWithWeight, cogNoWeight}
		for i, f := range clusterPosFuncs {
			hname := "/all/" + names[i]
			_ = f
			fmt.Println(hname)
			// h := hc.H1D(hname)
			// if h == nil {
			// 	log.Fatalf("could not get histogram %s\n", hname)
			// }
			// fillHisto(h, getClusterResidual(e, f))
		}
	}, maxEvents)

	// plotHistos(hc)
}

func getClusterResidual(ec *run2.EventClusters, clufunc ClusterPosFunc) []float64 {
	var res []float64
	var clu run2.Cluster
	for i := 0; i < ec.E.ClustersLength(); i++ {
		ec.E.Clusters(&clu, i)
		x, y := clufunc(&clu)
		pos := clu.Pos(nil)
		dx := x - float64(pos.X())
		dy := y - float64(pos.Y())
		d := math.Sqrt(dx*dx + dy*dy)
		res = append(res, d)
	}
	return res
}

func cogNoWeight(clu *run2.Cluster) (float64, float64) {
	return cog(clu.Pre(nil), false)
}

func cogWithWeight(clu *run2.Cluster) (float64, float64) {
	return cog(clu.Pre(nil), true)
}

// cog compute the center of gravity of the digits within precluster
func cog(pre *run2.PreCluster, useWeight bool) (float64, float64) {
	var digit run2.Digit
	var x, y, sumw float64
	for i := 0; i < pre.DigitsLength(); i++ {
		pre.Digits(&digit, i)
		deid := digit.Deid()
		manuid := int(digit.Manuid())
		seg := segmentation(int(deid), manuid < 1024)
		manuchannel := int(digit.Manuchannel())
		paduid, err := seg.FindPadByFEE(manuid, manuchannel)
		if seg.IsValid(paduid) == false || err != nil {
			log.Fatalf("got invalid pad for DE %v MANU %v CH %v : %v -> paduid %v", deid, manuid, manuchannel, err, paduid)
		}
		var w float64
		if useWeight == false {
			w = 1.0
		} else {
			w = float64(digit.Adc())
		}
		sumw += w

		x += seg.PadPositionX(paduid) * w
		y += seg.PadPositionY(paduid) * w
	}
	x /= sumw
	y /= sumw
	return x, y
}

func fillHisto(h *hbook.H1D, values []float64) {
	for _, v := range values {
		h.Fill(v, 1.0)
	}
}

func plotHisto(h *hbook.H1D, index int) {

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
	savePlot(p, index)
}

func savePlot(p *hplot.Plot, index int) {
	fname := strings.TrimSuffix(outputFileName, filepath.Ext(outputFileName))
	fname += strconv.Itoa(index) + ".pdf"
	err := p.Save(20*vg.Centimeter, -1, fname)
	if err != nil {
		log.Fatalf("Cannot save histogram:%s", err)
	}
}

// func createResidualHisto(hc *hist.Collection, path string, name string) {
// 	h := hbook.NewH1D(128, 0, 1)
// 	h.Ann["name"] = name
// 	hc.Add(path, h)
// }
//
// func createHistogramCollection() *hist.Collection {
//
// 	hc := hist.NewCollection("plot")
//
// 	createResidualHisto(hc, "/all/", "ww")
// 	createResidualHisto(hc, "/all/", "wow")
// 	createResidualHisto(hc, "/nodup/", "ww")
// 	createResidualHisto(hc, "/nodup/", "wow")
//
// 	hc.Print(os.Stdout)
//
// 	return hc
// }

// func plotHistos(hc *hist.Collection) {
//
// 	// for i, h := range histos {
// 	// 	if h == nil {
// 	// 		continue
// 	// 	}
// 	// 	fmt.Printf("%40s entries %4d Xmean %7.2f\n", h.Name(), h.Entries(), h.XMean())
// 	// 	plotHisto(h, i)
// 	// }
// }

func init() {
	clusterCmd.AddCommand(plotCmd)
	plotCmd.Flags().StringVarP(&outputFileName, "output", "o", "clusters.png", "Output image filename")
	segmentations = make(map[int]SegPair)
}
