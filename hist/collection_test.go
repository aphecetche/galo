package hist

import (
	"fmt"
	"os"
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
	c.Print(os.Stdout)
}

func TestRetrieveObject(t *testing.T) {
	c := createCollection()
	c.Print(os.Stdout)
	h, err := c.Get("htest2")
	if err != nil {
		t.Errorf("could not get htest2")
	} else {
		fmt.Printf("h=%s %p %T\n", (*h).Name(), h, h)
	}
}

func TestRetrieveH1D(t *testing.T) {
	c := createCollection()
	c.Print(os.Stdout)
	h, err := c.H1D("htest")
	if err != nil {
		t.Errorf("could not get htest")
	} else {
		fmt.Printf("h=%s %p %T\n", (*h).Name(), h, h)
		fmt.Printf("h[15]=%g\n", (*h).Value(15))
	}
	c.Print(os.Stdout)
}

func TestRetrieveH2D(t *testing.T) {
	c := createCollection()
	c.Print(os.Stdout)
	h2, err := c.H2D("htest2")
	if err != nil {
		t.Errorf("could not get htest2")
	} else {
		fmt.Printf("h2=%s %p %T\n", (*h2).Name(), h2, h2)
		fmt.Printf("h2[5,10]=%g\n", h2.Binning.Bins[155].Dist.X.SumW())
	}
	c.Print(os.Stdout)
}
