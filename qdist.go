package galo

// ChargeIntegrator wraps the single Integrate method.
type ChargeIntegrator interface {
	Integrate(x1, y1, x2, y2 float64) float64
}

// IntegrateFunc functions implement ChargeIntegrator.
type IntegrateFunc func(x1, y1, x2, y2 float64) float64

func (f IntegrateFunc) Integrate(x1, y1, x2, y2 float64) float64 {
	return f(x1, y1, x2, y2)
}

// Integrate a unit charge originating from (x,y) over the
// surface given by xmin,ymin,xmax,ymax
func ChargeOverBox(x, y float64, integ ChargeIntegrator, xmin, ymin, xmax, ymax float64) float64 {
	x1 := xmin - float64(x)
	y1 := ymin - float64(y)
	x2 := x1 + xmax - xmin
	y2 := y1 + ymax - ymin
	return integ.Integrate(x1, y1, x2, y2)
}

// ChargeSpreader wraps the single SpreadCharge method.
type ChargeSpreader interface {
	// SpreadCharge spreads the charge q originating at position (x,y)
	// over several digits that are returned
	SpreadCharge(q, x, y float64) []Digit
}
