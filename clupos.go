package galo

import (
	"math"
	"strconv"

	"github.com/aphecetche/pigiron/mapping"
)

type WeightModel int

const (
	NoWeight WeightModel = iota + 1
	RegularWeight
	SquaredWeight
)

type ClusterPositionerCOG struct {
	Wmod      WeightModel
	padPos    mapping.PadSizerPositioner
	padFinder mapping.PadByFEEFinder
	padLoc    mapping.PadFEELocator
}

func (clupos ClusterPositionerCOG) Position(clu Cluster) (float64, float64) {
	var x, y, sumw float64
	for _, d := range clu.Pre.Digits {
		dsid := clupos.padLoc.PadDualSampaID(d.ID)
		dsch := clupos.padLoc.PadDualSampaChannel(d.ID)
		paduid, err := clupos.padFinder.FindPadByFEE(dsid, dsch)
		if err != nil {
			panic(err)
		}
		var w float64
		if clupos.Wmod == NoWeight {
			w = 1.0
		} else {
			w = d.Q
		}
		sumw += w

		x += clupos.padPos.PadPositionX(paduid) * w
		y += clupos.padPos.PadPositionY(paduid) * w
	}
	x /= sumw
	y /= sumw
	return x, y
}

func (clupos ClusterPositionerCOG) Name() string {
	return "COG" + strconv.Itoa(int(clupos.Wmod))
}

func ClusterResidual(clu Cluster, cluposfunc ClusterPositioner) float64 {
	x, y := cluposfunc.Position(clu)
	pos := clu.Pos
	dx := x - float64(pos.X)
	dy := y - float64(pos.Y)
	return math.Sqrt(dx*dx + dy*dy)
}
