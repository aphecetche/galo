package mathieson

import (
	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/mapping"
)

func GenerateDigitGroup(deid mapping.DEID, x, y, charge float64) galo.DigitGroup {
	minRelCharge := 1E-3
	cs := NewMathiesonChargeSpreader(deid, minRelCharge)
	digits := cs.SpreadCharge(charge, x, y)
	return galo.DigitGroup{RefTime: 0, Digits: digits}
}
