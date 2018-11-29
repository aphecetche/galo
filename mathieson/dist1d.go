package mathieson

import (
	"math"
)

type Dist1D struct {
	inversePitch float64
	sk3          float64
	k2           float64
	k4           float64
}

// Newmathieson1d creates a 1D Dist function of given pitch
// and given K3 parameter.
func NewDist1D(pitch, k3 float64) *Dist1D {
	if pitch < 1E-9 {
		panic("pitch too small")
	}
	k2, k4 := computeK2K4FromK3(k3)
	inversePitch := 1.0 / pitch
	return &Dist1D{inversePitch: inversePitch,
		sk3: math.Sqrt(k3),
		k2:  k2,
		k4:  k4}
}

func (m *Dist1D) Pitch() float64 {
	return 1.0 / m.inversePitch
}

func (m *Dist1D) K3() float64 {
	return m.sk3 * m.sk3
}

func (m *Dist1D) F(x float64) float64 {
	t := math.Tanh(m.k2 * x)
	t = t * t
	return (1 - t) / (1 + m.K3()*t)
}

func (m *Dist1D) FWHM() float64 {
	a := math.Sqrt(2 + m.K3())
	w := 4.0 * math.Atanh(1.0/a)
	w /= math.Pi * (1 - 0.5*m.sk3)
	return w
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

// Integral computes the 1D integral of the Dist between x1 and x2.
func (m *Dist1D) Integral(x1, x2 float64) float64 {
	lambda1 := x1 * m.inversePitch
	lambda2 := x2 * m.inversePitch
	u1 := m.sk3 * math.Tanh(m.k2*lambda1)
	u2 := m.sk3 * math.Tanh(m.k2*lambda2)
	return m.k4 * (math.Atan(u1) - math.Atan(u2))
}
