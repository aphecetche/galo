package mathieson

import (
	"github.com/aphecetche/galo"
)

// 2D Mathieson for St1 (pitch 2.1 mm)
var MathiesonSt1 Mathieson2D = *(NewMathieson2D(0.21, 0.700*0.700, 0.755*0.755))

// 2D Mathieson for St2 and St345 (pitch 2.5 mm)
var MathiesonSt2345 Mathieson2D = *(NewMathieson2D(0.25, 0.7131*0.7131, 0.7642*0.7642))

func NewChargeIntegrator(deid int) galo.ChargeIntegrator {
	if deid < 300 {
		return MathiesonSt1
	}
	return MathiesonSt2345
}
