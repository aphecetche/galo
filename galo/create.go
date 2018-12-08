package main

import (
	"os"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/mathieson"
	"github.com/aphecetche/galo/yaml"
	"github.com/aphecetche/pigiron/mapping"
	"github.com/spf13/cobra"
)

type Pos struct {
	X, Y float64
}

func generateDigitGroup(deid mapping.DEID, positions []Pos, charges []float64) []galo.DigitGroup {
	minRelCharge := 1E-3
	cs := mathieson.NewMathiesonChargeSpreader(deid, minRelCharge)
	var digitGroups []galo.DigitGroup
	for i, pos := range positions {
		digits := cs.SpreadCharge(charges[i], pos.X, pos.Y)
		digitGroups = append(digitGroups, galo.DigitGroup{RefTime: 0, Digits: digits})
	}
	return digitGroups
}

func MockClustersFromDigitGroups(deid mapping.DEID, positions []Pos, charges []float64, dgs []galo.DigitGroup) galo.DEClusters {
	var clusters []galo.Cluster

	for i, dg := range dgs {
		pre := galo.PreCluster{DigitGroup: dg}
		// FIXME: cluster pos should be absolute, not relative to the de ?
		// or keep it as (X,Y) always at the level of cluster and let
		// the tracking stage do to geo transformation ?
		pos := galo.ClusterPos{X: positions[i].X, Y: positions[i].Y}
		clu := galo.Cluster{Pre: pre, Pos: pos, Q: galo.ClusterCharge(charges[i])}
		clusters = append(clusters, clu)
	}

	return galo.DEClusters{DeID: deid, Clusters: clusters}
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate clusters",
	Run: func(cmd *cobra.Command, args []string) {

		positions := []Pos{{x, y}, {x - 5, y - 5}}
		charges := []float64{charge * 2.0, charge / 2.0}

		deid := mapping.DEID(deidFlag)

		dgs := generateDigitGroup(deid, positions, charges)
		clusters := MockClustersFromDigitGroups(deid, positions, charges, dgs)

		// svg := svg.NewClusterEncoder(os.Stdout)
		// svg.DefaultStyle()
		// svg.MoveToOrigin()
		// svg.WithCharge()
		// defer svg.Close()
		// svg.Encode(&clusters)

		yaml := yaml.NewClusterEncoder(os.Stdout,
			func(deid mapping.DEID) mapping.PadFEELocator {
				return galo.SegCache.Segmentation(deid)
			})
		defer yaml.Close()
		yaml.Encode(&clusters)
	},
}

var (
	x, y     float64
	charge   float64
	deidFlag int
)

func init() {
	clusterCmd.AddCommand(createCmd)

	createCmd.Flags().IntVarP(&deidFlag, "deid", "d", 100, "detection element ID of the cluster to be generated")
	createCmd.Flags().Float64VarP(&x, "xpos", "x", 0.0, "x position of the cluster to be generated")
	createCmd.Flags().Float64VarP(&y, "ypos", "y", 0.0, "y position of the cluster to be generated")
	createCmd.Flags().Float64VarP(&charge, "charge", "q", 1.0, "charge of the cluster to be generated")
}
