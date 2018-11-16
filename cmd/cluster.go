package cmd

import (
	"github.com/spf13/cobra"
)

// clusterCmd represents the clusters command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Inspect clusters",
}

var maxEvents int

func init() {
	rootCmd.AddCommand(clusterCmd)

	// persistent flag is for this command and all its sub-commands
	clusterCmd.PersistentFlags().IntVarP(&maxEvents, "max-events", "m", 100000000, "Maximum number of events to process")
}

func init() {
	clusterCmd.AddCommand(plotCmd)
}
