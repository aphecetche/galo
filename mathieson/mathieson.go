package mathieson

import (
	"fmt"
	"math"
)

// Mathieson1D is a function that describe the charge distribution
// on cathodes of multiwire chambers.
//
// See e.g. Mathieson, E. (1988). Cathode Charge Distributions in Multiwire
// Chambers. 4: Empirical Formula for Small Anode - Cathode Separation.
// Nucl. Instrum. Meth., A270(2-3), 602–603.
// http://doi.org/10.1016/0168-9002(88)90736-X
type Mathieson1D struct {
	inversePitch float64
	sk3          float64
	k2           float64
	k4           float64
	atanfunc     func(x float64) float64
	tanhfunc     func(x float64) float64
}

// Atan approximations below come from the paper
//
// Efficient approximations for the arctangent function
// by S. Rajan ; Sichun Wang ; R. Inkol ; A. Joyal
// IEEE Signal Processing Magazine ( Volume: 23 , Issue: 3 , May 2006 )
// DOI: 10.1109/MSP.2006.1628884
//
// and are labelled using the equation numbers in that paper

const pi_over_4 float64 = math.Pi / 4.0

func AtanEq5(x float64) float64 {
	return pi_over_4*x + 0.285*x*(1-math.Abs(x))
}

func AtanEq7(x float64) float64 {
	return pi_over_4*x + 0.273*x*(1-math.Abs(x))
}

func AtanEq9(x float64) float64 {
	return pi_over_4*x - x*(math.Abs(x)-1)*(0.2447+0.0663*math.Abs(x))
}

func AtanEq11(x float64) float64 {
	return x / (1 + 0.28125*x*x)
}

func TanhApprox3(x float64) float64 {
	x2 := x * x
	if x2 > 25 {
		if x > 0 {
			return 1.0
		}
		return -1.0
	}
	// return x / (1 + x2/(3+x2/(5+x2/(7+x2/(9+x2/(11+x2/(13)))))))
	return x / (1 + x2/(3+x2/(5+x2/(7+x2/(9+x2/(11))))))
}

func TanhApprox2(x float64) float64 {
	// from https://www.math.utah.edu/~beebe/software/ieee/tanh.pdf
	// note we are using double precision here.
	// TODO see if single precision would be enough ?
	// (note that the constant terms are different, see the pdf above)
	// does not seem to help a lot wrt to math.Tanh (see corresponding benchmark)
	// (about half of the time spent in math.Exp anyway)

	const (
		xsmall  float64 = 1.29047841397589243466E-08
		xmedium float64 = 0.54930614433405484570
		xlarge  float64 = 19.06154746539849600897
		p0      float64 = -0.16134119023996228053E+04
		p1      float64 = -0.99225929672236083313E+02
		p2      float64 = -0.96437492777225469787
		q0      float64 = 0.48402357071988688686E+04
		q1      float64 = 0.22337720718962312926E+04
		q2      float64 = 0.11274474380534949335E+03
		q3      float64 = 1.0
	)

	xp := math.Abs(x)

	r := xp

	if xp > xsmall && xp < xmedium {
		g := xp * xp
		p := g*(p2*g+p1)*g + p0
		q := ((g+q2)*g+q1)*g + q0
		r = xp + g*xp*p/q
	} else if xp >= xmedium && xp < xlarge {
		t := 0.5 - 1/(1+math.Exp(2*xp))
		r = t + t
	} else if xp > xlarge {
		r = 1
	}

	if x < 0 {
		return -r
	}
	return r
}

func TanhApprox1(x float64) float64 {
	// from https://varietyofsound.wordpress.com/2011/02/14/efficient-tanh-computation-using-lamberts-continued-fraction/ but with a clipper,
	// as the approx seems to "blow up" badly after 4.5-5.0
	c := 5.0
	if x > c {
		return 1.0
	}
	if x < -c {
		return -1.0
	}
	x2 := x * x
	a := (((x2+378)*x2+17325)*x2 + 135135) * x
	b := ((28*x2+3150)*x2+62370)*x2 + 135135
	return a / b
}

// NewMathieson1D creates a 1D Dist function of given pitch
// and given K3 parameter.
func NewMathieson1D(pitch, k3 float64) *Mathieson1D {
	return NewMathieson1DApprox(pitch, k3, math.Atan, math.Tanh)
}

