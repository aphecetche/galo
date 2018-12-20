package mathieson_test

import (
	"math"
	"testing"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/mathieson"
	"github.com/aphecetche/pigiron/mapping"
	"github.com/gonum/floats"
	"gonum.org/v1/gonum/optimize"
)

func createCluster(positions []galo.ClusterPos, charges []float64) galo.DEClusters {
	deid := mapping.DEID(100)
	var dgs []galo.DigitGroup
	for i, p := range positions {
		dgs = append(dgs, mathieson.GenerateDigitGroup(deid, p.X, p.Y, charges[i]))
	}
	return galo.MockClustersFromDigitGroups(deid, positions, charges, dgs)
}

func TestMathiesonFit(t *testing.T) {

	positions := []galo.ClusterPos{{72, 24}}
	charges := []float64{50.0}

	declu := createCluster(positions, charges)

	fitter := mathieson.NewClusterFitter(mathieson.St1, &optimize.BFGS{})

	x, y := fitter.Position(&declu, 0)

	dx := x - positions[0].X
	dy := y - positions[0].Y
	d := math.Sqrt(dx*dx + dy*dy)

	tol := 1E-4 // 1 micron in centimeters

	if !floats.EqualWithinAbs(x, positions[0].X, tol) {
		t.Errorf("Want x=%10.4f Got %10.4f", positions[0].X, x)
	}
	if !floats.EqualWithinAbs(y, positions[0].Y, tol) {
		t.Errorf("Want y=%10.4f Got %10.4f", positions[0].Y, y)
	}
	if d > tol {
		t.Errorf("Want d=0 Got %10.4f", d)
	}
}

func BenchmarkMathiesonFit(b *testing.B) {
	positions := []galo.ClusterPos{{72, 24}}
	charges := []float64{50.0}
	declu := createCluster(positions, charges)

	methods := []struct {
		name   string
		method optimize.Method
	}{
		{"BFGS", &optimize.BFGS{}},
		{"NelderMead", &optimize.NelderMead{}},
	}

	for _, m := range methods {
		b.Run(m.name, func(b *testing.B) {
			fitter := mathieson.NewClusterFitter(mathieson.St1, m.method)
			for i := 0; i < b.N; i++ {
				_, _ = fitter.Position(&declu, 0)
			}
		})
	}
}
