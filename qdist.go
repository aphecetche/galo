package galo

// ChargeIntegrator wraps the single Integrate method.
type ChargeIntegrator interface {
	Integrate(x1, y1, x2, y2 float64) float64
}

// Integrate a unit charge originating from (x,y) over the
// surface given by (lowerLeft,topRight)
func ChargeOverBox(x, y float64, integ ChargeIntegrator, lowerLeft, topRight XY) float64 {
	x1 := lowerLeft.X - float64(x)
	y1 := topRight.Y - float64(y)
	x2 := x1 + (topRight.X - lowerLeft.X)
	y2 := y1 + (topRight.Y - lowerLeft.Y)
	return integ.Integrate(x1, y1, x2, y2)
}

// ChargeSpreader wraps the single SpreadCharge method.
type ChargeSpreader interface {
	// SpreadCharge spreads the charge q originating at position (x,y)
	// over several digits that are returned
	SpreadCharge(q, x, y float64) []Digit
}
