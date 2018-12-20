package galo

import (
	"fmt"
	"strings"
)

type PreCluster struct {
	DigitGroup
}

func (pre PreCluster) NofPads() int {
	return len(pre.Digits)
}

func (pre PreCluster) Charge() float64 {
	var q float64
	for _, d := range pre.Digits {
		q += d.Q
	}
	return q
}

func (pre PreCluster) String() string {
	var s string
	s += fmt.Sprintf("Q=%7.3f %3d pads\n", pre.Charge(), len(pre.Digits))
	for _, d := range pre.Digits {
		s += fmt.Sprintf("%s%v\n", strings.Repeat(" ", 10), d)
	}
	return s
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
