package run2

import "gonum.org/v1/gonum/floats"

// SameCluster returns true if the two clusters :
// - have the same precluster
// - have close enough positions
func SameCluster(ca, cb Cluster) bool {
	pa := ca.Pre(nil)
	pb := cb.Pre(nil)
	if !SamePreCluster(*pa, *pb) {
		return false
	}
	const tol = 1E-6
	return floats.EqualWithinAbs(float64(ca.Pos(nil).X()), float64(cb.Pos(nil).X()), tol) &&
		floats.EqualWithinAbs(float64(ca.Pos(nil).Y()), float64(cb.Pos(nil).Y()), tol)
}

// SameDigitLocation returns true if the digits are the same pad
func SameDigitLocation(da, db Digit) bool {
	if da.Deid() != db.Deid() {
		return false
	}
	if da.Manuid() != db.Manuid() {
		return false
	}
	if da.Manuchannel() != db.Manuchannel() {
		return false
	}
	return true
}

// SamePreCluster returns true if both preclusters have :
// - the same digits
// - in the same order
func SamePreCluster(a, b PreCluster) bool {
	if a.DigitsLength() != b.DigitsLength() {
		return false
	}
	var da, db Digit
	for i := 0; i < a.DigitsLength(); i++ {
		a.Digits(&da, i)
		b.Digits(&db, i)
		if !SameDigitLocation(da, db) {
			return false
		}
	}
	return true
}

// ShareDigits returns true if both precluster have at least
// one digit in common
func ShareDigits(a, b PreCluster) bool {
	var da, db Digit
	for i := 0; i < a.DigitsLength(); i++ {
		a.Digits(&da, i)
		for j := 0; j < b.DigitsLength(); j++ {
			b.Digits(&db, j)
			if SameDigitLocation(da, db) {
				return true
			}
		}
	}
	return false
}
