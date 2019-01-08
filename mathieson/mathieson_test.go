package mathieson_test

import (
	"log"
	"math"
	"testing"

	"github.com/aphecetche/galo/mathieson"
	"github.com/gonum/floats"
)

// test of Mathieson function
//
// input of the test should get a list of {center pad, neighbouring list of pads} to integrate over

func TestMathieson(t *testing.T) {

	m := mathieson.St1

	var x1 float64 = 0
	var y1 float64 = 0
	var x2 float64 = 1
	var y2 float64 = 2

	v := m.Integrate(x1, y1, x2, y2)
	expected := 0.12498849 * 2.0
	if !floats.EqualWithinAbs(v, expected, 1E-6) {
		t.Errorf("Wanted %7.2f Got %7.2f\n", expected, v)
	}
}

type approxFunc struct {
	name string
	f    func(float64) float64
}

type mathiesonIntegrateFunc struct {
	name string
	f    func(float64, float64) float64
}

var (
	result               float64
	tanhapproximations   []approxFunc
	atanapproximations   []approxFunc
	matintapproximations []mathiesonIntegrateFunc
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
	matintapproximations = []mathiesonIntegrateFunc{
		{"ref", integrateRef},
		{"approx0", integrateApprox0},
		{"tanh1atan", integrateTanh1Atan},
		{"tanh1ataneq9", integrateTanh1AtanEq9},
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

const (
	inversePitch float64 = 4.761904761904762
	k2           float64 = 1.0210176124166828
	k4           float64 = 0.40934889717686523
	sk3          float64 = 0.7
)

func integrate(x1, x2 float64, tanh, atan func(float64) float64) float64 {
	lambda1 := x1 * inversePitch
	lambda2 := x2 * inversePitch
	u1 := sk3 * tanh(k2*lambda1)
	u2 := sk3 * tanh(k2*lambda2)
	return k4 * (atan(u1) - atan(u2))
}

func integrateRef(x1, x2 float64) float64 {
	return integrate(x1, x2, math.Tanh, math.Atan)
}

func integrateApprox0(x1, x2 float64) float64 {
	return integrateBis(x1, x2, math.Tanh, math.Atan)
}

func integrateBis(x1, x2 float64, tanh, atan func(float64) float64) float64 {
	lambda1 := x1 * inversePitch
	lambda2 := x2 * inversePitch
	u1 := sk3 * tanh(k2*lambda1)
	u2 := sk3 * tanh(k2*lambda2)
	r := k4 * (atan(u1) - atan(u2))
	u12 := u1 * u2
	k := 0.0
	if u12 > 1.0 {
		if u1 > 0.0 {
			k = math.Pi
		} else {
			k = -math.Pi
		}
	}
	a := k4 * (atan((u1-u2)/(1.0+u12)) + k)
	if !floats.EqualWithinAbs(a, r, 1E-3) {
		log.Fatalf("a=%g r=%g u12=%g x1=%g x2=%g", a, r, u12, x1, x2)
	}
	return a
}

func integrateTanh1Atan(x1, x2 float64) float64 {
	return integrate(x1, x2, mathieson.TanhApprox1, math.Atan)
}

func integrateTanh1AtanEq9(x1, x2 float64) float64 {
	return integrate(x1, x2, mathieson.TanhApprox1, mathieson.AtanEq9)
}

func TestMathiesonIntegrateApprox(t *testing.T) {
	for _, approx := range matintapproximations {
		t.Run(approx.name, func(t *testing.T) {
			for x := 0.0; x < 20.0; x += 0.1 {
				x1 := x
				x2 := x + 1.0
				a := approx.f(x1, x2)
				expected := integrateRef(x1, x2)
				if !floats.EqualWithinAbs(a, expected, 1E-3) {
					t.Errorf("Wrong approx for %s(%g)=%g. Expected %g", approx.name,
						x, a, expected)
					break
				}
			}
		})
	}
}

func BenchmarkMathiesonIntegrateApprox(b *testing.B) {
	for _, approx := range matintapproximations {
		b.Run(approx.name, func(b *testing.B) {
			var r float64
			for n := 0; n < b.N; n++ {
				for x := 0.0; x < 20.0; x += 0.1 {
					x1 := x
					x2 := x + 1.0
					r = approx.f(x1, x2)
				}
			}
			result = r
		})
	}
}
