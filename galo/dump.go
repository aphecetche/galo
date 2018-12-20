package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aphecetche/galo"
	"github.com/spf13/cobra"
)

func LoopOverFile(fileName string, selName string, tcFunc func(index int, tc *galo.TaggedClusters, selected []int)) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dec := NewClusterDecoder(f, fileName)
	defer dec.Close()
	clusel := galo.NewClusterSelector(selName)
	if clusel == nil {
		log.Fatal("Do no know how to create cluster selector ", selName)
	}
	nevents, nsel := galo.ClusterLoop(dec, clusel, firstEvent, maxEvents, tcFunc)
	fmt.Println(nevents, " events processed. ", nsel, " matched")
}

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump clusters",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		LoopOverFile(args[0], cluSelArg, galo.DumpClusters)
	},
}

var outputFileName string

func init() {
	clusterCmd.AddCommand(dumpCmd)
}
