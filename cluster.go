package galo

import (
	"fmt"

	"github.com/aphecetche/pigiron/mapping"
)

type ClusterPos struct {
	X float64
	Y float64
}

type ClusterCharge float64

type Cluster struct {
	Pre PreCluster
	Pos ClusterPos
	Q   ClusterCharge
}

type DEClusters struct {
	DeID     mapping.DEID
	Clusters []Cluster
}

type Clusterizer interface {
	// Clusterize converts a precluster into one or several clusters.
	Clusterize(pre PreCluster) []Cluster
}

type ClusterSelector interface {
	// Select decides if a cluster is to be kept or not.
	Select(clu Cluster) bool
}

type ClusterCharger interface {
	// Charge (re)computes the charge of the cluster.
	Charge(clu Cluster) float64
}

type ClusterPositioner interface {
	// Position (re)computes the position of the cluster.
	Position(clu Cluster) (x, y float64)
}

type DEClustersDecoder interface {
	// Decode reads the next DEClusters from its input and stores it
	// in the value pointed by clu.
	Decode(declu *DEClusters) error
	Close()
}

type DEClustersEncoder interface {
	// Encode writes the encoding of clu to the stream.
	Encode(declu *DEClusters) error
	Close()
}

func (pos ClusterPos) String() string {
	return fmt.Sprintf("X %7.2f Y %7.2f", pos.X, pos.Y)
}
func (clu Cluster) String() string {
	s := fmt.Sprintf(" Q=%7.2f", clu.Q)
	s += fmt.Sprintf(" Pos=%v", clu.Pos)
	s += fmt.Sprintf(" Pre=%v", clu.Pre)
	return s
}
