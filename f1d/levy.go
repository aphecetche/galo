package f1d

import "math"

func Levy(x, mu, sigma float64) float64 {

	return math.Exp(-sigma/(2*(x-mu))) / math.Pow((x-mu), 1.5)

}
