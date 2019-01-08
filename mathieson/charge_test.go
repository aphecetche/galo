package mathieson_test

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/mathieson"
)

func sortByID(digits []galo.Digit) {
	sort.Slice(digits, func(i, j int) bool {
		return digits[i].ID < digits[j].ID
	})
}

func dumpDigits(digits []galo.Digit) {
	for _, d := range digits {
		fmt.Printf("%[1]v (%[1]T)\n", d)
	}
}

func dumpRelDigits(digits []galo.Digit) {
	for _, d := range digits {
		fmt.Printf("ID %6d QD ", d.ID)
		if d.Q > 0 {
			fmt.Printf("%5.2f %%", d.Q)
		} else {
			fmt.Printf(" [ ] ")
		}
		fmt.Printf("\n")
	}
}

func chargeSum(digits []galo.Digit) float64 {
	q := 0.0
	for _, d := range digits {
		q += d.Q
	}
	return q
}

func relativeDiffPercent(a, b float64) float64 {
	return math.Abs((a-b)/b) * 100.0
}

func relative(digits []galo.Digit, ref []galo.Digit) ([]galo.Digit, float64) {
	var diff []galo.Digit
	var refi int
	for _, d := range digits {
		refi = -1
		for i, r := range ref {
			if d.ID == r.ID {
				refi = i
				break
			}
		}
		if refi >= 0 {
			delta := relativeDiffPercent(d.Q, ref[refi].Q)
			diff = append(diff, galo.Digit{ID: d.ID, Q: delta})
		} else {
			// new pad
			diff = append(diff, galo.Digit{ID: d.ID, Q: -1.0})
		}
	}
	return diff, relativeDiffPercent(chargeSum(digits), chargeSum(ref))
}

func TestMinChargeEffect(t *testing.T) {
	q0 := 100.0
	x0 := 24.0
	y0 := 24.0

	for _, minRelCharge := range []float64{1E-2, 1E-3, 1E-4, 1E-5} {
		cs := mathieson.NewChargeSpreader(100, minRelCharge)
		ref := cs.SpreadCharge(q0, x0, y0)
		q := chargeSum(ref)
		fmt.Printf("minRelCharge %e Qtot %7.2f Qreldiffpercent %5.2f\n", minRelCharge, q, relativeDiffPercent(q, q0))
	}
}

func TestShiftPosition(t *testing.T) {

	minRelCharge := 1E-3

	cs := mathieson.NewChargeSpreader(100, minRelCharge)

	q0 := 100.0
	x0 := 24.0
	y0 := 24.0

	ref := cs.SpreadCharge(q0, x0, y0)

	fmt.Println("Qrel vs target=", relativeDiffPercent(chargeSum(ref), q0))
	sortByID(ref)
	dumpDigits(ref)

	for _, d := range []float64{1.0, 10.0, 20., 50.0, 100.0, 200.0} {
		x := x0 + d*1E-4
		y := y0
		fmt.Printf("delta %7.2f microns\n", d)
		digits := cs.SpreadCharge(q0, x, y)
		sortByID(digits)
		rd, q := relative(digits, ref)
		dumpRelDigits(rd)
		fmt.Println("DeltaQ(relative to ref)=", q)
	}
}
