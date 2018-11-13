package hist

import (
	"testing"

	"go-hep.org/x/hep/hbook"
)

func TestEmptyCollection(t *testing.T) {
	b := Collection{}
	if b.NObjects() != 0 {
		t.Errorf("Want 0 objects - Got %d\n", b.NObjects())
	}
}

func createCollection() Collection {
	c := Collection{}
	h := hbook.NewH1D(20, 0, 20)
	h.Annotation()["name"] = "htest"
	h.Fill(15, 15)
	h.Fill(15, 15)
	c.Add(h)
	h2 := hbook.NewH2D(10, 0, 10, 30, 0, 30)
	h2.Annotation()["name"] = "htest2"
	h2.Fill(5, 15, 75)
	c.Add(h2)
	return c
}

func TestAdd(t *testing.T) {
	c := createCollection()
	const expect int = 2
	if c.NObjects() != expect {
		t.Errorf("Want %d objects - Got %d\n", expect, c.NObjects())
	}
}

func TestRetrieveExistingObjectShouldNotFail(t *testing.T) {
	c := createCollection()
	_, err := c.Get("htest2")
	if err != nil {
		t.Errorf("could not get htest2")
	}
}

func TestRetrieveNonExistentObjectShouldReturnError(t *testing.T) {
	c := createCollection()
	_, err := c.Get("toto")
	want := ErrNonExistingObject
	if err != want {
		t.Errorf("want error %s - got %s", want, err)
	}
}

func TestRetrieveH1D(t *testing.T) {
	c := createCollection()
	h, err := c.H1D("htest")
	if err != nil {
		t.Errorf("could not get htest")
	}
	want := 30.0
	got := (*h).Value(15)
	if got != want {
		t.Errorf("want %v - got %v", want, got)
	}
}

func TestRetrieveH2D(t *testing.T) {
	c := createCollection()
	h2, err := c.H2D("htest2")
	if err != nil {
		t.Errorf("could not get htest2")
	}
	want := 75.0
	got := h2.Binning.Bins[155].Dist.X.SumW()
	if got != want {
		t.Errorf("want %v - got %v", want, got)
	}
}
