package hist

import (
	"bytes"
	"testing"

	"go-hep.org/x/hep/hbook"
)

func TestCreateCollection(t *testing.T) {

	hc := NewCollection("test")

	if hc.NKeys() != 0 {
		t.Errorf("Want no key - Got %d", hc.NKeys())
	}

	if hc.NObjects() != 0 {
		t.Errorf("Want no object - Got %d", hc.NObjects())
	}
}

func TestAddToCollection(t *testing.T) {

	hc := NewCollection("test")

	h := hbook.NewH1D(10, 0, 10)

	hc.Add("test/a/b/c", h)

	if hc.NKeys() != 1 {
		t.Errorf("Want one key - Got %d", hc.NKeys())
	}

	if hc.NObjects() != 1 {
		t.Errorf("Want one object - Got %d", hc.NObjects())
	}
}

func TestSortAllIdentifiers(t *testing.T) {
	hc := NewCollection("test")
	hc.Add("test/a/b/c", hbook.NewH1D(10, 0, 10))
	hc.Add("test/a/a/b/c", hbook.NewH1D(10, 0, 10))
	ids := hc.SortAllIdentifiers()
	if len(ids) != 2 {
		t.Errorf("Want two ids - Got %d", len(ids))
	}
}

func TestPrint(t *testing.T) {
	hc := NewCollection("test")
	h := hbook.NewH1D(10, 0, 10)
	h.Ann["name"] = "test H1"
	hc.Add("test/a/b/c", h)

	h2 := hbook.NewH2D(10, 0, 10, 20, 0, 20)
	h2.Ann["name"] = "test H2"
	hc.Add("test/a/a/b/c", h2)

	ids := hc.SortAllIdentifiers()
	if len(ids) != 2 {
		t.Errorf("Want two ids - Got %d", len(ids))
	}

	buf := new(bytes.Buffer)
	hc.Print(buf)

	expect := `Number of identifiers 2
KEY /test/a/a/b/c
        OBJ:test H2
KEY /test/a/b/c
        OBJ:test H1
`

	if expect != buf.String() {
		t.Errorf("Did not get expected output from print. Expected:\n%s\nGot:\n%s\n", expect, buf.String())
	}
}

func TestRetrieve(t *testing.T) {
	hc := NewCollection("test")
	h := hbook.NewH1D(10, 0, 10)
	h.Ann["name"] = "h1"
	hc.Add("test/a/b/c", h)
	bin := 5
	expected := 42.42
	h.Fill(5, 42)
	h.Fill(5, 0.42)
	hr := hc.H1D("/test/a/b/c/h1")
	v := hr.Value(bin)
	if v != expected {
		t.Errorf("Expected %v - Got %v\n", expected, v)
	}
}
