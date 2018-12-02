package galo

import "math"

func Moyal(x, cst, mu, sigma float64) float64 {
	v := (x - mu) / sigma
	return cst * math.Exp(-0.5*v-0.5*math.Exp(-v))
}

var (
	p1 = [5]float64{0.4259894875, -0.1249762550, 0.03984243700, -0.006298287635, 0.001511162253}
	q1 = [5]float64{1.0, -0.3388260629, 0.09594393323, -0.01608042283, 0.003778942063}
	p2 = [5]float64{0.1788541609, 0.1173957403, 0.01488850518, -0.001394989411, 0.0001283617211}
	q2 = [5]float64{1.0, 0.7428795082, 0.3153932961, 0.06694219548, 0.008790609714}
	p3 = [5]float64{0.1788544503, 0.09359161662, 0.006325387654, 0.00006611667319, -0.000002031049101}
	q3 = [5]float64{1.0, 0.6097809921, 0.2560616665, 0.04746722384, 0.006957301675}
	p4 = [5]float64{0.9874054407, 118.6723273, 849.2794360, -743.7792444, 427.0262186}
	q4 = [5]float64{1.0, 106.8615961, 337.6496214, 2016.712389, 1597.063511}
	p5 = [5]float64{1.003675074, 167.5702434, 4789.711289, 21217.86767, -22324.94910}
	q5 = [5]float64{1.0, 156.9424537, 3745.310488, 9834.698876, 66924.28357}
	p6 = [5]float64{1.000827619, 664.9143136, 62972.92665, 475554.6998, -5743609.109}
	q6 = [5]float64{1.0, 651.4101098, 56974.73333, 165917.4725, -2815759.939}
	a1 = [3]float64{0.04166666667, -0.01996527778, 0.02709538966}
	a2 = [2]float64{-1.845568670, -4.284640743}
)

func Landau(x, mu, sigma float64) float64 {
	if sigma <= 0 {
		return 0
	}
	xi := 1.0
	x0 := 0.0
	den := landau_pdf((x-mu)/sigma, xi, x0)
	return den
}

// LANDAU pdf : algorithm from CERNLIB G110 denlan
// same algorithm is used in GSL
func landau_pdf(x, xi, x0 float64) float64 {
	if xi <= 0 {
		return 0
	}
	v := (x - x0) / xi
	var u, ue, us, denlan float64
	if v < -5.5 {
		u = math.Exp(v + 1.0)
		if u < 1e-10 {
			return 0.0
		}
		ue = math.Exp(-1 / u)
		us = math.Sqrt(u)
		denlan = 0.3989422803 * (ue / us) * (1 + (a1[0]+(a1[1]+a1[2]*u)*u)*u)
	} else if v < -1 {
		u = math.Exp(-v - 1)
		denlan = math.Exp(-u) * math.Sqrt(u) *
			(p1[0] + (p1[1]+(p1[2]+(p1[3]+p1[4]*v)*v)*v)*v) /
			(q1[0] + (q1[1]+(q1[2]+(q1[3]+q1[4]*v)*v)*v)*v)
	} else if v < 1 {
		denlan = (p2[0] + (p2[1]+(p2[2]+(p2[3]+p2[4]*v)*v)*v)*v) /
			(q2[0] + (q2[1]+(q2[2]+(q2[3]+q2[4]*v)*v)*v)*v)
	} else if v < 5 {
		denlan = (p3[0] + (p3[1]+(p3[2]+(p3[3]+p3[4]*v)*v)*v)*v) /
			(q3[0] + (q3[1]+(q3[2]+(q3[3]+q3[4]*v)*v)*v)*v)
	} else if v < 12 {
		u = 1 / v
		denlan = u * u * (p4[0] + (p4[1]+(p4[2]+(p4[3]+p4[4]*u)*u)*u)*u) /
			(q4[0] + (q4[1]+(q4[2]+(q4[3]+q4[4]*u)*u)*u)*u)
	} else if v < 50 {
		u = 1 / v
		denlan = u * u * (p5[0] + (p5[1]+(p5[2]+(p5[3]+p5[4]*u)*u)*u)*u) /
			(q5[0] + (q5[1]+(q5[2]+(q5[3]+q5[4]*u)*u)*u)*u)
	} else if v < 300 {
		u = 1 / v
		denlan = u * u * (p6[0] + (p6[1]+(p6[2]+(p6[3]+p6[4]*u)*u)*u)*u) /
			(q6[0] + (q6[1]+(q6[2]+(q6[3]+q6[4]*u)*u)*u)*u)
	} else {
		u = 1 / (v - v*math.Log(v)/(v+1))
		denlan = u * u * (1 + (a2[0]+a2[1]*u)*u)
	}
	return denlan / xi
}

func Levy(x, mu, sigma float64) float64 {
	return math.Exp(-sigma/(2*(x-mu))) / math.Pow((x-mu), 1.5)
}

func Gaus(x, cst, mu, sigma float64) float64 {
	v := (x - mu) / sigma
	return cst * math.Exp(-0.5*v*v)
}
