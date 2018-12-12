package main

import (
	"fmt"
	"os"
	"strings"

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

		nevents := 0
		for {
			var declu galo.DEClusters
			err := dec.Decode(&declu)
			if len(declu.Clusters) > 0 {
				printHeader(nevents)
				printDEClusters(declu)
			}
			if err != nil {
				break
			}
			nevents++
			if nevents > maxEvents {
				break
			}
		}
		fmt.Println(nevents, " events processed")
	},
}

func printHeader(nevents int) {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Event %6d\n", nevents)
}

func printDEClusters(declu galo.DEClusters) {
	seg := galo.SegCache.Segmentation(declu.DeID)
	for _, c := range declu.Clusters {
		fmt.Printf("Q %7.2f POS %v\n", c.Q, c.Pos)
		for _, d := range c.Pre.Digits {
			fmt.Printf("%s%v\n", strings.Repeat(" ", 10), d.HumanReadable(seg))
		}
	}
}

var outputFileName string

func init() {
	clusterCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringVarP(&outputFileName, "output", "o", "clusters.png", "Output filename")
}
