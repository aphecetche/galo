package mathieson

import (
	"log"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
)

// mathiesonFitter implements galo.DEClusterPositioner.
type mathiesonFitter struct {
	mat2D  Mathieson2D // F implements galo.ChargeIntegrator
	method optimize.Method
}

func NewClusterFitter(mat Mathieson2D, method optimize.Method) *mathiesonFitter {
	return &mathiesonFitter{mat2D: mat, method: method}
}

func (cp *mathiesonFitter) Position(declu *galo.DEClusters, i int) (x, y float64) {
	clu := declu.Clusters[i]

	seg := galo.SegCache.Segmentation(declu.DeID)
	var charges []float64
	var boxes []geo.BBox

	for _, d := range clu.Pre.Digits {
		charges = append(charges, d.Q)
		b := mapping.ComputePadBBox(seg, d.ID)
		boxes = append(boxes, b)
	}

	n := 0

	fcn := func(pos []float64) float64 {
		n++
		x := pos[0]
		y := pos[1]
		q := pos[2] / 2.0 // FIXME: should take into account the fact that charge splitting is not perfect (i.e. not 50% on each cathode all the time).
		lnL := 0.0
		for i, b := range boxes {
			dq := q*galo.ChargeOverBox(x, y, &cp.mat2D, b) - charges[i]
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

	cog := galo.NewDEClusterPositioner("cog")

	x0, y0 := cog.Position(declu, i)
	q0 := clu.Pre.Charge()

	var p0 = []float64{x0, y0, q0}

	res, err := optimize.Minimize(p, p0, nil, cp.method)
	if err != nil {
		log.Fatal(err)
	}

	return res.X[0], res.X[1]
}

func (cp *mathiesonFitter) Name() string {
	return "MathiesonFitter"
}
