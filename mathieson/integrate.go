package mathieson

import (
	"math"
)

type Integrate1DFunc func(float64, float64) float64

func Integrator1D(pitch, k3 float64, impl IntegrateImpl) Integrate1DFunc {
	atanfunc, tanhfunc := getApproximations(impl)
	sk3 := math.Sqrt(k3)
	k2 := K2FromK3(k3)
	k4 := K4FromK3(k3)
	inversePitch := 1.0 / pitch
	return func(x1, x2 float64) float64 {
		return IntegrateDefault(x1, x2, inversePitch,
			sk3, k2, k4,
			atanfunc, tanhfunc)
	}
}

func IntegrateDefault(x1, x2 float64,
	inversePitch, sk3, k2, k4 float64, atan, tanh func(float64) float64) float64 {
	lambda1 := x1 * inversePitch
	lambda2 := x2 * inversePitch
	u1 := sk3 * tanh(k2*lambda1)
	u2 := sk3 * tanh(k2*lambda2)
	return k4 * (atan(u1) - atan(u2))
}
