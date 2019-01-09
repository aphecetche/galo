package mathieson

import "math"

type IntegrateImpl int

const (
	Unknown IntegrateImpl = iota
	IntegrateImplDefault
	IntegrateImplTanhApprox1
	IntegrateImplTanhApprox2
	IntegrateImplTanhApprox3
	IntegrateImplAtanEq5
	IntegrateImplAtanEq7
	IntegrateImplAtanEq9
	IntegrateImplAtanEq11
	IntegrateImplAtanEq5TanhApprox1
	IntegrateImplAtanEq7TanhApprox1
	IntegrateImplAtanEq9TanhApprox1
	IntegrateImplAtanEq11TanhApprox1
)

func (impl IntegrateImpl) String() string {
	switch impl {
	case IntegrateImplDefault:
		return "Default"
	case IntegrateImplTanhApprox1:
		return "TanhApprox1"
	case IntegrateImplTanhApprox2:
		return "TanhApprox2"
	case IntegrateImplTanhApprox3:
		return "TanhApprox3"
	case IntegrateImplAtanEq5:
		return "AtanEq5"
	case IntegrateImplAtanEq7:
		return "AtanEq7"
	case IntegrateImplAtanEq9:
		return "AtanEq9"
	case IntegrateImplAtanEq11:
		return "AtanEq11"
	case IntegrateImplAtanEq5TanhApprox1:
		return "AtanEq5TanhApprox1"
	case IntegrateImplAtanEq7TanhApprox1:
		return "AtanEq7TanhApprox1"
	case IntegrateImplAtanEq9TanhApprox1:
		return "AtanEq9TanhApprox1"
	case IntegrateImplAtanEq11TanhApprox1:
		return "AtanEq11TanhApprox1"
	default:
		return "Unknown"
	}
}

func getApproximations(impl IntegrateImpl) (func(float64) float64, func(float64) float64) {
	if impl == IntegrateImplDefault {
		return math.Atan, math.Tanh
	}
	if impl == IntegrateImplTanhApprox1 {
		return math.Atan, TanhApprox1
	}
	if impl == IntegrateImplTanhApprox2 {
		return math.Atan, TanhApprox2
	}
	if impl == IntegrateImplTanhApprox3 {
		return math.Atan, TanhApprox3
	}
	if impl == IntegrateImplAtanEq5 {
		return AtanEq5, math.Tanh
	}
	if impl == IntegrateImplAtanEq7 {
		return AtanEq7, math.Tanh
	}
	if impl == IntegrateImplAtanEq9 {
		return AtanEq9, math.Tanh
	}
	if impl == IntegrateImplAtanEq11 {
		return AtanEq11, math.Tanh
	}
	if impl == IntegrateImplAtanEq5TanhApprox1 {
		return AtanEq5, TanhApprox1
	}
	if impl == IntegrateImplAtanEq7TanhApprox1 {
		return AtanEq7, TanhApprox1
	}
	if impl == IntegrateImplAtanEq9TanhApprox1 {
		return AtanEq9, TanhApprox1
	}
	if impl == IntegrateImplAtanEq11TanhApprox1 {
		return AtanEq11, TanhApprox1
	}
	panic("should not get there")
	return nil, nil
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
