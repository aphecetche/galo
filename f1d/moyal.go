package f1d

import "math"

func Moyal(x, cst, mu, sigma float64) float64 {
	v := (x - mu) / sigma
	return cst * math.Exp(-0.5*v-0.5*math.Exp(-v))
}
