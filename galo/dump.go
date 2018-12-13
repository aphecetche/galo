package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aphecetche/galo"
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
		dec := NewClusterDecoder(f, args[0])
		defer dec.Close()
		clusel := galo.NewClusterSelector(cluSelArg)
		if clusel == nil {
			log.Fatal("Do no know how to create cluster selector ", cluSelArg)
		}
		nevents, nsel := galo.ClusterLoop(dec, clusel, firstEvent, maxEvents, galo.DumpClusters)
		fmt.Println(nevents, " events processed. ", nsel, " matched")
	},
}

var outputFileName string

func init() {
	clusterCmd.AddCommand(dumpCmd)
}
