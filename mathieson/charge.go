package mathieson

import (
	"log"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/mapping"
)

const minRelCharge float64 = 1E-4

type ChargeSpreadFunc func(q, x, y float64) []galo.Digit

func (f ChargeSpreadFunc) SpreadCharge(q, x, y float64) []galo.Digit {
	return f(q, x, y)
}

func padBox(cseg mapping.CathodeSegmentation, paduid mapping.PadUID) (lowerLeft, topRight galo.XY) {

	px := cseg.PadPositionX(paduid)
	py := cseg.PadPositionY(paduid)
	dx := cseg.PadSizeX(paduid)
	dy := cseg.PadSizeY(paduid)
	return galo.XY{px - dx/2.0, py - dy/2.0}, galo.XY{px + dx/2.0, py + dy/2.0}
}

func NewMathiesonChargeSpreader(deid int) ChargeSpreadFunc {
	cseg := mapping.NewCathodeSegmentation(deid, true)
	integ := NewChargeIntegrator(deid)
	return func(q, x, y float64) []galo.Digit {
		var digits galo.Digits
		deid := cseg.DetElemID()
		paduid, err := cseg.FindPadByPosition(x, y)
		if err != nil {
			log.Fatalf("Could not find at (%v,%v)  DE %v", x, y, deid)
		}
		neighbours := cseg.GetNeighbours(paduid)
		for _, nei := range neighbours {
			neighbours = append(neighbours, cseg.GetNeighbours(nei)...)
		}
		pids := make(map[mapping.PadUID]struct{})
		for _, nei := range neighbours {
			pids[nei] = struct{}{}
		}

		for paduid, _ := range pids {
			lowerLeft, topRight := padBox(cseg, paduid)
			dq := galo.ChargeOverBox(x, y, integ, lowerLeft, topRight)
			if dq < minRelCharge {
				continue
			}
			digits = append(digits, galo.Digit{ID: int(paduid), Q: q * dq})
		}
		return digits
	}
}
