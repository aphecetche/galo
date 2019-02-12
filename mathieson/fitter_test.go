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
	"go-hep.org/x/hep/hplot"
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
	cs := mathieson.NewChargeSpreader(deid, minRelCharge)
	return createClustersFromChargeSpreader(de, cs)
}

func createClustersFromChargeSpreader(de testDEHit, cs galo.ChargeSpreader) galo.DEClusters {
	deid := mapping.DEID(de.deID)
	var dgs []galo.DigitGroup
	var positions []galo.ClusterPos
	var charges []galo.ClusterCharge
	for _, hit := range de.hits {
		digits := cs.SpreadCharge(hit.q, hit.x, hit.y)
		dgs = append(dgs, galo.DigitGroup{RefTime: 0, Digits: digits})
		positions = append(positions, galo.ClusterPos{X: hit.x, Y: hit.y})
		charges = append(charges, galo.ClusterCharge(hit.q))
	}
	return galo.MockClustersFromDigitGroups(deid, positions, charges, dgs)
}

//TODO: create here several galo.Fitter types with various Integrate approximations
// type testFitter struct {
// 	name   string
// 	fitter *galo.Fitter
// }
//
// var (
// 	fitters []testFitter
// )
//
// func init() {
// 	fitters = []testFitter{
// 		{"ref", newFitter(100)},
// 	}
// }

func newIntegrator(f func(x1, x2 float64) float64) galo.IntegrateFunc {
	return func(x1, y1, x2, y2 float64) float64 {
		return 4.0 * f(x1, x2) * f(y1, y2)
	}
}

func newFitter(deid mapping.DEID) *galo.Fitter {
	return &galo.Fitter{Integrator: mathieson.NewChargeIntegrator(deid, mathieson.IntegrateImplDefault), Method: &optimize.NelderMead{}}
}

func newFastFitter(deid mapping.DEID) *galo.Fitter {
	return &galo.Fitter{Integrator: mathieson.NewChargeIntegrator(deid, mathieson.IntegrateImplAtanEq9TanhApprox1), Method: &optimize.NelderMead{}}
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
		clusters = append(clusters, galo.Cluster{Pre: galo.PreCluster{DigitGroup: galo.DigitGroup{RefTime: 0, Digits: digits}}, Pos: clu.Pos, Q: clu.Q})
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
		// 	fitter := newFastFitter(declu.DeID)
		fitter := newFitter(declu.DeID)
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

	ci := mathieson.NewChargeIntegrator(declu.DeID, mathieson.IntegrateImplDefault)
	for _, m := range methods {
		b.Run(m.name, func(b *testing.B) {
			fitter := galo.Fitter{Integrator: ci, Method: m.method}
			for i := 0; i < b.N; i++ {
				_, _ = fitter.Position(&declu, 0)
			}
		})
	}
}

func generateTestPoints(deid mapping.DEID, N int) []testDEHit {
	var th []testDEHit
	seg := galo.SegCache.Segmentation(deid)
	box := mapping.ComputeSegmentationBBox(seg)
	i := 0
	for i < N {
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
		i++
	}
	return th
}

func generateClusters(deid int, n int) []galo.DEClusters {
	testpoints := generateTestPoints(mapping.DEID(deid), n)
	minRelCharge := 1E-4
	cs := mathieson.NewChargeSpreader(mapping.DEID(deid), minRelCharge)
	var clusters []galo.DEClusters
	for _, tp := range testpoints {
		clusters = append(clusters, createClustersFromChargeSpreader(tp, cs))
	}
	return clusters
}

func TestGenerateTestPoints(t *testing.T) {
	deid := 100
	N := 100000
	testpoints := generateTestPoints(mapping.DEID(deid), N)
	if len(testpoints) != N {
		t.Errorf("Wanted %d testpoints. Got %d\n", N, len(testpoints))
	}

}
func TestGenerateClusters(t *testing.T) {
	N := 10000
	clusters := generateClusters(100, N)
	if len(clusters) < 800 {
		t.Errorf("Wanted %d clusters. Got %d\n", N, len(clusters))
	}
}

func TestMathiesonFitterApproximations(t *testing.T) {
	clusters := generateClusters(100, 1000)
	for _, approx := range matintapproximations {
		h := hbook.NewH1D(128, 0, 50)
		fitter := &galo.Fitter{Integrator: mathieson.NewChargeIntegrator(100, approx), Method: &optimize.NelderMead{}}
		t.Run(approx.String(), func(t *testing.T) {
			for _, declu := range clusters {
				res := 1E4 * galo.DEClusterResidual(&declu, 0, fitter)
				h.Fill(res, 1.0)
			}
		})
		p := hplot.New()
		h.Scale(1 / h.Integral())
		hh := hplot.NewH1D(h)
		p.Add(hh)
		galo.SavePlot(p, "TestMathiesonIntegrateApprox", approx.String())
	}
}

func BenchmarkMathiesonFitterApproximations(b *testing.B) {
	clusters := generateClusters(100, 100)
	for _, approx := range matintapproximations {
		fitter := &galo.Fitter{Integrator: mathieson.NewChargeIntegrator(100, approx), Method: &optimize.NelderMead{}}
		b.Run(approx.String(), func(b *testing.B) {
			var r float64
			for _, declu := range clusters {
				for i := 0; i < b.N; i++ {
					r = galo.DEClusterResidual(&declu, 0, fitter)
				}
			}
			result = r
		})
	}
}
