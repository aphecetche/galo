package galo

type ClusterPos struct {
	X float64
	Y float64
	Z float64
}

type Cluster struct {
	Pre PreCluster
	Pos ClusterPos
}

type Clusters []Cluster

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

type ClusterDecoder interface {
	// Decode reads the next cluster from its input and stores it
	// in the value pointed by clu.
	Decode(clu *Cluster) error
}
