package mathieson_test

import (
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
