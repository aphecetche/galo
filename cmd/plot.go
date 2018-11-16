package cmd

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
		run2.PlotClusters(f, maxEvents, outputFileName)
	},
}

func init() {
	clusterCmd.AddCommand(plotCmd)
	plotCmd.Flags().StringVarP(&outputFileName, "output", "o", "clusters.png", "Output image filename")
}
