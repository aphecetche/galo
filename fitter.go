package galo

import (
	"log"

	"github.com/aphecetche/pigiron/mapping"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
)

// Fitter implements DEClusterPositioner.
type Fitter struct {
	Integrator ChargeIntegrator
	Method     optimize.Method
}

func (f *Fitter) Position(declu *DEClusters, i int) (x, y float64) {
	clu := declu.Clusters[i]

	seg := SegCache.Segmentation(declu.DeID)
	var charges []float64
	var xmin, ymin, xmax, ymax []float64
	var x1, y1, x2, y2 float64

	for _, d := range clu.Pre.Digits {
		charges = append(charges, d.Q)
		mapping.ComputePadBBox(seg, d.ID, &x1, &y1, &x2, &y2)
		xmin = append(xmin, x1)
		xmax = append(xmax, x2)
		ymin = append(ymin, y1)
		ymax = append(ymin, y2)
	}

	n := 0

	fcn := func(pos []float64) float64 {
		n++
		x := pos[0]
		y := pos[1]
		q := pos[2] / 2.0 // FIXME: should take into account the fact that charge splitting is not perfect (i.e. not 50% on each cathode all the time).
		lnL := 0.0
		for i, _ := range charges {
			dq := q*ChargeOverBox(x, y, f.Integrator, xmin[i], ymin[i], xmax[i], ymax[i]) - charges[i]
			lnL += dq * dq
		}
		return lnL
	}

	// to make the minimization, needs a couple of things :
	// - the objective function fcn to minimize
	// - the gradient of fcn
	// - the hessian of fcn
	// - an optimization method

	grad := func(grad, x []float64) {
		fd.Gradient(grad, fcn, x, nil)
	}

	hess := func(h mat.MutableSymmetric, x []float64) {
		fd.Hessian(h.(*mat.SymDense), fcn, x, nil)
	}

	p := optimize.Problem{
		Func: fcn,
		Grad: grad,
		Hess: hess,
	}

	cog := NewDEClusterPositioner("cog")

	x0, y0 := cog.Position(declu, i)
	q0 := clu.Pre.Charge()

	var p0 = []float64{x0, y0, q0}

	res, err := optimize.Minimize(p, p0, nil, f.Method)
	if err != nil {
		log.Fatal(err)
	}

	return res.X[0], res.X[1]
}

func (f *Fitter) Name() string {
	return "Fitter"
}
