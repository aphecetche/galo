package main

import (
	"github.com/aphecetche/galo"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Just reads input file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		LoopOverFile(args[0], "any", func(index int, tc *galo.TaggedClusters, selected []int) {})
	},
}

func init() {
	clusterCmd.AddCommand(readCmd)
}
