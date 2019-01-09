package mathieson

import (
	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/mapping"
)

type Station struct {
	K3x   float64
	K3y   float64
	Pitch float64 //cm
}

var (
	St1    Station = Station{0.49, 0.570025, 0.21}
	St2345 Station = Station{0.50851161, 0.58400164, 0.25}
)

func integrator1DPair(st Station, impl IntegrateImpl) (Integrate1DFunc, Integrate1DFunc) {
	return Integrator1D(st.Pitch, st.K3x, impl),
		Integrator1D(st.Pitch, st.K3y, impl)
}

func newIntegrateFunc(fx, fy Integrate1DFunc) galo.IntegrateFunc {
	return func(x1, y1, x2, y2 float64) float64 {
		return 4.0 * fx(x1, x2) * fy(y1, y2)
	}
}

func NewChargeIntegrator(deid mapping.DEID, impl IntegrateImpl) galo.ChargeIntegrator {
	var fx, fy Integrate1DFunc
	st := St2345
	if deid < 300 {
		st = St1
	}
	fx, fy = integrator1DPair(st, impl)
	return newIntegrateFunc(fx, fy)
}
