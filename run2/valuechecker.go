package run2

import "fmt"

type ValueChecker interface {
	IsBad() bool
	IsReallyBad() bool
	IsExtremelyBad() bool
}

type ValueCheck struct {
	want float64
	got  float64
	msg  string
}

func newValueCheck(msg string, a, b float64) *ValueCheck {
	return &ValueCheck{want: a, got: b, msg: msg}
}

func (v *ValueCheck) Ratio() float64 {
	return v.want / v.got
}

func (v *ValueCheck) IsBad() bool {
	return v.Ratio() < 0.7 || v.Ratio() > 1.5
}

func (v *ValueCheck) IsReallyBad() bool {
	return v.Ratio() < 0.5 || v.Ratio() > 2.0
}

func (v *ValueCheck) IsExtremelyBad() bool {
	return v.Ratio() < 0.25 || v.Ratio() > 100.0
}

func (v *ValueCheck) String() string {
	ratio := v.Ratio()
	var ratioMsg string
	if v.IsBad() {
		ratioMsg = "!"
	}
	if v.IsReallyBad() {
		ratioMsg = "!!"
	}
	if v.IsExtremelyBad() {
		ratioMsg = "!!!"
	}
	return fmt.Sprintf("%s orig=%7.2f recalc=%7.2f ratio=%10.2f %s", v.msg, v.want, v.got, ratio, ratioMsg)
}
