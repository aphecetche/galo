package main

import (
	"os"

	"github.com/aphecetche/galo/run2"
	"github.com/spf13/cobra"
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
		dec := NewClusterDecoder(f, args[0])
		defer dec.Close()
		run2.PlotClusters(dec, maxEvents, outputFileName)
	},
}
var silent bool

func init() {
	clusterCmd.AddCommand(plotCmd)
	plotCmd.Flags().BoolVarP(&silent, "silent", "s", false, "no text output")
}
