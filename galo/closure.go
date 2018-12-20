package main

import (
	"github.com/spf13/cobra"
)

// Closure generates n charge distributions originating
// in N (x,y) points in detection element
// detElemId and then fits them to get their position (xf,yf).
// Compares (xf,yf) to (x,y), which should be zero ideally.
func Closure(detElemId int, n int) {

}

// closureCmd represents the closure command
var closureCmd = &cobra.Command{
	Use:   "closure",
	Short: "Closure test : generate charge distribution and fit it",
	Run: func(cmd *cobra.Command, args []string) {
		Closure(detElemId, nSample)
	},
}

var detElemId int
var nSample int

func init() {
	rootCmd.AddCommand(closureCmd)

	closureCmd.Flags().IntVarP(&detElemId, "deid", "d", 100, "Detection element ID to consider")
	closureCmd.Flags().IntVarP(&nSample, "nsample", "n", 1, "Number of samples to generate")

}
