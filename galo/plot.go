package main

import (
	"strconv"

	"github.com/aphecetche/galo"
	"github.com/spf13/cobra"
)

// plotCmd represents the plot command
var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot clusters",
	Run: func(cmd *cobra.Command, args []string) {
		hc := galo.CreateHistogramCollection()
		positioner := galo.NewDEClusterPositioner(positionerArg)
		LoopOverFile(args[0], cluSelArg, func(index int, tc *galo.TaggedClusters, selected []int) {
			galo.FillHistogramCollection(tc, selected, hc, positioner)
		})
		plots := galo.PlotHistogramCollection(hc)
		for i, p := range plots {
			galo.SavePlot(p, outputFileName, "toto"+strconv.Itoa(i))
		}
	},
}

var quiet bool

var positionerArg string

func init() {
	clusterCmd.AddCommand(plotCmd)
	plotCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "no text output")
	plotCmd.Flags().StringVarP(&positionerArg, "positioner", "p", "COG", "Positioner to use")
}
