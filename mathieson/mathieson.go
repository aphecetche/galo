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
}

// NewMathieson1D creates a 1D Dist function of given pitch
// and given K3 parameter.
func NewMathieson1D(pitch, k3 float64) *Mathieson1D {
	if pitch < 1E-9 {
		panic("pitch too small")
	}
	k2 := K2FromK3(k3)
	inversePitch := 1.0 / pitch
	sk3 := math.Sqrt(k3)
	return &Mathieson1D{inversePitch: inversePitch,
		sk3: sk3,
		k2:  k2,
	}
}

func (m Mathieson1D) String() string {
	return fmt.Sprintf("pitch=%7.2f k2=%7.2f sk3=%7.2f k4=%7.2f",
		1.0/m.inversePitch, m.k2, m.sk3, K4FromK3(m.K3()))
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

func K2FromK3(k3 float64) float64 {
	sk3 := math.Sqrt(k3)
	return (math.Pi / 2.0) * (1.0 - 0.5*sk3)
}

func K4FromK3(k3 float64) float64 {
	k2 := K2FromK3(k3)
	sk3 := math.Sqrt(k3)
	c1 := k2 * sk3 / 4.0 / math.Atan(sk3)
	return c1 / k2 / sk3
}
