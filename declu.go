package galo

import (
	"math"
	"strconv"
	"strings"

	"github.com/aphecetche/pigiron/mapping"
)

// DEClusters represents a group of clusters for one detection element.
type DEClusters struct {
	// DeID is the detection element id of those clusters.
	DeID mapping.DEID
	// Clusters is the list of clusters of this group.
	Clusters []Cluster
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

// DEClusterPositioner computes the 2D position of clusters
type DEClusterPositioner interface {
	// Position (re)computes the position of the cluster.
	Position(declu *DEClusters, i int) (x, y float64)
	// Name of the positioner.
	Name() string
}

func NewDEClusterPositioner(name string) DEClusterPositioner {
	s := strings.ToUpper(name)
	if s == "COG" {
		return DEClusterPositionerCOG{Wmod: NoWeight}
	}
	return nil
}

type WeightModel int

const (
	NoWeight WeightModel = iota + 1
	RegularWeight
	SquaredWeight
)

type DEClusterPositionerCOG struct {
	Wmod WeightModel
}

func (declupos DEClusterPositionerCOG) Position(declu *DEClusters, i int) (float64, float64) {
	var x, y, sumw float64
	seg := SegCache.Segmentation(declu.DeID)
	clu := declu.Clusters[i]
	for _, d := range clu.Pre.Digits {
		dsid := seg.PadDualSampaID(d.ID)
		dsch := seg.PadDualSampaChannel(d.ID)
		paduid, err := seg.FindPadByFEE(dsid, dsch)
		if err != nil {
			panic(err)
		}
		var w float64
		if declupos.Wmod == NoWeight {
			w = 1.0
		} else {
			w = d.Q
		}
		sumw += w

		x += seg.PadPositionX(paduid) * w
		y += seg.PadPositionY(paduid) * w
	}
	x /= sumw
	y /= sumw
	return x, y
}

func (clupos DEClusterPositionerCOG) Name() string {
	return "COG" + strconv.Itoa(int(clupos.Wmod))
}

func DEClusterResidual(declu *DEClusters, i int, declupos DEClusterPositioner) float64 {
	x, y := declupos.Position(declu, i)
	pos := declu.Clusters[i].Pos
	dx := x - float64(pos.X)
	dy := y - float64(pos.Y)
	return math.Sqrt(dx*dx + dy*dy)
}
