package main

import (
	"fmt"
	"os"

	"github.com/aphecetche/galo/run2"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump clusters",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		nevents := run2.ForEachEvent(f, run2.DumpEventClusters, maxEvents)
		fmt.Println(nevents, " events processed")
	},
}

var outputFileName string

func init() {
	clusterCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringVarP(&outputFileName, "output", "o", "clusters.png", "Output filename")
}
