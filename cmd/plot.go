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
	"image/color"
	"io"
	"log"
	"math"
	"os"

	"github.com/aphecetche/galo/dataformats/run2"
	"github.com/aphecetche/pigiron/mapping"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
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

func getClusterPositions(e *run2.Event) plotter.XYs {
	var clusterPositions plotter.XYs
	var clu run2.Cluster
	for i := 0; i < e.ClustersLength(); i++ {
		e.Clusters(&clu, i)
		pos := clu.Pos(nil)
		clusterPositions = append(clusterPositions, struct{ X, Y float64 }{X: float64(pos.X()), Y: float64(pos.Y())})
	}
	return clusterPositions
}

func getClusterResidual(e *run2.Event, prefunc func(*run2.PreCluster) (float64, float64)) []float64 {
	var res []float64
	var clu run2.Cluster
	for i := 0; i < e.ClustersLength(); i++ {
		e.Clusters(&clu, i)
		pos := clu.Pos(nil)
		pre := clu.Pre(nil)
		x, y := prefunc(pre)
		dx := x - float64(pos.X())
		dy := y - float64(pos.Y())
		d := math.Sqrt(dx*dx + dy*dy)
		res = append(res, d)
	}
	return res
}

type SegPair struct {
	Bending, NonBending mapping.Segmentation
}

var segmentations map[int]SegPair

func segmentation(deid int, bending bool) mapping.Segmentation {
	seg := segmentations[deid]
	if seg.Bending == nil {
		segmentations[deid] = SegPair{
			Bending:    mapping.NewSegmentation(deid, true),
			NonBending: mapping.NewSegmentation(deid, false),
		}
		seg = segmentations[deid]
	}
	if bending {
		return seg.Bending
	}

	return seg.NonBending
}

func cogNoWeight(pre *run2.PreCluster) (float64, float64) {
	return cog(pre, false)
}

func cogWithWeight(pre *run2.PreCluster) (float64, float64) {
	return cog(pre, true)
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

func plotRun2Clusters(r io.Reader) {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %s", err)
	}
	p.Title.Text = "Cluster positions"
	// var cpos plotter.XYs
	var res plotter.Values
	var resW plotter.Values

	run2.ForEachEvent(r, func(e *run2.Event) {
		// cpos = append(cpos, getClusterPositions(e)...)
		res = append(res, getClusterResidual(e, cogNoWeight)...)
		resW = append(resW, getClusterResidual(e, cogWithWeight)...)
	}, maxEvents)

	// s, err := plotter.NewScatter(cpos)
	// if err != nil {
	// 	log.Fatalf("could not create scatter plot: %s", err)
	// }
	// p.Add(s)

	h, err := plotter.NewHist(res, 512)
	if err != nil {
		log.Fatalf("could not create histogram: %s", err)
	}
	hw, err := plotter.NewHist(resW, 256)
	if err != nil {
		log.Fatalf("could not create histogram: %s", err)
	}

	blue := color.RGBA{B: 255, A: 255}
	red := color.RGBA{R: 255, A: 255}
	h.LineStyle.Color = blue
	hw.LineStyle.Color = red

	h.FillColor = blue
	hw.FillColor = red

	h.Normalize(float64(len(res)))
	hw.Normalize(float64(len(res)))
	p.Add(hw, h)

	p.X.Max = 1.2
	// p.Y.Min = 1E-1
	// p.Y.Scale = plot.LogScale{}

	// Save the plot to a file
	if err := p.Save(15*vg.Centimeter, 15*vg.Centimeter, outputFileName); err != nil {
		panic(err)
	}
}

func init() {
	clusterCmd.AddCommand(plotCmd)
	plotCmd.Flags().StringVarP(&outputFileName, "output", "o", "clusters.png", "Output image filename")
	segmentations = make(map[int]SegPair)
}
