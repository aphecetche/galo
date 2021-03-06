package galo

import (
	"fmt"

	"github.com/aphecetche/pigiron/mapping"
)

type Digit struct {
	ID      mapping.PadUID // digit id is the corresponding pad uid (relative to a detection element ID)
	Q       float64        //TODO: should take only 10 bits as the original ADC value
	Toffset byte           // time offset relative to the group ref time
}

type DigitGroup struct {
	RefTime int // reference timestamp for the group digits
	Digits  []Digit
}

func (d Digit) String() string {
	return fmt.Sprintf("ID %6d Q %7.3f", d.ID, d.Q)
}

// SameDigitLocation returns true if the digits are the same pad
func SameDigitLocation(da, db Digit) bool {
	return da.ID == db.ID
}
