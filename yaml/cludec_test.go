package yaml

import (
	"fmt"
	"log"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestReadOnePixel(t *testing.T) {
	pixelYAML := `
x: 25.2
y: 74.13
dx: 0.315
dy: 0.21
`
	var pix yaPixel
	err := yaml.Unmarshal([]byte(pixelYAML), &pix)
	if err != nil {
		log.Fatalf("Error:%v", err)
	}
	if pix.X != 25.2 ||
		pix.Y != 74.13 ||
		pix.DX != 0.315 ||
		pix.DY != 0.21 {
		t.Errorf("wrong pixel read in")
	}
}

func TestReadPixelSeq(t *testing.T) {
	pixelsYAML := `
- x: 25.2
  y: 74.13
  dx: 0.315
  dy: 0.21
- x: 25.83
  y: 74.55000000000001
  dx: 0.315
  dy: 0.21
`
	var pixels []yaPixel
	err := yaml.Unmarshal([]byte(pixelsYAML), &pixels)
	if err != nil {
		log.Fatalf("Error:%v", err)
	}
	if len(pixels) != 2 {
		t.Errorf("Wrong number of pixels")
	}
	pix0 := pixels[0]
	if pix0.X != 25.2 ||
		pix0.Y != 74.13 ||
		pix0.DX != 0.315 ||
		pix0.DY != 0.21 {
		t.Errorf("wrong 1st pixel read in")
	}
	pix1 := pixels[1]
	if pix1.X != 25.83 ||
		pix1.Y != 74.55 ||
		pix1.DX != 0.315 ||
		pix1.DY != 0.21 {
		t.Errorf("wrong 2nd pixel read in")
	}
}

func TestReadStep(t *testing.T) {
	stepYAML := `
pixels:
  - x: 25.2
    y: 74.13
    dx: 0.315
    dy: 0.21
`
	var step yaStep

	err := yaml.Unmarshal([]byte(stepYAML), &step)
	if err != nil {
		log.Fatalf("Error:%v", err)
	}
	if len(step.Pixels) != 1 {
		t.Errorf("Wrong number of pixels : %v", len(step.Pixels))
	}
	pix0 := step.Pixels[0]
	if pix0.X != 25.2 ||
		pix0.Y != 74.13 ||
		pix0.DX != 0.315 ||
		pix0.DY != 0.21 {
		t.Errorf("wrong 1st pixel read in")
	}
}

type stepStruct struct {
	Steps []yaStep
}

func TestReadStepSeq(t *testing.T) {
	stepYAML := `
steps:
  - pixels:
    - x: 25.2
      y: 74.13
      dx: 0.315
      dy: 0.21
`
	var step stepStruct
	err := yaml.Unmarshal([]byte(stepYAML), &step)
	if err != nil {
		log.Fatalf("Error:%v", err)
	}
	if len(step.Steps) != 1 {
		t.Errorf("Wrong number of steps : %v", len(step.Steps))
	}
	pixels := step.Steps[0].Pixels
	if len(pixels) != 1 {
		t.Errorf("Want 1 pixels - Got %d", len(pixels))
	}
	// if pix0.X != 25.2 ||
	// 	pix0.Y != 74.13 ||
	// 	pix0.DX != 0.315 ||
	// 	pix0.DY != 0.21 {
	// 	t.Errorf("wrong 1st pixel read in")
	// }
}

func checkDigits(t *testing.T, digits []yaDigit) error {

	if len(digits) != 2 {
		return fmt.Errorf("Wanted 2 digits - Got %d", len(digits))
	}

	d0 := digits[0]

	if d0.Deid != 100 ||
		d0.Dsid != 235 ||
		d0.Dsch != 16 ||
		d0.Adc != 294 ||
		d0.Charge != 4.661163 {
		return fmt.Errorf("wrong 1st digit read in")
	}
	d1 := digits[1]

	if d1.Deid != 100 ||
		d1.Dsid != 235 ||
		d1.Dsch != 61 ||
		d1.Adc != 538 ||
		d1.Charge != 36.61433 {
		return fmt.Errorf("wrong 2nd digit read in")
	}
	return nil
}

func TestReadDigitSeq(t *testing.T) {
	digitsYAML := `
- deid: 100
  dsid: 235
  dsch: 16
  adc: 294
  charge: 4.661163
- deid: 100
  dsid: 235
  dsch: 61
  adc: 538
  charge: 36.61433
`
	var digits []yaDigit

	err := yaml.Unmarshal([]byte(digitsYAML), &digits)
	if err != nil {
		log.Fatalf("Error:%v", err)
	}

	err = checkDigits(t, digits)

	if err != nil {
		t.Errorf(err.Error())
	}

}

type preStruct struct {
	Pre yaPre
}

func TestReadPre(t *testing.T) {
	preYAML := `
pre:
  digitgroup:
    digits:
     - deid: 100
       dsid: 235
       dsch: 16
       adc: 294
       charge: 4.661163
     - deid: 100
       dsid: 235
       dsch: 61
       adc: 538
       charge: 36.61433
`
	var pre preStruct

	err := yaml.Unmarshal([]byte(preYAML), &pre)
	if err != nil {
		log.Fatalf("Error:%v", err)
	}

	if len(pre.Pre.DigitGroup.Digits) != 2 {
		t.Errorf("Want 2 digits - got %d", len(pre.Pre.DigitGroup.Digits))
	}
	err = checkDigits(t, pre.Pre.DigitGroup.Digits)

	if err != nil {
		t.Errorf(err.Error())
	}

}

type prePosChargeStruct struct {
	Pre    yaPre
	Pos    yaPos
	Charge float32
}

func TestReadPrePosCharge(t *testing.T) {
	preposchargeYAML := `
pre:
  digitgroup:
    digits:
     - deid: 100
       dsid: 235
       dsch: 16
       adc: 294
       charge: 4.661163
     - deid: 100
       dsid: 235
       dsch: 61
       adc: 538
       charge: 36.61433
pos:
  x: 25.45887
  y: 74.50821
charge: 57.69108
`
	var preposcharge prePosChargeStruct

	err := yaml.Unmarshal([]byte(preposchargeYAML), &preposcharge)
	if err != nil {
		log.Fatalf("Error:%v", err)
	}

	if len(preposcharge.Pre.DigitGroup.Digits) != 2 {
		t.Errorf("Want 2 digits - got %d", len(preposcharge.Pre.DigitGroup.Digits))
	}
	err = checkDigits(t, preposcharge.Pre.DigitGroup.Digits)

	if err != nil {
		t.Errorf(err.Error())
	}

	if preposcharge.Pos.X != 25.45887 ||
		preposcharge.Pos.Y != 74.50821 {
		t.Errorf("preposcharge Pos is incorrect")
	}
	if preposcharge.Charge != 57.69108 {
		t.Errorf("preposcharge Charge is incorrect")
	}
}