// NewMathieson1DApprox creates a 1D Dist function of given pitch
// and given K3 parameter, with given approximation functions for Atan
// and Tanh.
func NewMathieson1DApprox(pitch, k3 float64, atanfunc, tanhfunc func(float64) float64) *Mathieson1D {
	if pitch < 1E-9 {
		panic("pitch too small")
	}
	k2, k4 := computeK2K4FromK3(k3)
	inversePitch := 1.0 / pitch
	return &Mathieson1D{inversePitch: inversePitch,
		sk3:      math.Sqrt(k3),
		k2:       k2,
		k4:       k4,
		atanfunc: atanfunc,
		tanhfunc: tanhfunc,
	}
}

func (m Mathieson1D) String() string {
	return fmt.Sprintf("pitch=%7.2f k2=%7.2f sk3=%7.2f k4=%7.2f",
		1.0/m.inversePitch, m.k2, m.sk3, m.k4)
}

// Pitch returns the pitch (distance between the wires of the cathode pad
// chamber).
func (m Mathieson1D) Pitch() float64 {
	return 1.0 / m.inversePitch
}

// K3 returns the K3 parameter of the Mathieson.
func (m Mathieson1D) K3() float64 {
	return m.sk3 * m.sk3
}

// F returns the evaluation of the Mathieson at point x.
func (m Mathieson1D) F(x float64) float64 {
	t := math.Tanh(m.k2 * x)
	t = t * t
	return (1 - t) / (1 + m.K3()*t)
}

// FWHM returns an estimate of the full width half maximum.
func (m Mathieson1D) FWHM() float64 {
	a := math.Sqrt(2 + m.K3())
	w := 4.0 * math.Atanh(1.0/a)
	w /= math.Pi * (1 - 0.5*m.sk3)
	return w
}

// Integrate computes the 1D integral of the Dist between x1 and x2.
func (m *Mathieson1D) Integrate(x1, x2 float64) float64 {
	lambda1 := x1 * m.inversePitch
	lambda2 := x2 * m.inversePitch
	u1 := m.sk3 * m.tanhfunc(m.k2*lambda1)
	u2 := m.sk3 * m.tanhfunc(m.k2*lambda2)
	return m.k4 * (m.atanfunc(u1) - m.atanfunc(u2))
}

// Mathieson2D is the product of two Mathieson1D.
type Mathieson2D struct {
	X, Y Mathieson1D
}

// NewMathieson2D creates a 2D Dist function of given pitch and
// given K3 parameters (one for each direction).
func NewMathieson2D(pitch, k3x, k3y float64) *Mathieson2D {
	return &Mathieson2D{X: *NewMathieson1D(pitch, k3x),
		Y: *NewMathieson1D(pitch, k3y)}
}

// NewMathieson2DApprox creates a 2D Dist function of given pitch and
// given K3 parameters (one for each direction), with given
// approximation functions for Atan and Tanh
func NewMathieson2DApprox(pitch, k3x, k3y float64, atanfunc, tanhfunc func(float64) float64) *Mathieson2D {
	return &Mathieson2D{X: *NewMathieson1DApprox(pitch, k3x, atanfunc, tanhfunc),
		Y: *NewMathieson1DApprox(pitch, k3y, atanfunc, tanhfunc)}
}

// Integrate computes the 2D integral of the Dist over the area (x1,y1)->(x2,y2).
func (m *Mathieson2D) Integrate(x1, y1, x2, y2 float64) float64 {
	return 4.0 * m.X.Integrate(x1, x2) * m.Y.Integrate(y1, y2)
}

func (m *Mathieson2D) String() string {
	s := fmt.Sprintf("X=%s\n", m.X.String())
	s += fmt.Sprintf("Y=%s\n", m.Y.String())
	return s
}

// ComputeK2K4FromK3 computes the K3 parameter of the Dist,
// given its K2 and K4 parameters.
func computeK2K4FromK3(k3 float64) (float64, float64) {
	sk3 := math.Sqrt(k3)
	k2 := (math.Pi / 2.0) * (1.0 - 0.5*sk3)
	c1 := k2 * sk3 / 4.0 / math.Atan(sk3)
	k4 := c1 / k2 / sk3
	return k2, k4
}
