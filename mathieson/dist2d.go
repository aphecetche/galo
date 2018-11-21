package mathieson

type Dist2D struct {
	X, Y Dist1D
}

// NewDist2D creates a 2D Dist function of given pitch and
// given K3 parameters (one for each direction).
func NewDist2D(pitch, k3x, k3y float64) *Dist2D {
	return &Dist2D{X: *NewDist1D(pitch, k3x),
		Y: *NewDist1D(pitch, k3y)}
}

// Integral computes the 2D integral of the Dist over the area (x1,y1)->(x2,y2).
func (m *Dist2D) Integral(x1, x2, y1, y2 float64) float64 {
	return 4.0 * m.X.Integral(x1, x2) * m.Y.Integral(y1, y2)
}
