package mathieson_test

import (
	"math"
	"testing"

	"github.com/aphecetche/galo/mathieson"
	"github.com/gonum/floats"
)

// test of Mathieson function
//
// input of the test should get a list of {center pad, neighbouring list of pads} to integrate over

func TestMathieson(t *testing.T) {

	var x1 float64 = 0
	var y1 float64 = 0
	var x2 float64 = 1
	var y2 float64 = 2

	ci := mathieson.NewChargeIntegrator(100, mathieson.IntegrateImplDefault)
	v := ci.Integrate(x1, y1, x2, y2)
	expected := 0.12498849 * 2.0
	if !floats.EqualWithinAbs(v, expected, 1E-6) {
		t.Errorf("Wanted %7.2f Got %7.2f\n", expected, v)
	}
}

type approxFunc struct {
	name string
	f    func(float64) float64
}

var (
	result               float64
	tanhapproximations   []approxFunc
	atanapproximations   []approxFunc
	matintapproximations []mathieson.IntegrateImpl
)

func init() {
	tanhapproximations = []approxFunc{
		{"math", math.Tanh},
		{"approx1", mathieson.TanhApprox1},
		{"approx2", mathieson.TanhApprox2},
		{"approx3", mathieson.TanhApprox3},
	}
	atanapproximations = []approxFunc{
		{"math", math.Atan},
		{"eq5", mathieson.AtanEq5},
		{"eq7", mathieson.AtanEq7},
		{"eq9", mathieson.AtanEq9},
		{"eq11", mathieson.AtanEq11},
	}
	matintapproximations = []mathieson.IntegrateImpl{
		mathieson.IntegrateImplDefault,
		mathieson.IntegrateImplTanhApprox1,
		mathieson.IntegrateImplAtanEq5,
		mathieson.IntegrateImplAtanEq7,
		mathieson.IntegrateImplAtanEq9,
		mathieson.IntegrateImplAtanEq11,
	}
}

func TestTanhApprox(t *testing.T) {
	for _, approx := range tanhapproximations {
		t.Run(approx.name, func(t *testing.T) {
			for x := -20.0; x < 20.0; x += 0.01 {
				a := approx.f(x)
				expected := math.Tanh(x)
				if !floats.EqualWithinAbs(a, expected, 1E-3) {
					t.Errorf("Wrong approx for %s(%g)=%g. Expected %g", approx.name,
						x, a, expected)
					break
				}
			}
		})
	}
}

func TestAtanApprox(t *testing.T) {
	for _, approx := range atanapproximations {
		t.Run(approx.name, func(t *testing.T) {
			for x := -0.8; x < 0.8; x += 0.01 {
				a := approx.f(x)
				expected := math.Atan(x)
				if !floats.EqualWithinAbs(a, expected, 2E-2) {
					t.Errorf("Wrong approx for %s(%g)=%g. Expected %g", approx.name,
						x, a, expected)
					break
				}
			}
		})
	}
}

func BenchmarkTanhApprox(b *testing.B) {
	for _, approx := range tanhapproximations {
		b.Run(approx.name, func(b *testing.B) {
			var r float64
			for n := 0; n < b.N; n++ {
				for x := -20.0; x < 20.0; x += 0.01 {
					r = approx.f(x)
				}
			}
			result = r
		})
	}
}

func BenchmarkAtanApprox(b *testing.B) {
	for _, approx := range atanapproximations {
		b.Run(approx.name, func(b *testing.B) {
			var r float64
			for n := 0; n < b.N; n++ {
				for x := -0.8; x < 0.8; x += 0.01 {
					r = approx.f(x)
				}
			}
			result = r
		})
	}
}

func TestMathiesonIntegrate1D(t *testing.T) {
	pitch := mathieson.St1.Pitch
	k3 := mathieson.St1.K3x
	def := mathieson.Integrator1D(pitch, k3, mathieson.IntegrateImplDefault)
	for _, approx := range matintapproximations {
		impl := mathieson.Integrator1D(pitch, k3, approx)
		t.Run(approx.String(), func(t *testing.T) {
			for x := 0.0; x < 20.0; x += 0.01 {
				x1 := x
				x2 := x + 1.0
				a := impl(x1, x2)
				expected := def(x1, x2)
				if !floats.EqualWithinAbs(a, expected, 1E-2) {
					t.Errorf("Wrong approx for %s(%g)=%g. Expected %g", approx.String(),
						x, a, expected)
					break
				}
			}
		})
	}
}

func BenchmarkMathiesonIntegrate1D(b *testing.B) {
	pitch := mathieson.St1.Pitch
	k3 := mathieson.St1.K3x
	for _, approx := range matintapproximations {
		impl := mathieson.Integrator1D(pitch, k3, approx)
		b.Run(approx.String(), func(b *testing.B) {
			var r float64
			for n := 0; n < b.N; n++ {
				for x := 0.0; x < 20.0; x += 0.01 {
					x1 := x
					x2 := x + 1.0
					r = impl(x1, x2)
				}
			}
			result = r
		})
	}
}
