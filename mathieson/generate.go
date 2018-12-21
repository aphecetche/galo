package mathieson

import (
	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/mapping"
)

func GenerateDigitGroup(deid mapping.DEID, x, y, charge float64, minRelCharge float64) galo.DigitGroup {
	cs := NewChargeSpreader(deid, minRelCharge)
	digits := cs.SpreadCharge(charge, x, y)
	return galo.DigitGroup{RefTime: 0, Digits: digits}
}
