package mathieson_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/mathieson"
	"github.com/aphecetche/pigiron/mapping"
	"github.com/gonum/floats"
	"go-hep.org/x/hep/hbook"
	"gonum.org/v1/gonum/optimize"
)

type testHit struct {
	x, y, q float64
}

type testDEHit struct {
	deID mapping.DEID
	hits []testHit
}

func createCluster(de testDEHit) galo.DEClusters {
	minRelCharge := 1E-4
	deid := mapping.DEID(de.deID)
	var dgs []galo.DigitGroup
	var positions []galo.ClusterPos
	var charges []galo.ClusterCharge
	for _, hit := range de.hits {
		dgs = append(dgs, mathieson.GenerateDigitGroup(deid, hit.x, hit.y, hit.q, minRelCharge))
		positions = append(positions, galo.ClusterPos{X: hit.x, Y: hit.y})
		charges = append(charges, galo.ClusterCharge(hit.q))
	}
	return galo.MockClustersFromDigitGroups(deid, positions, charges, dgs)
}

func newFitter(deid mapping.DEID) galo.DEClusterPositioner {
	return newFitterApprox(deid, math.Atan, math.Tanh)
}

func newFitterApprox(deid mapping.DEID, atanfunc, tanhfunc func(float64) float64) galo.DEClusterPositioner {
	var method optimize.Method = &optimize.NelderMead{}
	if deid < 500 {
		m := mathieson.NewMathieson2DApprox(0.21, 0.7*0.7, 0.755*0.755, atanfunc, tanhfunc)
		return mathieson.NewClusterFitter(*m, method)
	}
	m := mathieson.NewMathieson2DApprox(0.25, 0.7131*0.7131, 0.7642*0.7642, atanfunc, tanhfunc)
	return mathieson.NewClusterFitter(*m, method)
}

func TestBasicMathiesonFit(t *testing.T) {

	for _, tp := range []testDEHit{
		{100, []testHit{{24.0, 72.0, 50.0}}},
		{500, []testHit{{20.0, 0, 50.0}}},
	} {
		declu := createCluster(tp)
		fitter := newFitter(declu.DeID)
		testFitClusters(declu, fitter, t)
	}
}

func noisify(declu galo.DEClusters, noiseFraction float64) galo.DEClusters {
	var digits []galo.Digit
	var clusters []galo.Cluster

	for _, clu := range declu.Clusters {
		for _, d := range clu.Pre.Digits {
			dq := d.Q * (1.0 + rand.NormFloat64()*noiseFraction)
			if dq < 0 {
				dq = 0.0
			}
			digits = append(digits, galo.Digit{ID: d.ID, Q: dq})
		}
		clusters = append(clusters, galo.Cluster{Pre: galo.PreCluster{galo.DigitGroup{0, digits}}, Pos: clu.Pos, Q: clu.Q})
	}
	return galo.DEClusters{DeID: declu.DeID, Clusters: clusters}
}

func TestNoisyMathiesonFit(t *testing.T) {
	N := 1000
	for _, tp := range []testDEHit{
		{100, []testHit{{24.0, 72.0, 50.0}}},
		// {500, []testHit{{20.0, 0, 50.0}}},
	} {

		h := hbook.NewH1D(128, 0, 1E3)
		declu := createCluster(tp)
		fitter := newFitterApprox(declu.DeID, mathieson.AtanEq11, mathieson.TanhApprox1)
		for _, noise := range []float64{10} {
			for i := 0; i < N; i++ {
				noisy := noisify(declu, noise/100.0)
				res := 1E4 * galo.DEClusterResidual(&noisy, 0, fitter)
				h.Fill(res, 1.0)
			}
			p := galo.PlotResidual(h)
			s := fmt.Sprintf("Noise%3.1fPercent", noise)
			galo.SavePlot(p, "TestNoisyMathiesonFit", s)
		}
	}
}

func testFitClusters(declu galo.DEClusters, fitter galo.DEClusterPositioner,
	t *testing.T) {
	for i, clu := range declu.Clusters {
		x, y := fitter.Position(&declu, i)
		x0 := clu.Pos.X
		y0 := clu.Pos.Y
		dx := x - x0
		dy := y - y0
		d := math.Sqrt(dx*dx + dy*dy)
		tol := 1E-4 // 1 micron in centimeters
		if !floats.EqualWithinAbs(x, x0, tol) {
			t.Errorf("Want x=%10.4f Got %10.4f", x0, x)
		}
		if !floats.EqualWithinAbs(y, y0, tol) {
			t.Errorf("Want y=%10.4f Got %10.4f", y0, y)
		}
		if d > tol {
			t.Errorf("Want d=0 Got %10.4f", d)
		}
	}
}

func BenchmarkMathiesonFit(b *testing.B) {
	declu := createCluster(testDEHit{deID: 100, hits: []testHit{{24, 72, 50}}})

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

// func BenchmarkMathiesonApprox(b *testing.B) {
//
// 	for _, tp := range []testDEHit{
// 		{100, []testHit{{24.0, 72.0, 50.0}}},
// 	} {
// 		declu := createCluster(tp)
// 		fitter := newFitter(declu.DEID, atanfunc, tanhfunc)
// 		b.Run(m.name, func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				_, _ = fitter.Position(&declu, 0)
// 			}
// 		})
// 	}
// }

func generateTestPoints(deid mapping.DEID, N int) []testDEHit {
	var th []testDEHit
	seg := galo.SegCache.Segmentation(deid)
	box := mapping.ComputeSegmentationBBox(seg)
	for i := 0; i < N; i++ {
		q := 100.0 + rand.Float64()*50.0
		x := box.Xmin() + rand.Float64()*box.Width()
		y := box.Ymin() + rand.Float64()*box.Height()
		b, nb, err := seg.FindPadPairByPosition(x, y)
		if err != nil {
			// not a valid pad
			continue
		}
		if !seg.IsValid(b) && !seg.IsValid(nb) {
			// discard monocathode stuff
			continue
		}
		th = append(th, testDEHit{deid, []testHit{{x, y, q}}})
	}
	return th
}

func TestMathiesonApprox(t *testing.T) {
	h := hbook.NewH1D(128, 0, 1E3)
	testpoints := generateTestPoints(100, 100)
	for _, tp := range testpoints {
		declu := createCluster(tp)
		refFitter := newFitter(declu.DeID)
		fitter := newFitterApprox(declu.DeID, math.Atan, mathieson.TanhApprox2)
		res := 1E4 * galo.DEClusterResidual(&declu, 0, fitter)
		refRes := 1E4 * galo.DEClusterResidual(&declu, 0, refFitter)
		fmt.Printf("res=%8.4f microns ref=%8.4f mult=%d\n",
			res,
			refRes,
			declu.Clusters[0].Pre.NofPads())
		h.Fill(res, 1.0)
	}
	p := galo.PlotResidual(h)
	galo.SavePlot(p, "TestNoisyMathieson", "Approx")
}
