package main

import (
	"fmt"

	"github.com/aphecetche/pigiron/mapping"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate clusters",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("reimplement me !")
		// q := 1.0
		// cseg := segcache.CathodeSegmentation(deid, true)
		// chargeDist := qdist.MathiesonSt1
		//
		// pre := generate.PreClusterFromCharge(q, x, y, cseg, chargeDist)
		// fmt.Printf("%v", pre)
	},
}

var (
	x, y     float64
	charge   float64
	deid     int
	segcache mapping.SegCache
)

func init() {
	clusterCmd.AddCommand(createCmd)

	createCmd.Flags().IntVarP(&deid, "deid", "d", 100, "detection element ID of the cluster to be generated")
	createCmd.Flags().Float64VarP(&x, "xpos", "x", 0.0, "x position of the cluster to be generated")
	createCmd.Flags().Float64VarP(&y, "ypos", "y", 0.0, "y position of the cluster to be generated")
	createCmd.Flags().Float64VarP(&charge, "charge", "q", 1.0, "charge of the cluster to be generated")
}
