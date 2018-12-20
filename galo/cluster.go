package main

import (
	"github.com/spf13/cobra"
)

// clusterCmd represents the clusters command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Various operations related to clusters",
}

var maxEvents int
var cluSelArg string
var firstEvent int

func init() {
	rootCmd.AddCommand(clusterCmd)

	// persistent flag is for this command and all its sub-commands
	clusterCmd.PersistentFlags().IntVarP(&maxEvents, "max-events", "m", 100000000, "Maximum number of events to process")
	clusterCmd.PersistentFlags().StringVarP(&cluSelArg, "select", "s", "any", "Clusters to be used")
	clusterCmd.PersistentFlags().IntVarP(&firstEvent, "first-event", "a", 0, "First event to consider")
}
