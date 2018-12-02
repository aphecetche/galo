package galo

import "fmt"

type Digit struct {
	ID int     // digit id is the corresponding pad uid
	Q  float64 //TODO: should take only 10 bits as the original ADC value
}

type Digits []Digit

// // DigitReader wraps the basic Read (digits) method.
// type DigitReader interface {
// 	Read(digits Digits) (n int, err error)
// }
//
// // DigitWriter wraps the basic Write (digits) method.
// type DigitWriter interface {
// 	Write(digits Digits) (n int, err error)
// }

func (d Digit) String() string {
	return fmt.Sprintf("ID %6d Q %7.3f", d.ID, d.Q)
}
