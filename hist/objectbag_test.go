package hist

import (
	"testing"

	"go-hep.org/x/hep/hbook"
)

func TestEmptyBag(t *testing.T) {
	b := objectBag{}
	if b.NObjects() != 0 {
		t.Errorf("Want 0 objects - Got %d\n", b.NObjects())
	}
}

func TestAdd(t *testing.T) {
	b := objectBag{}
	h := hbook.NewH1D(200, 0, 10)
	h.Annotation()["name"] = "htest"
	b.Add(h)
	if b.NObjects() != 1 {
		t.Errorf("Want 1 objects - Got %d\n", b.NObjects())
	}
}
