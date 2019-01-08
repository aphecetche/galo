package main

import (
	"log"
	"os"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/mathieson"
	"github.com/aphecetche/pigiron/mapping"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate clusters",
	Run: func(cmd *cobra.Command, args []string) {

		positions := []galo.ClusterPos{{x, y}, {x - 5, y - 5}}
		var charges []galo.ClusterCharge

		deid := mapping.DEID(deidFlag)

		var dgs []galo.DigitGroup

		minRelCharge := 1E-4
		q := charge / 2.0

		for _, p := range positions {
			dgs = append(dgs, mathieson.GenerateDigitGroup(deid, p.X, p.Y, q, minRelCharge))
			charges = append(charges, galo.ClusterCharge(q))
			q *= 4.0
		}
		clusters := galo.MockClustersFromDigitGroups(deid, positions, charges, dgs)

		dest, err := os.Create(output)
		if err != nil {
			log.Fatal(err)
		}
		defer dest.Close()
		to := NewClusterEncoder(dest, output)
		if to == nil {
			log.Fatalf("Could not get encoder for output file %v", output)
		}
		defer to.Close()

		to.Encode(&clusters)
	},
}

var (
	x, y     float64
	charge   float64
	deidFlag int
	output   string
)

func init() {
	clusterCmd.AddCommand(createCmd)

	createCmd.Flags().IntVarP(&deidFlag, "deid", "d", 100, "detection element ID of the cluster to be generated")
	createCmd.Flags().Float64VarP(&x, "xpos", "x", 0.0, "x position of the cluster to be generated")
	createCmd.Flags().Float64VarP(&y, "ypos", "y", 0.0, "y position of the cluster to be generated")
	createCmd.Flags().Float64VarP(&charge, "charge", "q", 1.0, "charge of the cluster to be generated")
	createCmd.Flags().StringVarP(&output, "output", "o", "", "Destination file")
}
