package mathieson

import (
	"fmt"

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

		paduids := make(map[mapping.PadUID]struct{})

		const sx float64 = 0.5 // cm
		const sy float64 = 0.5 // cm // FIXME: where to get those from ?

		xmin := x - sx
		xmax := x + sx
		ymin := y - sy
		ymax := y + sy

		seg.ForEachPadInArea(xmin, ymin, xmax, ymax, func(paduid mapping.PadUID) {
			paduids[paduid] = struct{}{}
		})

		fmt.Println("# paduids=", len(paduids), paduids)

		var xpadmin, ypadmin, xpadmax, ypadmax float64
		for paduid, _ := range paduids {
			mapping.ComputePadBBox(seg, mapping.PadUID(paduid), &xpadmin, &ypadmin, &xpadmax, &ypadmax)
			dq := 0.5 * galo.ChargeOverBox(x, y, integ, xpadmin, ypadmin, xpadmax, ypadmax)
			if dq < minRelCharge {
				continue
			}
			digits = append(digits, galo.Digit{ID: paduid, Q: q * dq})
		}
		return digits
	}
}
