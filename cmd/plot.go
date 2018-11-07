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
	"io"
	"log"
	"os"

	"github.com/aphecetche/galo/dataformats/run2"
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
		b := e.Clusters(&clu, i)
		if b == false {
			log.Fatalf("could not get cluster %d", i)
		}

		pos := clu.Pos(nil)

		clusterPositions = append(clusterPositions, struct{ X, Y float64 }{X: float64(pos.X()), Y: float64(pos.Y())})
		//	pre := clu.Pre(nil)
	}
	return clusterPositions
}

func plotRun2Clusters(r io.Reader) {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %s", err)
	}
	p.Title.Text = "Cluster positions"
	var cpos plotter.XYs

	run2.ForEachEvent(r, func(e *run2.Event) {
		cpos = append(cpos, getClusterPositions(e)...)
	}, maxEvents)

	s, err := plotter.NewScatter(cpos)
	if err != nil {
		log.Fatalf("could not create scatter plot: %s", err)
	}

	p.Add(s)

	// Save the plot to a file
	if err := p.Save(15*vg.Centimeter, 15*vg.Centimeter, outputFileName); err != nil {
		panic(err)
	}
}

func init() {
	clusterCmd.AddCommand(plotCmd)
	plotCmd.Flags().StringVarP(&outputFileName, "output", "o", "clusters.png", "Output image filename")
}
