package mathieson

import (
	"fmt"
	"testing"
)

// test of Mathieson function
//
// input of the test should get a list of {center pad, neighbouring list of pads} to integrate over

func TestMathieson(t *testing.T) {

	m := St1

	var x1 float64 = 0
	var y1 float64 = 0
	var x2 float64 = 1
	var y2 float64 = 2

	fmt.Println(m.Integral(x1, y1, x2, y2))
}
