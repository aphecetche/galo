package galo

import (
	"fmt"
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
		s += d.String() + "\n"
	}
	return s
}
