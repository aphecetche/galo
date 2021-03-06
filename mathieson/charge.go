package mathieson

import (
	"log"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/mapping"
	_ "github.com/aphecetche/pigiron/mapping/impl4"
)

type ChargeSpreadFunc func(q, x, y float64) []galo.Digit

func (f ChargeSpreadFunc) SpreadCharge(q, x, y float64) []galo.Digit {
	return f(q, x, y)
}

func NewChargeSpreader(deid mapping.DEID, minRelCharge float64) ChargeSpreadFunc {
	seg := galo.SegCache.Segmentation(deid)
	integ := NewChargeIntegrator(deid, IntegrateImplDefault)
	return func(q, x, y float64) []galo.Digit {
		var digits []galo.Digit
		deid := seg.DetElemID()
		pb, pnb, err := seg.FindPadPairByPosition(x, y)
		if err != nil {
			log.Fatalf("Could not find at (%v,%v)  DE %v", x, y, deid)
		}
		neighbours := seg.GetNeighbours(pb)
		neighbours = append(neighbours, seg.GetNeighbours(pnb)...)
		for _, nei := range neighbours {
			neighbours = append(neighbours, seg.GetNeighbours(nei)...)
		}
		pids := make(map[mapping.PadUID]struct{})
		for _, nei := range neighbours {
			pids[nei] = struct{}{}
		}

		for paduid, _ := range pids {
			dq := 0.5 * galo.ChargeOverBox(x, y, integ, mapping.ComputePadBBox(seg, paduid))
			if dq < minRelCharge {
				continue
			}
			digits = append(digits, galo.Digit{ID: paduid, Q: q * dq})
		}
		return digits
	}
}
