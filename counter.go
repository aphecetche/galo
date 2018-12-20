package galo

import (
	"fmt"
	"strconv"
)

type CounterCollection struct {
	m map[string]int
}

func NewCounterCollection() *CounterCollection {
	return &CounterCollection{
		m: make(map[string]int),
	}
}

func (cc *CounterCollection) Incr(label string) {
	cc.m[label]++
}

func (cc *CounterCollection) String() string {
	s := ""
	l := 0
	for k := range cc.m {
		if len(k) > l {
			l = len(k)
		}
	}
	format := "#%" + strconv.Itoa(l) + "s=%6d"
	for k, v := range cc.m {
		s += fmt.Sprintf(format, k, v)
		s += "\n"
	}
	return s
}
