package galo

import (
	"fmt"

	"github.com/aphecetche/pigiron/mapping"
	"github.com/gonum/floats"
)

// Cluster represents a MCH cluster.
type Cluster struct {
	// Pre is the PreCluster this cluster originates from.
	Pre PreCluster
	// Pos is the cluster position.
	Pos ClusterPos
	// Q is the total charge of the cluster.
	Q ClusterCharge
}

// ClusterPos represents the 2D position of a cluster.
// The position is relative to one detection element.
type ClusterPos struct {
	X float64
	Y float64
}

// ClusterCharge represents the total charge
// of a cluster.
type ClusterCharge float64

// Clusterizer wraps the single Clusterizer method.
type Clusterizer interface {
	// Clusterize converts a precluster into one or several clusters.
	Clusterize(pre PreCluster) []Cluster
}

// ClusterCharger computes the position of a cluster.
type ClusterCharger interface {
	// Charge (re)computes the charge of the cluster.
	Charge(clu Cluster) float64
	// Name of the charger
	Name() string
}

// String gets a string representation of the position.
func (pos ClusterPos) String() string {
	return fmt.Sprintf("X %7.2f Y %7.2f", pos.X, pos.Y)
}

// String gets a string representation of the cluster.
func (clu Cluster) String() string {
	s := fmt.Sprintf(" Q=%7.2f", clu.Q)
	s += fmt.Sprintf(" Pos=%v", clu.Pos)
	s += fmt.Sprintf(" Pre=%v", clu.Pre)
	return s
}

// SameCluster returns true if the two clusters :
// - have the same precluster
// - have close enough positions
func SameCluster(ca, cb Cluster) bool {
	pa := ca.Pre
	pb := cb.Pre
	if !SamePreCluster(pa, pb) {
		return false
	}
	const tol = 1E-6
	return floats.EqualWithinAbs(float64(ca.Pos.X), float64(cb.Pos.X), tol) &&
		floats.EqualWithinAbs(float64(ca.Pos.Y), float64(cb.Pos.Y), tol)
}

func MockClustersFromDigitGroups(deid mapping.DEID, positions []ClusterPos, charges []ClusterCharge, dgs []DigitGroup) DEClusters {
	var clusters []Cluster

	for i, dg := range dgs {
		pre := PreCluster{DigitGroup: dg}
		pos := ClusterPos{X: positions[i].X, Y: positions[i].Y}
		clu := Cluster{Pre: pre, Pos: pos, Q: charges[i]}
		clusters = append(clusters, clu)
	}

	return DEClusters{DeID: deid, Clusters: clusters}
}
