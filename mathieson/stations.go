package mathieson

import (
	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/mapping"
)

// 2D Mathieson for St1 (pitch 2.1 mm)
var St1 Mathieson2D = *(NewMathieson2D(0.21, 0.700*0.700, 0.755*0.755))

// 2D Mathieson for St2 and St345 (pitch 2.5 mm)
var St2345 Mathieson2D = *(NewMathieson2D(0.25, 0.7131*0.7131, 0.7642*0.7642))

func NewChargeIntegrator(deid mapping.DEID) galo.ChargeIntegrator {
	if deid < 300 {
		return St1
	}
	return St2345
}
