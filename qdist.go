package galo

import "github.com/aphecetche/pigiron/geo"

// ChargeIntegrator wraps the single Integrate method.
type ChargeIntegrator interface {
	Integrate(x1, y1, x2, y2 float64) float64
}

// Integrate a unit charge originating from (x,y) over the
// surface given by (lowerLeft,topRight)
func ChargeOverBox(x, y float64, integ ChargeIntegrator, bbox geo.BBox) float64 {
	x1 := bbox.Xmin() - float64(x)
	y1 := bbox.Ymin() - float64(y)
	x2 := x1 + bbox.Width()
	y2 := y1 + bbox.Height()
	return integ.Integrate(x1, y1, x2, y2)
}

// ChargeSpreader wraps the single SpreadCharge method.
type ChargeSpreader interface {
	// SpreadCharge spreads the charge q originating at position (x,y)
	// over several digits that are returned
	SpreadCharge(q, x, y float64) []Digit
}
