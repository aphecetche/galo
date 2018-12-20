package run2

import (
	"log"
)

func computeDigitCharge(dig Digit, x, y float32) float64 {
	log.Fatalf("reimplement me")
	return 1.0
	// deid := int(dig.Deid())
	// manuid := mapping.DualSampaID(dig.Manuid())
	// manuch := int(dig.Manuchannel())
	// integ := mathieson.NewChargeIntegrator(deid)
	// bending := true
	// if manuid >= 1024 {
	// 	bending = false
	// }
	// cseg := segcache.CathodeSegmentation(deid, bending)
	// paduid, err := cseg.FindPadByFEE(manuid, manuch)
	// if err != nil {
	// 	log.Fatalf("Could not find at (%v,%v)  DE %v", x, y, deid)
	// }
	// return generate.ChargeOverPad(float64(x), float64(y), integ, paduid, cseg)
}
