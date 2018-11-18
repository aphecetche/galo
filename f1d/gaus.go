package f1d

import "math"

func Gaus(x, cst, mu, sigma float64) float64 {
	v := (x - mu) / sigma
	return cst * math.Exp(-0.5*v*v)
}
