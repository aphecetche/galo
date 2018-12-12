package galo

import (
	"fmt"

	"github.com/aphecetche/pigiron/mapping"
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

// DEClusters represents a group of clusters for one detection element.
type DEClusters struct {
	// DeID is the detection element id of those clusters.
	DeID mapping.DEID
	// Clusters is the list of clusters of this group.
	Clusters []Cluster
}

// Clusterizer wraps the single Clusterizer method.
type Clusterizer interface {
	// Clusterize converts a precluster into one or several clusters.
	Clusterize(pre PreCluster) []Cluster
}

// ClusterSelector selects or discard a cluster based on some criteria.
type ClusterSelector interface {
	// Select decides if a cluster is to be kept or not.
	Select(clu Cluster) bool
	// Name of the selector.
	Name() string
}

// ClusterCharger computes the position of a cluster.
type ClusterCharger interface {
	// Charge (re)computes the charge of the cluster.
	Charge(clu Cluster) float64
	// Name of the charger
	Name() string
}

// ClusterPositioner computes the 2D position of a cluster.
type ClusterPositioner interface {
	// Position (re)computes the position of the cluster.
	Position(clu Cluster) (x, y float64)
	// Name of the positioner.
	Name() string
}

// DEClustersDecoder decodes DEClusters from an underlying stream.
type DEClustersDecoder interface {
	// Decode reads the next DEClusters from its input and stores it
	// in the value pointed by clu.
	Decode(declu *DEClusters) error
	// Close may be necessary for those decoder implementations.
	Close()
}

// DEClustersDecoder encodes DEClusters onto the underlying stream.
type DEClustersEncoder interface {
	// Encode writes the encoding of clu to the stream.
	Encode(declu *DEClusters) error
	// Close may be necessary for those decoder implementations to flush data.
	Close()
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

// // String gets a string representation of the DEClusters.
// func (declu DEClusters) String() string {
// 	s := fmt.Sprintf("DE %4d - %d clusters", declu.DeID, len(declu.Clusters))
// 	for _, c := range declu.Clusters {
// 		s += fmt.Sprintf("%v", c)
// 	}
// 	return s
// }
