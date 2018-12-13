package galo

import "gonum.org/v1/gonum/floats"

// SameCluster returns true if the two clusters :
// - have the same precluster
// - have close enough positions
func SameCluster(ca, cb Cluster) bool {
	pa := ca.Pre
	pb := cb.Pre
	if !SamePreCluster(pa, pb) {
		return false
	}
	const tol = 1E-6
	return floats.EqualWithinAbs(float64(ca.Pos.X), float64(cb.Pos.X), tol) &&
		floats.EqualWithinAbs(float64(ca.Pos.Y), float64(cb.Pos.Y), tol)
}

// SameDigitLocation returns true if the digits are the same pad
func SameDigitLocation(da, db Digit) bool {
	return da.ID == db.ID
}

// SamePreCluster returns true if both preclusters have :
// - the same digits
// - in the same order
func SamePreCluster(a, b PreCluster) bool {
	if a.NofPads() != b.NofPads() {
		return false
	}
	for i := 0; i < a.NofPads(); i++ {
		da := a.Digits[i]
		db := b.Digits[i]
		if !SameDigitLocation(da, db) {
			return false
		}
	}
	return true
}

// ShareDigits returns true if both precluster have at least
// one digit in common
func ShareDigits(a, b PreCluster) bool {
	for i := 0; i < a.NofPads(); i++ {
		da := a.Digits[i]
		for j := 0; j < b.NofPads(); j++ {
			db := b.Digits[j]
			if SameDigitLocation(da, db) {
				return true
			}
		}
	}
	return false
}
