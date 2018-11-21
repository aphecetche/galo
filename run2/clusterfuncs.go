package run2

import (
	"fmt"
	"log"

	"github.com/aphecetche/pigiron/mapping"
)

type WeightModel int

const (
	NoWeight WeightModel = iota + 1
	RegularWeight
	SquaredWeight
)

// ClusterPosFunc computes the (x,y) position of a cluster
type ClusterPosFunc struct {
	F    func(*Cluster) (float64, float64)
	Name string
}

type ClusterSelFunc struct {
	F    func(*EventClusters, int) bool
	Name string
}

func strangeClusters(ec *EventClusters, i int) bool {
	silent := true
	var clu Cluster
	ec.E.Clusters(&clu, i)
	n := clu.Pre(nil).DigitsLength()
	if n > 200 {
		cluster2SVG(ec, i)
		if !silent {
			fmt.Println("WARNING", "cluster with", n, "digits")
			DumpEventClusters(ec)
			fmt.Println("")
		}
	}
	return false
}

func allClusters(ec *EventClusters, i int) bool {
	return true
}

func simpleClusters(ec *EventClusters, i int) bool {
	return ec.IsSimple(i)
}

func splitClusters(ec *EventClusters, i int) bool {
	return ec.IsSplit(i)
}

func dupClusters(ec *EventClusters, i int) bool {
	return ec.IsDup(i)
}

func cogNoWeight(clu *Cluster) (float64, float64) {
	return cog(clu.Pre(nil), NoWeight)
}

func cogWithRegularWeight(clu *Cluster) (float64, float64) {
	return cog(clu.Pre(nil), RegularWeight)
}

func cogWithSquaredWeight(clu *Cluster) (float64, float64) {
	return cog(clu.Pre(nil), SquaredWeight)
}

// cog compute the center of gravity of the digits within precluster
func cog(pre *PreCluster, weight WeightModel) (float64, float64) {
	var digit Digit
	var x, y, sumw float64
	for i := 0; i < pre.DigitsLength(); i++ {
		pre.Digits(&digit, i)
		deid := digit.Deid()
		manuid := int(digit.Manuid())
		seg := segcache.Segmentation(int(deid), manuid < 1024)
		manuchannel := int(digit.Manuchannel())
		paduid, err := seg.FindPadByFEE(manuid, manuchannel)
		if seg.IsValid(paduid) == false || err != nil {
			log.Fatalf("got invalid pad for DE %v MANU %v CH %v : %v -> paduid %v", deid, manuid, manuchannel, err, paduid)
		}
		var w float64
		if weight == NoWeight {
			w = 1.0
		} else {
			w = float64(digit.Adc())
		}
		sumw += w

		x += seg.PadPositionX(paduid) * w
		y += seg.PadPositionY(paduid) * w
	}
	x /= sumw
	y /= sumw
	return x, y
}

var segcache mapping.SegCache
